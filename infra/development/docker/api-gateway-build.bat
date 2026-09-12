set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

REM Build to a temp path, then rename atomically so Tilt never syncs a half-written binary.
go build -o build/.tmp-api-gateway ./services/api-gateway/cmd/main.go
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
move /Y build\.tmp-api-gateway build\api-gateway
