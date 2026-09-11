package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RideFare struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string
	TotalPriceInCents float64
}
