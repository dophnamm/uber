package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID       primitive.ObjectID
	UserID   string
	RideFare RideFare
}

type TripRepository interface {
	Create(ctx context.Context, t Trip) (Trip, error)
}

type TripService interface {
	Create(ctx context.Context, rf RideFare) (*Trip, error)
}
