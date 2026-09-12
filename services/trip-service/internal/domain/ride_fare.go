package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RideFare struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID            string             `json:"userId" bson:"userId"`
	PackageSlug       string             `json:"packageSlug" bson:"packageSlug"`
	TotalPriceInCents float64            `json:"totalPriceInCents" bson:"totalPriceInCents"`
}
