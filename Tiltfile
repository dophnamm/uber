# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### Helpers ###
def dotenv(path):
    if not os.path.exists(path):
        path = path + '.example'
    result = {}
    for line in str(read_file(path)).splitlines():
        line = line.strip()
        if line == '' or line.startswith('#') or '=' not in line:
            continue
        k, v = line.split('=', 1)
        k = k.strip()
        v = v.strip()
        if len(v) > 1 and v[0] == v[-1] and (v[0] == '"' or v[0] == "'"):
            v = v[1:-1]
        # Expand $VAR / ${VAR} against keys already parsed, longest name first so
        # that $MONGO_PORT is not partially matched by a shorter key like $MONGO.
        for ek in sorted(result.keys(), key=len, reverse=True):
            v = v.replace('${%s}' % ek, result[ek]).replace('$%s' % ek, result[ek])
        result[k] = v
    return result

def configmap(name, vars):
    lines = ['apiVersion: v1', 'kind: ConfigMap', 'metadata:', '  name: ' + name, 'data:']
    for k in sorted(vars.keys()):
        lines.append('  %s: "%s"' % (k, vars[k].replace('\\', '\\\\').replace('"', '\\"')))
    return blob('\n'.join(lines) + '\n')

def render(path, vars):
    s = str(read_file(path))
    for k, v in vars.items():
        s = s.replace('{{%s}}' % k, str(v))
    return blob(s)
### End of Helpers ###

### K8s Config ###
# Uncomment to use secrets
# k8s_yaml('./infra/development/k8s/secrets.yaml')

k8s_yaml('./infra/development/k8s/app-config.yaml')
### End of K8s Config ###

### API Gateway ###
# Build to a temp path, then rename atomically. Writing straight to
# build/api-gateway lets Tilt's watcher fire mid-write and sync a half-written
# binary, which makes entr fail with `Text file busy` and crash the container.
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/.tmp-api-gateway ./services/api-gateway/cmd/main.go && mv -f build/.tmp-api-gateway build/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")

gateway_env = dotenv('./services/api-gateway/.env')
watch_file('./services/api-gateway/.env')

docker_build_with_restart(
  'ride-sharing/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    # Stage the binary, then rename it into place inside the container. Syncing
    # straight onto /app/build/api-gateway lets entr exec a file Tilt is still
    # writing, which fails with ETXTBSY and kills the container.
    sync('./build/api-gateway', '/tmp/api-gateway'),
    # ./shared is Go source compiled into the static binary, so nothing reads it
    # at runtime. It stays synced only so edits there use live_update instead of
    # falling back to a full image rebuild.
    sync('./shared', '/app/shared'),
    run('mv -f /tmp/api-gateway /app/build/api-gateway'),
  ],
)

k8s_yaml(configmap('api-gateway-config', gateway_env))

k8s_yaml(render('./infra/development/k8s/api-gateway-deployment.yaml', gateway_env))
k8s_resource('api-gateway',
             port_forwards='%s:%s' % (gateway_env['PORT'], gateway_env['PORT']),
             objects=['api-gateway-config:configmap'],
             resource_deps=['api-gateway-compile'], labels="services")
### End of API Gateway ###

### MongoDB ###
trip_env = dotenv('./services/trip-service/.env')
watch_file('./services/trip-service/.env')

k8s_yaml(render('./infra/development/k8s/mongodb-deployment.yaml', trip_env))
k8s_resource('mongodb',
             port_forwards='%s:%s' % (trip_env['MONGO_PORT'], trip_env['MONGO_PORT']),
             objects=['mongodb-data:persistentvolumeclaim'],
             labels="database")
### End of MongoDB ###

### Trip Service ###
# Build to a temp path, then rename atomically. Writing straight to
# build/trip-service lets Tilt's watcher fire mid-write and sync a half-written
# binary, which makes entr fail with `Text file busy` and crash the container.
trip_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/.tmp-trip-service ./services/trip-service/cmd/main.go && mv -f build/.tmp-trip-service build/trip-service'
if os.name == 'nt':
  trip_compile_cmd = './infra/development/docker/trip-build.bat'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './shared'], labels="compiles")

docker_build_with_restart(
  'ride-sharing/trip-service',
  '.',
  entrypoint=['/app/build/trip-service'],
  dockerfile='./infra/development/docker/trip-service.Dockerfile',
  only=[
    './build/trip-service',
    './shared',
  ],
  live_update=[
    # Stage the binary, then rename it into place inside the container. Syncing
    # straight onto /app/build/trip-service lets entr exec a file Tilt is still
    # writing, which fails with ETXTBSY and kills the container.
    sync('./build/trip-service', '/tmp/trip-service'),
    # ./shared is Go source compiled into the static binary, so nothing reads it
    # at runtime. It stays synced only so edits there use live_update instead of
    # falling back to a full image rebuild.
    sync('./shared', '/app/shared'),
    run('mv -f /tmp/trip-service /app/build/trip-service'),
  ],
)

k8s_yaml(configmap('trip-service-config', trip_env))

k8s_yaml(render('./infra/development/k8s/trip-service-deployment.yaml', trip_env))
k8s_resource('trip-service',
             port_forwards='%s:%s' % (trip_env['PORT'], trip_env['PORT']),
             objects=['trip-service-config:configmap'],
             resource_deps=['trip-service-compile', 'mongodb'], labels="services")
### End of Trip Service ###

### Web Frontend ###
docker_build(
  'ride-sharing/web',
  '.',
  dockerfile='./infra/development/docker/web.Dockerfile',
  # Without `only`, the context is the whole repo and every Go rebuild (which
  # rewrites build/) triggers a full web image rebuild.
  only=['./web'],
)

k8s_yaml(render('./infra/development/k8s/web-deployment.yaml', {'GATEWAY_PORT': gateway_env['PORT']}))
k8s_resource('web', port_forwards='3000:3000', labels="frontend")
### End of Web Frontend ###