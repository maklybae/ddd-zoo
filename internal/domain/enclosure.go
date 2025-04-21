package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrEnclosureFull        = errors.New("enclosure is full")
	ErrAnimalInEnclosure    = errors.New("animal is already in enclosure")
	ErrAnimalNotInEnclosure = errors.New("animal is not in enclosure")
)

type (
	EnclosureID   uuid.UUID
	EnclosureType string
	EnclosureSize int
)

func (eid EnclosureID) String() string {
	return uuid.UUID(eid).String()
}

func (eid EnclosureID) UUID() uuid.UUID {
	return uuid.UUID(eid)
}

// Value Object.
type EnclosureOccupancy struct {
	Capacity int
	Animals  map[*Animal]struct{}
}

func (eo EnclosureOccupancy) CountAnimals() int {
	return len(eo.Animals)
}

func (eo EnclosureOccupancy) AddAnimal(a *Animal) (newOccupancy EnclosureOccupancy, err error) {
	if eo.CountAnimals() >= eo.Capacity {
		return eo, ErrEnclosureFull
	}

	if _, exists := eo.Animals[a]; exists {
		return eo, ErrAnimalInEnclosure
	}

	eo.Animals[a] = struct{}{}

	return EnclosureOccupancy{}, nil
}

func (eo EnclosureOccupancy) RemoveAnimal(a *Animal) (newOccupancy EnclosureOccupancy, err error) {
	if _, exists := eo.Animals[a]; !exists {
		return EnclosureOccupancy{}, ErrAnimalNotInEnclosure
	}

	delete(eo.Animals, a)

	return eo, nil
}

type Enclosure struct {
	ID        EnclosureID
	Type      EnclosureType
	Size      EnclosureSize
	Occupancy EnclosureOccupancy
}

func (e *Enclosure) AddAnimal(a *Animal) error {
	eo, err := e.Occupancy.AddAnimal(a)
	if err != nil {
		return fmt.Errorf("could not add animal to enclosure: %w", err)
	}

	e.Occupancy = eo

	return nil
}

func (e *Enclosure) RemoveAnimal(a *Animal) error {
	eo, err := e.Occupancy.RemoveAnimal(a)
	if err != nil {
		return fmt.Errorf("could not remove animal from enclosure: %w", err)
	}

	e.Occupancy = eo

	return nil
}

func (e *Enclosure) Clean() error {
	return nil
}
