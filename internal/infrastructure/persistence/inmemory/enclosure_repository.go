package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/maklybae/ddd-zoo/internal/domain"
)

// Статическая проверка реализации интерфейса
var _ domain.EnclosureRepository = (*EnclosureRepository)(nil)

type EnclosureRepository struct {
	enclosures map[domain.EnclosureID]*domain.Enclosure
	mutex      sync.RWMutex
}

func NewEnclosureRepository() *EnclosureRepository {
	return &EnclosureRepository{
		enclosures: make(map[domain.EnclosureID]*domain.Enclosure),
	}
}

func (r *EnclosureRepository) GetEnclosure(ctx context.Context, id domain.EnclosureID) (*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enclosure, exists := r.enclosures[id]
	if !exists {
		return nil, fmt.Errorf("enclosure with id %s not found", id)
	}

	return enclosure, nil
}

func (r *EnclosureRepository) AddEnclosure(ctx context.Context, enclosure *domain.Enclosure) error {
	if enclosure.ID == domain.EnclosureID(uuid.Nil) {
		return fmt.Errorf("enclosure id cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.enclosures[enclosure.ID]; exists {
		return fmt.Errorf("enclosure with id %s already exists", enclosure.ID)
	}

	r.enclosures[enclosure.ID] = enclosure
	return nil
}

func (r *EnclosureRepository) DeleteEnclosure(ctx context.Context, id domain.EnclosureID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	enclosure, exists := r.enclosures[id]
	if !exists {
		return fmt.Errorf("enclosure with id %s not found", id)
	}

	// Проверяем, содержит ли вольер животных
	if enclosure.Occupancy.CountAnimals() > 0 {
		return fmt.Errorf("cannot delete enclosure with id %s because it contains animals", id)
	}

	delete(r.enclosures, id)
	return nil
}

func (r *EnclosureRepository) UpdateEnclosure(ctx context.Context, enclosure *domain.Enclosure) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.enclosures[enclosure.ID]; !exists {
		return fmt.Errorf("enclosure with id %s not found", enclosure.ID)
	}

	r.enclosures[enclosure.ID] = enclosure
	return nil
}

func (r *EnclosureRepository) GetAllEnclosures(ctx context.Context) ([]*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enclosures := make([]*domain.Enclosure, 0, len(r.enclosures))
	for _, enclosure := range r.enclosures {
		enclosures = append(enclosures, enclosure)
	}

	return enclosures, nil
}

func (r *EnclosureRepository) CountEnclosures(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.enclosures), nil
}

func (r *EnclosureRepository) CountFreeEnclosures(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, enclosure := range r.enclosures {
		if enclosure.Occupancy.CountAnimals() < enclosure.Occupancy.Capacity {
			count++
		}
	}

	return count, nil
}

// GetEnclosuresByType возвращает все вольеры определенного типа
func (r *EnclosureRepository) GetEnclosuresByType(ctx context.Context, enclosureType domain.EnclosureType) ([]*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var enclosures []*domain.Enclosure
	for _, enclosure := range r.enclosures {
		if enclosure.Type == enclosureType {
			enclosures = append(enclosures, enclosure)
		}
	}

	return enclosures, nil
}

// GetEnclosuresWithSpace возвращает все вольеры, в которых есть свободное место
func (r *EnclosureRepository) GetEnclosuresWithSpace(ctx context.Context) ([]*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var enclosures []*domain.Enclosure
	for _, enclosure := range r.enclosures {
		if enclosure.Occupancy.CountAnimals() < enclosure.Occupancy.Capacity {
			enclosures = append(enclosures, enclosure)
		}
	}

	return enclosures, nil
}
