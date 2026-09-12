package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID   string             `json:"userId" bson:"userId"`
	RideFare RideFare           `json:"rideFare" bson:"rideFare"`
}
