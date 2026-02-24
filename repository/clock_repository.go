package repository

import "time"

type SystemClockRepository struct{}

func NewSystemClockRepository() *SystemClockRepository {
	return &SystemClockRepository{}
}

func (r *SystemClockRepository) Now() time.Time {
	return time.Now().UTC()
}
