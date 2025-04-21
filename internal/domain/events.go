package domain

import (
	"time"
)

// AnimalMovedEvent is triggered when an animal is moved to a new enclosure.
type AnimalMovedEvent struct {
	AnimalID      AnimalID
	AnimalName    AnimalName
	AnimalSpecies AnimalSpecies
	FromEnclosure *Enclosure
	ToEnclosure   *Enclosure
	Timestamp     time.Time
}

func (e *AnimalMovedEvent) EventName() string {
	return "animal.moved"
}

// FeedingTimeEvent is triggered when it's time to feed an animal.
type FeedingTimeEvent struct {
	ScheduleID    FeedingScheduleID
	AnimalID      AnimalID
	AnimalName    AnimalName
	AnimalSpecies AnimalSpecies
	Food          Food
	FeedingTime   time.Time
	Timestamp     time.Time
}

func (e *FeedingTimeEvent) EventName() string {
	return "feeding.time"
}
