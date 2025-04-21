package services

import (
	"context"
	"fmt"

	"github.com/maklybae/ddd-zoo/internal/domain"
	"github.com/maklybae/ddd-zoo/pkg/events"
)

type AnimalTransferService interface {
	TransferAnimal(ctx context.Context, animalID domain.AnimalID, toEnclosureID domain.EnclosureID) error
}

type AnimalTransfer struct {
	animalRepository    domain.AnimalRepository
	enclosureRepository domain.EnclosureRepository
	eventHandler        events.EventHandler
	timeProvider        TimeProvider
}

func NewAnimalTransfer(
	animalRepository domain.AnimalRepository,
	enclosureRepository domain.EnclosureRepository,
	eventHandler events.EventHandler,
	timeProvider TimeProvider,
) *AnimalTransfer {
	return &AnimalTransfer{
		animalRepository:    animalRepository,
		enclosureRepository: enclosureRepository,
		eventHandler:        eventHandler,
		timeProvider:        timeProvider,
	}
}

func (at *AnimalTransfer) TransferAnimal(ctx context.Context, animalID domain.AnimalID, toEnclosureID domain.EnclosureID) error {
	animal, err := at.animalRepository.GetAnimal(ctx, animalID)
	if err != nil {
		return fmt.Errorf("getting animal: %w", err)
	}

	toEnclosure, err := at.enclosureRepository.GetEnclosure(ctx, toEnclosureID)
	if err != nil {
		return fmt.Errorf("getting enclosure: %w", err)
	}

	// Store the old enclosure for the event
	fromEnclosure := animal.Enclosure

	if err = animal.Enclosure.RemoveAnimal(animal); err != nil {
		return fmt.Errorf("removing animal from enclosure: %w", err)
	}

	if err := toEnclosure.AddAnimal(animal); err != nil {
		return fmt.Errorf("adding animal to enclosure: %w", err)
	}

	if err := animal.MoveToEnclosure(toEnclosure); err != nil {
		return fmt.Errorf("moving animal to enclosure: %w", err)
	}

	if err := at.animalRepository.UpdateAnimal(ctx, animal); err != nil {
		return fmt.Errorf("updating animal: %w", err)
	}

	if err := at.enclosureRepository.UpdateEnclosure(ctx, toEnclosure); err != nil {
		return err
	}

	// Publish the AnimalMovedEvent
	movedEvent := domain.AnimalMovedEvent{
		AnimalID:      animal.ID,
		AnimalName:    animal.Name,
		AnimalSpecies: animal.Species,
		FromEnclosure: fromEnclosure,
		ToEnclosure:   toEnclosure,
		Timestamp:     at.timeProvider.Now(),
	}

	if err := at.eventHandler.Handle(ctx, &movedEvent); err != nil {
		return fmt.Errorf("publishing animal moved event: %w", err)
	}

	return nil
}
