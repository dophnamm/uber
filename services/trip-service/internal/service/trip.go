package service

import (
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
)

type TripService struct {
	tripRepo repository.TripRepository
}

func NewTripService(tripRepo repository.TripRepository) *TripService {
	return &TripService{
		tripRepo: tripRepo,
	}
}
