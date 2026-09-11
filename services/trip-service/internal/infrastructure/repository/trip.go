package repository

import "ride-sharing/services/trip-service/internal/domain"

type repository struct {
	trips     map[string]*domain.Trip
	rideFares map[string]*domain.RideFare
}

func NewRepository() *repository {
	return &repository{
		trips:     make(map[string]*domain.Trip),
		rideFares: make(map[string]*domain.RideFare),
	}
}
