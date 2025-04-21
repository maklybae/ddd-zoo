package services

import (
	"context"
	"fmt"

	"github.com/maklybae/ddd-zoo/internal/domain"
)

type ZooStatisticsService interface {
	GetAnimalCount(ctx context.Context) (int, error)
	GetEnclosureCount(ctx context.Context) (int, error)
	GetFreeEnclosureCount(ctx context.Context) (int, error)
	GetFeedingScheduleCount(ctx context.Context) (int, error)
}

type ZooStatistics struct {
	animalRepository          domain.AnimalRepository
	enclosureRepository       domain.EnclosureRepository
	feedingScheduleRepository domain.FeedingScheduleRepository
}

func NewZooStatistics(
	animalRepository domain.AnimalRepository,
	enclosureRepository domain.EnclosureRepository,
	feedingScheduleRepository domain.FeedingScheduleRepository,
) *ZooStatistics {
	return &ZooStatistics{
		animalRepository:          animalRepository,
		enclosureRepository:       enclosureRepository,
		feedingScheduleRepository: feedingScheduleRepository,
	}
}

func (zs *ZooStatistics) GetAnimalCount(ctx context.Context) (int, error) {
	count, err := zs.animalRepository.CountAnimals(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (zs *ZooStatistics) GetEnclosureCount(ctx context.Context) (int, error) {
	count, err := zs.enclosureRepository.CountEnclosures(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting enclosure count: %w", err)
	}

	return count, nil
}

func (zs *ZooStatistics) GetFreeEnclosureCount(ctx context.Context) (int, error) {
	count, err := zs.enclosureRepository.CountFreeEnclosures(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting free enclosure count: %w", err)
	}

	return count, nil
}

func (zs *ZooStatistics) GetFeedingScheduleCount(ctx context.Context) (int, error) {
	count, err := zs.feedingScheduleRepository.CountFeedingSchedules(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting feeding schedule count: %w", err)
	}

	return count, nil
}
