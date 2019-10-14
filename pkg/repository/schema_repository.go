package repository

import (
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
)

type SchemaRepository interface {
	LoadSchemas(locationDir string) (map[string]*entity.Schema, error)
}
