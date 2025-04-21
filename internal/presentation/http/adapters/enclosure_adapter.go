package adapters

import (
	"github.com/google/uuid"
	"github.com/maklybae/ddd-zoo/internal/domain"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
)

func DomainEnclosureToAPI(enclosure *domain.Enclosure) v1.Enclosure {
	if enclosure == nil {
		return v1.Enclosure{}
	}

	animals := make([]v1.Animal, 0)

	if enclosure.Occupancy.Animals != nil {
		for animal := range enclosure.Occupancy.Animals {
			if animal != nil {
				animals = append(animals, DomainAnimalToAPI(animal))
			}
		}
	}

	return v1.Enclosure{
		Id:             enclosure.ID.UUID(),
		Animals:        &animals,
		Type:           string(enclosure.Type),
		Size:           int(enclosure.Size),
		CurrentAnimals: enclosure.Occupancy.CountAnimals(),
		MaxCapacity:    enclosure.Occupancy.Capacity,
	}
}

func APIToNewDomainEnclosure(input v1.EnclosureInput) (*domain.Enclosure, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &domain.Enclosure{
		ID:   domain.EnclosureID(id),
		Type: domain.EnclosureType(input.Type),
		Size: domain.EnclosureSize(input.Size),
		Occupancy: domain.EnclosureOccupancy{
			Capacity: input.MaxCapacity,
			Animals:  make(map[*domain.Animal]struct{}),
		},
	}, nil
}

func DomainEnclosureToAPIList(enclosures []*domain.Enclosure) []v1.Enclosure {
	if enclosures == nil {
		return []v1.Enclosure{}
	}

	result := make([]v1.Enclosure, len(enclosures))
	for i, enclosure := range enclosures {
		result[i] = DomainEnclosureToAPI(enclosure)
	}

	return result
}
