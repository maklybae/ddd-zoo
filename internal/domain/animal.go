package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAnimalHealthy = errors.New("animal is already healthy")
	ErrNilEnclosure  = errors.New("enclosure is nil")
)

type (
	AnimalID      uuid.UUID
	AnimalName    string
	AnimalSpecies string
	AnimalStatus  int

	Gender    string
	BirthDate time.Time
)

const (
	AnimalStatusHealthy AnimalStatus = iota
	AnimalStatusSick
)

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type Animal struct {
	ID           AnimalID
	Name         AnimalName
	Species      AnimalSpecies
	BirthDate    BirthDate
	FavoriteFood Food
	Status       AnimalStatus
	Enclosure    *Enclosure
}

func (a *Animal) Feed(food Food) error {
	_ = food
	return nil
}

func (a *Animal) Treat() error {
	if a.Status != AnimalStatusSick {
		return ErrAnimalHealthy
	}

	a.Status = AnimalStatusHealthy

	return nil
}

func (a *Animal) MoveToEnclosure(e *Enclosure) error {
	if e == nil {
		return ErrNilEnclosure
	}

	a.Enclosure = e

	return nil
}
