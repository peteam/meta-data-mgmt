package repository

import (
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"gopkg.in/couchbase/gocb.v1"
)

type DbRepository interface {
	Insert(item *entity.Foo) error
	Upsert(item *entity.Foo) error
	Remove(userId string, itemId string) error
	List() ([]*entity.ResourceType, error)
	Count(userId string) (int, error)
	InsertResource(resource *entity.Resource) (*entity.Resource, error)
	DeleteResource(cbKey string) error
	SearchResource(query *gocb.SearchQuery) (gocb.SearchResults, error)
	Retrieve(resourceId string, contentType string) (*entity.Content, error)
	RetrievewithOptional(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.Content, error)
	MultiRetrievewithOptional(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.MultiContent, error)
	PaginationCount(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (int, error)
	Readyz(item *entity.Health) error
	Healthz() error
}
