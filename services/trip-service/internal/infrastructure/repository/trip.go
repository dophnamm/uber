package repository

import (
	"ride-sharing/services/trip-service/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
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
	return nil
}
