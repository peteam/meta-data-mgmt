package service

import (
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"gopkg.in/couchbase/gocb.v1"
)

var (
	maxLimitStr = config.Viper.GetString(
		"service.max.limit")
	maxError = config.Viper.GetString(
		"service.max.error")
)

//-- SYSTEM GENEREATED: MODIFY AS REQUIRED ----
type MetaDataMgmtService interface {
	GetContent(resourceId string, contentType string) (*entity.Content, error)
	GetContentOptionalService(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.Content, error)
	GetMultiContentOptionalService(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.MultiContent, error)
	TotalContentCount(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (int, error)
	CountResource(resourceType string) (int, error)
	ListResource() ([]*entity.ResourceType, error)
	AddResource(entity *entity.Resource) (*entity.Resource, error)
	DeleteResource(cbKey string) error
	SearchResourceByFields(pageSize int, pageNumber int, searchEntity *entity.SearchByFieldsBodyEntity) (gocb.SearchResults, error)
	Readyz() error
	Healthz() error
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
type Service struct {
	repo repository.DbRepository
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func NewService(r repository.DbRepository) *Service {
	logger.BootstrapLogger.Debug("Entering Service.NewService() ...")
	return &Service{
		repo: r,
	}
}

func (s *Service) CountResource(resourceType string) (int, error) {
	logger.Logger.Debug("Entering Service.CountResourceTypes() ...")
	return s.repo.Count(resourceType)
}

func (s *Service) ListResource() ([]*entity.ResourceType, error) {
	logger.Logger.Debug("Entering Service.ListResourceTypes() ...")
	return s.repo.List()
}

/*GetContent performs retrieval of a  content entry for a give identifier through Service instance
 *
 */
func (s *Service) GetContent(resourceId string, contentType string) (*entity.Content, error) {
	logger.Logger.Debug("Entering Service.GetMovies() ...")
	return s.repo.Retrieve(resourceId, contentType)
}

/*GetContentOptional performs retrieval of a  content entry for a give identifier through Service instance
 *
 */
func (s *Service) GetContentOptionalService(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.Content, error) {
	logger.Logger.Debug("Entering Service.GetMovies() ...")
	return s.repo.RetrievewithOptional(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

/*GetMultiContentOptionalService performs retrieval of a  content entry for a give identifier through Service instance
 *
 */
func (s *Service) GetMultiContentOptionalService(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.MultiContent, error) {
	logger.Logger.Debug("Entering Service.GetMovies() ...")
	return s.repo.MultiRetrievewithOptional(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

/*TotalContentCount performs retrieval of a count query  for a give query params through Service instance
 *
 */
func (s *Service) TotalContentCount(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (int, error) {
	logger.Logger.Debug("Entering Service.GetMovies() ...")
	return s.repo.PaginationCount(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

//Readyz performs readiness check for the service and underlying components
func (s *Service) Readyz() error {
	logger.Logger.Debug("Entering Service.Readyz() ...")
	health := &entity.Health{
		Health: "dataService_service_health_value",
	}
	if s.repo.Readyz(health) != nil {
		return entity.ErrReadyzFailure
	}
	if s.repo.Healthz() != nil {
		return entity.ErrHealthzFailure
	}
	return nil
}

//Healthz performs health check for the service and underlying components
func (s *Service) Healthz() error {
	logger.Logger.Debug("Entering Service.Healthz() ...")
	if s.repo.Healthz() != nil {
		return entity.ErrHealthzFailure
	}
	return nil
}
