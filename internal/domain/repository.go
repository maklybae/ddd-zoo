package domain

import (
	"context"
	"time"
)

type AnimalRepository interface {
	GetAnimal(ctx context.Context, id AnimalID) (animal *Animal, err error)
	AddAnimal(ctx context.Context, animal *Animal) error
	DeleteAnimal(ctx context.Context, id AnimalID) error
	UpdateAnimal(ctx context.Context, animal *Animal) error
	GetAllAnimals(ctx context.Context) (animals []*Animal, err error)

	CountAnimals(ctx context.Context) (count int, err error)
	GetAnimalsByEnclosure(ctx context.Context, enclosureID EnclosureID) ([]*Animal, error)
	CountHealthyAnimals(ctx context.Context) (count int, err error)
	CountSickAnimals(ctx context.Context) (count int, err error)
}

type EnclosureRepository interface {
	GetEnclosure(ctx context.Context, id EnclosureID) (enclosure *Enclosure, err error)
	AddEnclosure(ctx context.Context, enclosure *Enclosure) error
	DeleteEnclosure(ctx context.Context, id EnclosureID) error
	UpdateEnclosure(ctx context.Context, enclosure *Enclosure) error
	GetAllEnclosures(ctx context.Context) (enclosures []*Enclosure, err error)

	CountEnclosures(ctx context.Context) (count int, err error)
	CountFreeEnclosures(ctx context.Context) (count int, err error)
	GetEnclosuresByType(ctx context.Context, enclosureType EnclosureType) ([]*Enclosure, error)
	GetEnclosuresWithSpace(ctx context.Context) ([]*Enclosure, error)
}

type FeedingScheduleRepository interface {
	GetFeedingSchedule(ctx context.Context, id FeedingScheduleID) (feedingSchedule *FeedingSchedule, err error)
	AddFeedingSchedule(ctx context.Context, feedingSchedule *FeedingSchedule) error
	DeleteFeedingSchedule(ctx context.Context, id FeedingScheduleID) error
	UpdateFeedingSchedule(ctx context.Context, feedingSchedule *FeedingSchedule) error
	GetAllFeedingSchedules(ctx context.Context) (feedingSchedules []*FeedingSchedule, err error)

	CountFeedingSchedules(ctx context.Context) (count int, err error)
	GetFeedingSchedulesForAnimal(ctx context.Context, animalID AnimalID) ([]*FeedingSchedule, error)
	GetCompletedFeedingSchedules(ctx context.Context) ([]*FeedingSchedule, error)
	GetPendingFeedingSchedules(ctx context.Context) ([]*FeedingSchedule, error)
	GetFeedingSchedulesForTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*FeedingSchedule, error)
	CountCompletedFeedingsToday(ctx context.Context, now time.Time) (int, error)
	CountPendingFeedingsToday(ctx context.Context, now time.Time) (int, error)
}
