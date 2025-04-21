package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/maklybae/ddd-zoo/internal/domain"
)

// Статическая проверка реализации интерфейса
var _ domain.AnimalRepository = (*AnimalRepository)(nil)

type AnimalRepository struct {
	animals map[domain.AnimalID]*domain.Animal
	mutex   sync.RWMutex
}

func NewAnimalRepository() *AnimalRepository {
	return &AnimalRepository{
		animals: make(map[domain.AnimalID]*domain.Animal),
	}
}

func (r *AnimalRepository) GetAnimal(ctx context.Context, id domain.AnimalID) (*domain.Animal, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	animal, exists := r.animals[id]
	if !exists {
		return nil, fmt.Errorf("animal with id %s not found", id)
	}

	return animal, nil
}

func (r *AnimalRepository) AddAnimal(ctx context.Context, animal *domain.Animal) error {
	if animal.ID == domain.AnimalID(uuid.Nil) {
		return fmt.Errorf("animal id cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.animals[animal.ID]; exists {
		return fmt.Errorf("animal with id %s already exists", animal.ID)
	}

	r.animals[animal.ID] = animal
	return nil
}

func (r *AnimalRepository) DeleteAnimal(ctx context.Context, id domain.AnimalID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.animals[id]; !exists {
		return fmt.Errorf("animal with id %s not found", id)
	}

	delete(r.animals, id)
	return nil
}

func (r *AnimalRepository) UpdateAnimal(ctx context.Context, animal *domain.Animal) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.animals[animal.ID]; !exists {
		return fmt.Errorf("animal with id %s not found", animal.ID)
	}

	r.animals[animal.ID] = animal
	return nil
}

func (r *AnimalRepository) GetAllAnimals(ctx context.Context) ([]*domain.Animal, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	animals := make([]*domain.Animal, 0, len(r.animals))
	for _, animal := range r.animals {
		animals = append(animals, animal)
	}

	return animals, nil
}

func (r *AnimalRepository) CountAnimals(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.animals), nil
}

// GetAnimalsByEnclosure возвращает всех животных, находящихся в указанном вольере
func (r *AnimalRepository) GetAnimalsByEnclosure(ctx context.Context, enclosureID domain.EnclosureID) ([]*domain.Animal, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var animals []*domain.Animal
	for _, animal := range r.animals {
		if animal.Enclosure != nil && animal.Enclosure.ID == enclosureID {
			animals = append(animals, animal)
		}
	}

	return animals, nil
}

// CountHealthyAnimals возвращает количество здоровых животных
func (r *AnimalRepository) CountHealthyAnimals(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, animal := range r.animals {
		if animal.Status == domain.AnimalStatusHealthy {
			count++
		}
	}

	return count, nil
}

// CountSickAnimals возвращает количество больных животных
func (r *AnimalRepository) CountSickAnimals(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, animal := range r.animals {
		if animal.Status == domain.AnimalStatusSick {
			count++
		}
	}

	return count, nil
}
