package service

import (
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
)

// AddResource add resource into couchbase
func (s *Service) AddResource(entity *entity.Resource) (*entity.Resource, error) {
	logger.Logger.Debug("Entering Service.InsertResource() ...")

	return s.repo.InsertResource(entity)
}

// DeleteResource delete resource into couchbase
func (s *Service) DeleteResource(cbKey string) error {
	logger.Logger.Debug("Entering Service.InsertResource() ...")

	return s.repo.DeleteResource(cbKey)
}
