package services

import (
	"context"
	"fmt"
	"time"

	"github.com/maklybae/ddd-zoo/internal/domain"
)

type FeedingOrganization struct {
	animalRepository          domain.AnimalRepository
	feedingScheduleRepository domain.FeedingScheduleRepository
}

func NewFeedingOrganization(
	animalRepository domain.AnimalRepository,
	feedingScheduleRepository domain.FeedingScheduleRepository,
) *FeedingOrganization {
	return &FeedingOrganization{
		animalRepository:          animalRepository,
		feedingScheduleRepository: feedingScheduleRepository,
	}
}

func (fo *FeedingOrganization) FeedAll(ctx context.Context, now time.Time) error {
	animals, err := fo.feedingScheduleRepository.GetAllFeedingSchedules(ctx)
	if err != nil {
		return fmt.Errorf("getting all feeding schedules: %w", err)
	}

	for _, feedingSchedule := range animals {
		if feedingSchedule.IsReady(now) {
			if err := feedingSchedule.Animal.Feed(feedingSchedule.Food); err != nil {
				return fmt.Errorf("feeding animal: %w", err)
			}

			if err := feedingSchedule.Done(); err != nil {
				return fmt.Errorf("marking feeding schedule as done: %w", err)
			}
		}

		if err := fo.animalRepository.UpdateAnimal(ctx, feedingSchedule.Animal); err != nil {
			return fmt.Errorf("updating animal: %w", err)
		}

		if err := fo.feedingScheduleRepository.UpdateFeedingSchedule(ctx, feedingSchedule); err != nil {
			return fmt.Errorf("updating feeding schedule: %w", err)
		}
	}

	return nil
}
