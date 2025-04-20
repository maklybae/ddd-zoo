package services

import (
	"context"
	"fmt"

	"github.com/maklybae/ddd-zoo/internal/domain"
)

type AnimalTransfer struct {
	animalRepository    domain.AnimalRepository
	enclosureRepository domain.EnclosureRepository
}

func NewAnimalTransfer(animalRepository domain.AnimalRepository, enclosureRepository domain.EnclosureRepository) *AnimalTransfer {
	return &AnimalTransfer{
		animalRepository:    animalRepository,
		enclosureRepository: enclosureRepository,
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

	return nil
}
