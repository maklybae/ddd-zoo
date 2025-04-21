package adapters

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maklybae/ddd-zoo/internal/domain"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
)

var (
	ErrAnimalNotFound          = errors.New("animal not found")
	ErrFeedingScheduleNotFound = errors.New("feeding schedule not found")
)

func DomainFeedingScheduleToAPI(schedule *domain.FeedingSchedule) v1.FeedingSchedule {
	if schedule == nil {
		return v1.FeedingSchedule{}
	}

	var animal v1.Animal
	if schedule.Animal != nil {
		animal = DomainAnimalToAPI(schedule.Animal)
	}

	return v1.FeedingSchedule{
		Id:          schedule.ID.UUID(),
		Animal:      animal,
		FeedingTime: time.Time(schedule.Time),
		FoodType:    string(schedule.Food),
		Completed:   schedule.Status == domain.FeedingStatusDone,
	}
}

func APIToNewDomainFeedingSchedule(input v1.FeedingScheduleInput, animal *domain.Animal) (*domain.FeedingSchedule, error) {
	if animal == nil {
		return nil, ErrAnimalNotFound
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &domain.FeedingSchedule{
		ID:     domain.FeedingScheduleID(id),
		Animal: animal,
		Food:   domain.Food(input.FoodType),
		Time:   domain.FeedingScheduleTime(input.FeedingTime),
		Status: domain.FeedingStatusNotDone,
	}, nil
}

func DomainFeedingScheduleToAPIList(schedules []*domain.FeedingSchedule) []v1.FeedingSchedule {
	if schedules == nil {
		return []v1.FeedingSchedule{}
	}

	result := make([]v1.FeedingSchedule, len(schedules))
	for i, schedule := range schedules {
		result[i] = DomainFeedingScheduleToAPI(schedule)
	}

	return result
}
