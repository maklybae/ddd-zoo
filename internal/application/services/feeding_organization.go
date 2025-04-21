package services

import (
	"context"
	"fmt"
	"time"

	"github.com/maklybae/ddd-zoo/internal/domain"
	"github.com/maklybae/ddd-zoo/pkg/events"
)

type FeedingOrganizationService interface {
	FeedAll(ctx context.Context, now time.Time) error
}

type FeedingOrganization struct {
	animalRepository          domain.AnimalRepository
	feedingScheduleRepository domain.FeedingScheduleRepository
	eventDispatcher           events.Dispatcher
	timeProvider              TimeProvider
}

func NewFeedingOrganization(
	animalRepository domain.AnimalRepository,
	feedingScheduleRepository domain.FeedingScheduleRepository,
	eventDispatcher events.Dispatcher,
	timeProvider TimeProvider,
) *FeedingOrganization {
	return &FeedingOrganization{
		animalRepository:          animalRepository,
		feedingScheduleRepository: feedingScheduleRepository,
		eventDispatcher:           eventDispatcher,
		timeProvider:              timeProvider,
	}
}

func (fo *FeedingOrganization) FeedAll(ctx context.Context, now time.Time) error {
	feedingSchedules, err := fo.feedingScheduleRepository.GetAllFeedingSchedules(ctx)
	if err != nil {
		return fmt.Errorf("getting all feeding schedules: %w", err)
	}

	for _, feedingSchedule := range feedingSchedules {
		if feedingSchedule.IsReady(now) {
			// Publish FeedingTimeEvent before feeding the animal
			feedingEvent := domain.FeedingTimeEvent{
				ScheduleID:    feedingSchedule.ID,
				AnimalID:      feedingSchedule.Animal.ID,
				AnimalName:    feedingSchedule.Animal.Name,
				AnimalSpecies: feedingSchedule.Animal.Species,
				Food:          feedingSchedule.Food,
				FeedingTime:   time.Time(feedingSchedule.Time),
				Timestamp:     fo.timeProvider.Now(),
			}

			fo.eventDispatcher.Dispatch(ctx, &feedingEvent)

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
