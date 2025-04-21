package http

import (
	"github.com/maklybae/ddd-zoo/internal/application/services"
	"github.com/maklybae/ddd-zoo/internal/domain"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
)

type Server struct {
	animalRepo             domain.AnimalRepository
	enclosureRepo          domain.EnclosureRepository
	feedingScheduleRepo    domain.FeedingScheduleRepository
	transferSvc            services.AnimalTransferService
	feedingOrganizationSvc services.FeedingOrganizationService
	statisticsSvc          services.ZooStatisticsService
	timeProvider           services.TimeProvider
}

var _ v1.ServerInterface = (*Server)(nil)

func NewServer(
	animalRepo domain.AnimalRepository,
	enclosureRepo domain.EnclosureRepository,
	feedingScheduleRepo domain.FeedingScheduleRepository,
	transferSvc services.AnimalTransferService,
	feedingOrganizationSvc services.FeedingOrganizationService,
	statisticsSvc services.ZooStatisticsService,
	timeProvider services.TimeProvider,
) *Server {
	return &Server{
		animalRepo:             animalRepo,
		enclosureRepo:          enclosureRepo,
		feedingScheduleRepo:    feedingScheduleRepo,
		transferSvc:            transferSvc,
		feedingOrganizationSvc: feedingOrganizationSvc,
		statisticsSvc:          statisticsSvc,
		timeProvider:           timeProvider,
	}
}
