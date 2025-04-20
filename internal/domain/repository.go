package domain

import "context"

type AnimalRepository interface {
	GetAnimal(ctx context.Context, id AnimalID) (animal *Animal, err error)
	AddAnimal(ctx context.Context, animal *Animal) error
	DeleteAnimal(ctx context.Context, id AnimalID) error
	UpdateAnimal(ctx context.Context, animal *Animal) error
	GetAllAnimals(ctx context.Context) (animals []*Animal, err error)
}

type EnclosureRepository interface {
	GetEnclosure(ctx context.Context, id EnclosureID) (enclosure *Enclosure, err error)
	AddEnclosure(ctx context.Context, enclosure *Enclosure) error
	DeleteEnclosure(ctx context.Context, id EnclosureID) error
	UpdateEnclosure(ctx context.Context, enclosure *Enclosure) error
	GetAllEnclosures(ctx context.Context) (enclosures []*Enclosure, err error)
}

type FeedingScheduleRepository interface {
	GetFeedingSchedule(ctx context.Context, id FeedingScheduleID) (feedingSchedule *FeedingSchedule, err error)
	AddFeedingSchedule(ctx context.Context, feedingSchedule *FeedingSchedule) error
	DeleteFeedingSchedule(ctx context.Context, id FeedingScheduleID) error
	UpdateFeedingSchedule(ctx context.Context, feedingSchedule *FeedingSchedule) error
	GetAllFeedingSchedules(ctx context.Context) (feedingSchedules []*FeedingSchedule, err error)
}
