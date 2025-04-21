package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/maklybae/ddd-zoo/internal/domain"
)

// Статическая проверка реализации интерфейса
var _ domain.FeedingScheduleRepository = (*FeedingScheduleRepository)(nil)

type FeedingScheduleRepository struct {
	schedules map[domain.FeedingScheduleID]*domain.FeedingSchedule
	mutex     sync.RWMutex
}

func NewFeedingScheduleRepository() *FeedingScheduleRepository {
	return &FeedingScheduleRepository{
		schedules: make(map[domain.FeedingScheduleID]*domain.FeedingSchedule),
	}
}

func (r *FeedingScheduleRepository) GetFeedingSchedule(ctx context.Context, id domain.FeedingScheduleID) (*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	schedule, exists := r.schedules[id]
	if !exists {
		return nil, fmt.Errorf("feeding schedule with id %s not found", id)
	}

	return schedule, nil
}

func (r *FeedingScheduleRepository) AddFeedingSchedule(ctx context.Context, schedule *domain.FeedingSchedule) error {
	if schedule.ID == domain.FeedingScheduleID(uuid.Nil) {
		return fmt.Errorf("feeding schedule id cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.schedules[schedule.ID]; exists {
		return fmt.Errorf("feeding schedule with id %s already exists", schedule.ID)
	}

	r.schedules[schedule.ID] = schedule
	return nil
}

func (r *FeedingScheduleRepository) DeleteFeedingSchedule(ctx context.Context, id domain.FeedingScheduleID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.schedules[id]; !exists {
		return fmt.Errorf("feeding schedule with id %s not found", id)
	}

	delete(r.schedules, id)
	return nil
}

func (r *FeedingScheduleRepository) UpdateFeedingSchedule(ctx context.Context, schedule *domain.FeedingSchedule) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.schedules[schedule.ID]; !exists {
		return fmt.Errorf("feeding schedule with id %s not found", schedule.ID)
	}

	r.schedules[schedule.ID] = schedule
	return nil
}

func (r *FeedingScheduleRepository) GetAllFeedingSchedules(ctx context.Context) ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	schedules := make([]*domain.FeedingSchedule, 0, len(r.schedules))
	for _, schedule := range r.schedules {
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *FeedingScheduleRepository) CountFeedingSchedules(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.schedules), nil
}

// GetFeedingSchedulesForAnimal возвращает все расписания кормлений для конкретного животного
func (r *FeedingScheduleRepository) GetFeedingSchedulesForAnimal(ctx context.Context, animalID domain.AnimalID) ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var schedules []*domain.FeedingSchedule
	for _, schedule := range r.schedules {
		if schedule.Animal.ID == animalID {
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// GetCompletedFeedingSchedules возвращает все выполненные расписания кормлений
func (r *FeedingScheduleRepository) GetCompletedFeedingSchedules(ctx context.Context) ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var schedules []*domain.FeedingSchedule
	for _, schedule := range r.schedules {
		if schedule.Status == domain.FeedingStatusDone {
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// GetPendingFeedingSchedules возвращает все ожидающие расписания кормлений
func (r *FeedingScheduleRepository) GetPendingFeedingSchedules(ctx context.Context) ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var schedules []*domain.FeedingSchedule
	for _, schedule := range r.schedules {
		if schedule.Status == domain.FeedingStatusNotDone {
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// GetFeedingSchedulesForTimeRange возвращает все расписания кормлений в заданном временном диапазоне
func (r *FeedingScheduleRepository) GetFeedingSchedulesForTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var schedules []*domain.FeedingSchedule
	for _, schedule := range r.schedules {
		scheduleTime := time.Time(schedule.Time)
		if (scheduleTime.Equal(startTime) || scheduleTime.After(startTime)) &&
			(scheduleTime.Equal(endTime) || scheduleTime.Before(endTime)) {
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// CountCompletedFeedingsToday возвращает количество выполненных кормлений за сегодня
func (r *FeedingScheduleRepository) CountCompletedFeedingsToday(ctx context.Context, now time.Time) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	count := 0
	for _, schedule := range r.schedules {
		scheduleTime := time.Time(schedule.Time)
		if schedule.Status == domain.FeedingStatusDone &&
			scheduleTime.After(startOfDay) && scheduleTime.Before(endOfDay) {
			count++
		}
	}

	return count, nil
}

// CountPendingFeedingsToday возвращает количество ожидающих кормлений на сегодня
func (r *FeedingScheduleRepository) CountPendingFeedingsToday(ctx context.Context, now time.Time) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	count := 0
	for _, schedule := range r.schedules {
		scheduleTime := time.Time(schedule.Time)
		if schedule.Status == domain.FeedingStatusNotDone &&
			scheduleTime.After(startOfDay) && scheduleTime.Before(endOfDay) {
			count++
		}
	}

	return count, nil
}
