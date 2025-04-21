package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maklybae/ddd-zoo/internal/domain"
	"github.com/maklybae/ddd-zoo/internal/presentation/http/adapters"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Get all animals
// (GET /api/v1/animals)
func (server *Server) GetApiV1Animals(c *gin.Context) {
	animals, err := server.animalRepo.GetAllAnimals(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert domain animals to API animals
	apiAnimals := adapters.DomainAnimalToAPIList(animals)

	c.JSON(http.StatusOK, v1.AnimalListResponse{
		Animals: apiAnimals,
	})
}

// Add a new animal
// (POST /api/v1/animals)
func (server *Server) PostApiV1Animals(c *gin.Context) {
	// Parse the request body
	var input v1.AnimalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert API input to domain animal
	animal, err := adapters.APIToNewDomainAnimal(input)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	enclosure, err := server.enclosureRepo.GetEnclosure(c.Request.Context(), animal.Enclosure.ID)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	animal.Enclosure = enclosure

	// Create a new animal
	err = server.animalRepo.AddAnimal(c.Request.Context(), animal)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	var gender v1.AnimalGender = v1.AnimalGenderFemale
	if animal.Gender == domain.Male {
		gender = v1.AnimalGenderMale
	}

	var status v1.AnimalStatus = v1.AnimalStatusHealthy
	if animal.Status == domain.AnimalStatusSick {
		status = v1.AnimalStatusSick
	}

	c.JSON(http.StatusCreated, v1.Animal{
		Id:           animal.ID.UUID(),
		EnclosureId:  animal.Enclosure.ID.UUID(),
		FavoriteFood: string(animal.FavoriteFood),
		Gender:       gender,
		Species:      string(animal.Species),
		Name:         string(animal.Name),
		BirthDate:    time.Time(animal.BirthDate),
		Status:       status,
	})
}

// Delete an animal
// (DELETE /api/v1/animals/{animalId})
func (server *Server) DeleteApiV1AnimalsAnimalId(c *gin.Context, animalId openapi_types.UUID) {
	animalIdDomain := domain.AnimalID(animalId)

	if err := server.animalRepo.DeleteAnimal(c.Request.Context(), animalIdDomain); err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// Get animal by ID
// (GET /api/v1/animals/{animalId})
func (server *Server) GetApiV1AnimalsAnimalId(c *gin.Context, animalId openapi_types.UUID) {
	animalIdDomain := domain.AnimalID(animalId)

	animal, err := server.animalRepo.GetAnimal(c.Request.Context(), animalIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert domain animal to API animal
	apiAnimal := adapters.DomainAnimalToAPI(animal)

	c.JSON(http.StatusOK, apiAnimal)
}

// Move an animal to a new enclosure
// (POST /api/v1/animals/{animalId}/move)
func (server *Server) PostApiV1AnimalsAnimalIdMove(c *gin.Context, animalId openapi_types.UUID) {
	animalIdDomain := domain.AnimalID(animalId)

	// Parse the request body
	var input v1.MoveAnimalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	newEnclosureId := domain.EnclosureID(input.NewEnclosureId)

	// Transfer the animal
	err := server.transferSvc.TransferAnimal(c.Request.Context(), animalIdDomain, newEnclosureId)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Get the updated animal
	animal, err := server.animalRepo.GetAnimal(c.Request.Context(), animalIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Return the updated animal
	apiAnimal := adapters.DomainAnimalToAPI(animal)
	c.JSON(http.StatusOK, apiAnimal)
}

// Treat a sick animal
// (POST /api/v1/animals/{animalId}/treat)
func (server *Server) PostApiV1AnimalsAnimalIdTreat(c *gin.Context, animalId openapi_types.UUID) {
	animalIdDomain := domain.AnimalID(animalId)

	// Get the animal
	animal, err := server.animalRepo.GetAnimal(c.Request.Context(), animalIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	if err = animal.Treat(); err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Save the updated animal
	err = server.animalRepo.UpdateAnimal(c.Request.Context(), animal)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Return the updated animal
	apiAnimal := adapters.DomainAnimalToAPI(animal)
	c.JSON(http.StatusOK, apiAnimal)
}

// Get all enclosures
// (GET /api/v1/enclosures)
func (server *Server) GetApiV1Enclosures(c *gin.Context) {
	enclosures, err := server.enclosureRepo.GetAllEnclosures(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert domain enclosures to API enclosures
	apiEnclosures := adapters.DomainEnclosureToAPIList(enclosures)

	c.JSON(http.StatusOK, v1.EnclosureListResponse{
		Enclosures: apiEnclosures,
	})
}

// Add a new enclosure
// (POST /api/v1/enclosures)
func (server *Server) PostApiV1Enclosures(c *gin.Context) {
	// Parse the request body
	var input v1.EnclosureInput
	if err := c.ShouldBindJSON(&input); err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Create a new enclosure
	enclosure, err := adapters.APIToNewDomainEnclosure(input)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Save the enclosure
	err = server.enclosureRepo.AddEnclosure(c.Request.Context(), enclosure)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Return the created enclosure
	apiEnclosure := adapters.DomainEnclosureToAPI(enclosure)
	c.JSON(http.StatusCreated, apiEnclosure)
}

// Delete an enclosure
// (DELETE /api/v1/enclosures/{enclosureId})
func (server *Server) DeleteApiV1EnclosuresEnclosureId(c *gin.Context, enclosureId openapi_types.UUID) {
	enclosureIdDomain := domain.EnclosureID(enclosureId)

	// Delete the enclosure
	err := server.enclosureRepo.DeleteEnclosure(c.Request.Context(), enclosureIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// Get enclosure by ID
// (GET /api/v1/enclosures/{enclosureId})
func (server *Server) GetApiV1EnclosuresEnclosureId(c *gin.Context, enclosureId openapi_types.UUID) {
	enclosureIdDomain := domain.EnclosureID(enclosureId)

	// Get the enclosure
	enclosure, err := server.enclosureRepo.GetEnclosure(c.Request.Context(), enclosureIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert domain enclosure to API enclosure
	apiEnclosure := adapters.DomainEnclosureToAPI(enclosure)

	c.JSON(http.StatusOK, apiEnclosure)
}

// Clean an enclosure
// (POST /api/v1/enclosures/{enclosureId}/clean)
func (server *Server) PostApiV1EnclosuresEnclosureIdClean(c *gin.Context, enclosureId openapi_types.UUID) {
	enclosureIdDomain := domain.EnclosureID(enclosureId)

	// Get the enclosure
	enclosure, err := server.enclosureRepo.GetEnclosure(c.Request.Context(), enclosureIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Trigger the cleaned event
	enclosure.Clean()

	// Save the updated enclosure
	err = server.enclosureRepo.UpdateEnclosure(c.Request.Context(), enclosure)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Return the updated enclosure
	apiEnclosure := adapters.DomainEnclosureToAPI(enclosure)
	c.JSON(http.StatusOK, apiEnclosure)
}

// Get all feeding schedules
// (GET /api/v1/feeding-schedules)
func (server *Server) GetApiV1FeedingSchedules(c *gin.Context) {
	schedules, err := server.feedingScheduleRepo.GetAllFeedingSchedules(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert domain schedules to API schedules
	apiSchedules := adapters.DomainFeedingScheduleToAPIList(schedules)

	c.JSON(http.StatusOK, v1.FeedingScheduleListResponse{
		Schedules: apiSchedules,
	})
}

// Add a new feeding schedule
// (POST /api/v1/feeding-schedules)
func (server *Server) PostApiV1FeedingSchedules(c *gin.Context) {
	// Parse the request body
	var input v1.FeedingScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	animalId := domain.AnimalID(input.AnimalId)

	// Get the animal
	animal, err := server.animalRepo.GetAnimal(c.Request.Context(), animalId)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Create a new feeding schedule
	schedule, err := adapters.APIToNewDomainFeedingSchedule(input, animal)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Save the feeding schedule
	err = server.feedingScheduleRepo.AddFeedingSchedule(c.Request.Context(), schedule)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Return the created feeding schedule
	apiSchedule := adapters.DomainFeedingScheduleToAPI(schedule)
	c.JSON(http.StatusCreated, apiSchedule)
}

// Delete a feeding schedule
// (DELETE /api/v1/feeding-schedules/{scheduleId})
func (server *Server) DeleteApiV1FeedingSchedulesScheduleId(c *gin.Context, scheduleId openapi_types.UUID) {
	scheduleIdDomain := domain.FeedingScheduleID(scheduleId)

	err := server.feedingScheduleRepo.DeleteFeedingSchedule(c.Request.Context(), scheduleIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// Get feeding schedule by ID
// (GET /api/v1/feeding-schedules/{scheduleId})
func (server *Server) GetApiV1FeedingSchedulesScheduleId(c *gin.Context, scheduleId openapi_types.UUID) {
	scheduleIdDomain := domain.FeedingScheduleID(scheduleId)

	schedule, err := server.feedingScheduleRepo.GetFeedingSchedule(c.Request.Context(), scheduleIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	// Convert domain schedule to API schedule
	apiSchedule := adapters.DomainFeedingScheduleToAPI(schedule)

	c.JSON(http.StatusOK, apiSchedule)
}

// Mark a feeding schedule as completed
// (POST /api/v1/feeding-schedules/{scheduleId}/complete)
func (server *Server) PostApiV1FeedingSchedulesScheduleIdComplete(c *gin.Context, scheduleId openapi_types.UUID) {
	scheduleIdDomain := domain.FeedingScheduleID(scheduleId)

	// Process the completion through the service
	schedule, err := server.feedingScheduleRepo.GetFeedingSchedule(c.Request.Context(), scheduleIdDomain)
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	schedule.Done()

	// Return the updated schedule
	apiSchedule := adapters.DomainFeedingScheduleToAPI(schedule)
	c.JSON(http.StatusOK, apiSchedule)
}

// Get zoo statistics
// (GET /api/v1/statistics)
func (server *Server) GetApiV1Statistics(c *gin.Context) {
	// Get the statistics
	totalAnimals, err := server.statisticsSvc.GetAnimalCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	healthyAnimals, err := server.statisticsSvc.GetHealthyAnimalCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	sickAnimals, err := server.statisticsSvc.GetSickAnimalCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	totalEnclosures, err := server.statisticsSvc.GetEnclosureCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	freeEnclosures, err := server.statisticsSvc.GetFreeEnclosureCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	feedingSchedulesCount, err := server.statisticsSvc.GetFeedingScheduleCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	completedFeedingsToday, err := server.statisticsSvc.GetCompletedFeedingsTodayCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	pendingFeedingsToday, err := server.statisticsSvc.GetPendingFeedingsTodayCount(c.Request.Context())
	if err != nil {
		server.SendBadRequestResponse(c, err, nil)
		return
	}

	stats := v1.ZooStatistics{
		TotalAnimals:           totalAnimals,
		HealthyAnimals:         healthyAnimals,
		SickAnimals:            sickAnimals,
		TotalEnclosures:        totalEnclosures,
		FreeEnclosures:         freeEnclosures,
		FeedingSchedulesCount:  feedingSchedulesCount,
		CompletedFeedingsToday: completedFeedingsToday,
		PendingFeedingsToday:   pendingFeedingsToday,
	}

	c.JSON(http.StatusOK, stats)
}
