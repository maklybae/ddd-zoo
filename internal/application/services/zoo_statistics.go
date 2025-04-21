package services

import (
	"context"
	"fmt"
	"time"

	"github.com/maklybae/ddd-zoo/internal/domain"
)

type ZooStatisticsService interface {
	GetAnimalCount(ctx context.Context) (int, error)
	GetEnclosureCount(ctx context.Context) (int, error)
	GetFreeEnclosureCount(ctx context.Context) (int, error)
	GetFeedingScheduleCount(ctx context.Context) (int, error)
	GetHealthyAnimalCount(ctx context.Context) (int, error)
	GetSickAnimalCount(ctx context.Context) (int, error)
	GetCompletedFeedingsTodayCount(ctx context.Context) (int, error)
	GetPendingFeedingsTodayCount(ctx context.Context) (int, error)
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

func (zs *ZooStatistics) GetHealthyAnimalCount(ctx context.Context) (int, error) {
	count, err := zs.animalRepository.CountHealthyAnimals(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting healthy animal count: %w", err)
	}

	return count, nil
}

func (zs *ZooStatistics) GetSickAnimalCount(ctx context.Context) (int, error) {
	count, err := zs.animalRepository.CountSickAnimals(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting sick animal count: %w", err)
	}

	return count, nil
}

func (zs *ZooStatistics) GetCompletedFeedingsTodayCount(ctx context.Context) (int, error) {
	// Get current time from time provider if needed
	now := time.Now()

	count, err := zs.feedingScheduleRepository.CountCompletedFeedingsToday(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("getting completed feedings count: %w", err)
	}

	return count, nil
}

func (zs *ZooStatistics) GetPendingFeedingsTodayCount(ctx context.Context) (int, error) {
	// Get current time from time provider if needed
	now := time.Now()

	count, err := zs.feedingScheduleRepository.CountPendingFeedingsToday(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("getting pending feedings count: %w", err)
	}

	return count, nil
}
