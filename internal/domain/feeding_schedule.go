package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFeedingStatusIsDone = errors.New("feeding status is already done")
)

type (
	FeedingScheduleID   uuid.UUID
	FeedingScheduleTime time.Time
	FeedingStatus       bool
)

const (
	FeedingStatusDone    FeedingStatus = true
	FeedingStatusNotDone FeedingStatus = false
)

func (fs FeedingStatus) Done() (newFeedingStatus FeedingStatus, err error) {
	if fs == FeedingStatusDone {
		return FeedingStatusDone, ErrFeedingStatusIsDone
	}

	fs = FeedingStatusDone

	return fs, nil
}

type FeedingSchedule struct {
	ID     FeedingScheduleID
	Animal *Animal
	Food   Food
	Time   FeedingScheduleTime
	Status FeedingStatus
}

func (fs *FeedingSchedule) ChangeTime(newTime FeedingScheduleTime) error {
	fs.Time = newTime
	return nil
}

func (fs *FeedingSchedule) Done() error {
	status, err := fs.Status.Done()
	if err != nil {
		return fmt.Errorf("marking feeding schedule as done: %w", err)
	}

	fs.Status = status

	return nil
}
