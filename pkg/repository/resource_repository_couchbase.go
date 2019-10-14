package repository

import (
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"gopkg.in/couchbase/gocb.v1"
)

// InsertResource insert resource into couchbase
func (r *CbRepository) InsertResource(resource *entity.Resource) (*entity.Resource, error) {
	logger.Logger.Debug("Entering CbRepository.InsertResource() ...")
	//index := config.Viper.GetString("database.couchbase.index.metadatamgmt.resourceType.timestamp.desc")
	logger.Logger.Debug(resource)

	_, err := r.Bucket.Insert(resource.Key, &resource.Body, 0)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// DeleteResource insert resource into couchbase
func (r *CbRepository) DeleteResource(cbKey string) error {
	logger.Logger.Debug("Entering CbRepository.DeleteResource() ...")
	//index := config.Viper.GetString("database.couchbase.index.metadatamgmt.resourceType.timestamp.desc")
	logger.Logger.Debug(cbKey)
	cas, err := r.Bucket.Get(cbKey, nil)
	if err != nil {
		return err
	}
	_, err = r.Bucket.Remove(cbKey, cas)
	if err != nil {
		return err
	}

	return nil
}

// InsertResource insert resource into couchbase
func (r *CbRepository) SearchResource(query *gocb.SearchQuery) (gocb.SearchResults, error) {
	logger.Logger.Debug("Entering CbRepository.SearchResource() ...")

	searchResult, err := r.Bucket.ExecuteSearchQuery(query)
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}
