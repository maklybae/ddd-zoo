package domain

import (
	"time"

	"github.com/maklybae/ddd-zoo/pkg/events"
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

var _ events.Event = (*AnimalMovedEvent)(nil)

func (e *AnimalMovedEvent) Name() string {
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

var _ events.Event = (*FeedingTimeEvent)(nil)

func (e *FeedingTimeEvent) Name() string {
	return "feeding.time"
}
