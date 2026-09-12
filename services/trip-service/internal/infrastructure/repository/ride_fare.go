package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"ride-sharing/services/trip-service/internal/domain"
)

type RideFareRepository interface {
	Create(rf *domain.RideFare) (*domain.RideFare, error)
}

type rideFareRepository struct {
	collection *mongo.Collection
}

func NewRideFareRepository(db *mongo.Database) RideFareRepository {
	return &rideFareRepository{
		collection: db.Collection("ride_fares"),
	}
}

func (r *rideFareRepository) Create(rf *domain.RideFare) (*domain.RideFare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trip, err := r.collection.InsertOne(ctx, rf)
	if err != nil {
		return nil, err
	}

	id, ok := trip.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("could not convert inserted ID to ObjectID")
	}

	rf.ID = id

	return rf, nil
}
