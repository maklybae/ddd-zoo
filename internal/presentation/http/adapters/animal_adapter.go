package adapters

import (
	"time"

	"github.com/google/uuid"
	"github.com/maklybae/ddd-zoo/internal/domain"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
)

func DomainAnimalToAPI(animal *domain.Animal) v1.Animal {
	if animal == nil {
		return v1.Animal{}
	}

	var enclosureID uuid.UUID
	if animal.Enclosure != nil {
		enclosureID = uuid.UUID(animal.Enclosure.ID)
	}

	status := v1.AnimalStatusHealthy
	if animal.Status == domain.AnimalStatusSick {
		status = v1.AnimalStatusSick
	}

	gender := v1.AnimalGenderMale
	if animal.Gender == domain.Female {
		gender = v1.AnimalGenderFemale
	}

	birthDate, _ := time.Parse(time.RFC3339, time.Time(animal.BirthDate).Format(time.RFC3339))

	return v1.Animal{
		Id:           animal.ID.UUID(),
		EnclosureId:  enclosureID,
		Species:      string(animal.Species),
		Name:         string(animal.Name),
		BirthDate:    birthDate,
		Gender:       gender,
		FavoriteFood: string(animal.FavoriteFood),
		Status:       status,
	}
}

func APIToNewDomainAnimal(input v1.AnimalInput) (*domain.Animal, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	birthDate := input.BirthDate

	gender := domain.Male
	if input.Gender == v1.AnimalInputGenderFemale {
		gender = domain.Female
	}

	status := domain.AnimalStatusHealthy
	if input.Status == v1.AnimalInputStatusSick {
		status = domain.AnimalStatusSick
	}

	return &domain.Animal{
		ID:           domain.AnimalID(id),
		Species:      domain.AnimalSpecies(input.Species),
		Name:         domain.AnimalName(input.Name),
		BirthDate:    domain.BirthDate(birthDate),
		Gender:       gender,
		FavoriteFood: domain.Food(input.FavoriteFood),
		Status:       status,
	}, nil
}

func DomainAnimalToAPIList(animals []*domain.Animal) []v1.Animal {
	if animals == nil {
		return []v1.Animal{}
	}

	result := make([]v1.Animal, len(animals))
	for i, animal := range animals {
		result[i] = DomainAnimalToAPI(animal)
	}

	return result
}
