package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"ride-sharing/services/trip-service/internal/domain"
)

type TripRepository interface {
	Create(t *domain.Trip) error
}

type tripRepository struct {
	collection *mongo.Collection
}

func NewTripRepository(db *mongo.Database) TripRepository {
	return &tripRepository{
		collection: db.Collection("trips"),
	}
}

func (r *tripRepository) Create(t *domain.Trip) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trip, err := r.collection.InsertOne(ctx, t)
	if err != nil {
		return err
	}

	id, ok := trip.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("could not convert inserted ID to ObjectID")
	}

	t.ID = id

	return nil
}
