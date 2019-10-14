package service

import (
	"errors"
	"fmt"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"github.com/xeipuuv/gojsonschema"
)

//-- SYSTEM GENEREATED: MODIFY AS REQUIRED ----
type SchemaLocalService interface {
	ReloadAllSchemas(location string) (map[string]*entity.Schema, error)
	Validate(json string, jsonSchema string) (*entity.Schema, error)
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
type SchemaService struct {
	repo repository.SchemaRepository
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func NewSchemaService(r repository.SchemaRepository) *SchemaService {
	logger.BootstrapLogger.Debug("Entering Service.NewSchemaService() ...")
	return &SchemaService{
		repo: r,
	}
}

//-- SYSTEM GENEREATED: MODIFY AS REQUIRED ----
func (s *SchemaService) ReloadAllSchemas(location string) (map[string]*entity.Schema, error) {
	logger.Logger.Debug("Entering SchemaService.ReloadAllSchemas() ...")

	return s.repo.LoadSchemas(location)
}

func (s *SchemaService) Validate(json string, jsonSchema string) (*entity.Schema, error) {
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	documentLoader := gojsonschema.NewStringLoader(json)

	schemaEntity := &entity.Schema{}
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		schemaEntity.ValidateResult = result
		return schemaEntity, err
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		invalidJSONError := errors.New("Invalid JSON Data agains schema")
		schemaEntity.ValidateResult = result
		return schemaEntity, invalidJSONError
	}
	schemaEntity.ValidateResult = result

	return schemaEntity, nil
}
