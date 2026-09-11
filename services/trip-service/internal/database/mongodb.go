package database

import (
	"context"
	"fmt"
	"time"

	"ride-sharing/services/trip-service/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectTimeout = 10 * time.Second

func InitDB(ctx context.Context, cfg *config.Config) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(cfg.MongoUri)

	if cfg.MongoUserName != "" {
		opts.SetAuth(options.Credential{
			AuthSource:  cfg.MongoAuthSource,
			Username:    cfg.MongoUserName,
			Password:    cfg.MongoPassword,
			PasswordSet: cfg.MongoPassword != "",
		})
	}

	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("Connecting to mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("Pinging mongodb: %w", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	return client.Database(cfg.MongoDBName), nil
}
