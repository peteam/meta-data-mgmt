package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"github.com/yalp/jsonpath"
)

type SchemaLocalRepository struct {
}

func NewSchemaRepository() *SchemaLocalRepository {
	logger.BootstrapLogger.Debug("Entering Repository.SchemaLocalRepository() ...")

	return &SchemaLocalRepository{}
}

// LoadSchemas - loading json schema from local disk
func (sr *SchemaLocalRepository) LoadSchemas(schemaLocation string) (map[string]*entity.Schema, error) {
	//logger.BootstrapLogger.Infof("Loading schema from %s", schemaLocation)

	files, err := ioutil.ReadDir(schemaLocation)
	if err != nil {
		//logger.BootstrapLogger.Error(err)
		log.Fatal(err)
	}

	schemaMap := make(map[string]*entity.Schema)

	for _, f := range files {
		fileAbsPath := filepath.Join(schemaLocation, f.Name()) // relative path
		bytes, err := ioutil.ReadFile(fileAbsPath)
		if err != nil {
			return nil, err
		}

		var data interface{}
		json.Unmarshal(bytes, &data)

		urn, err := jsonpath.Read(data, "$.properties.urn.const")
		if err != nil {
			return nil, err
		}
		urnStr := fmt.Sprintf("%v", urn)
		entitySchema := entity.Schema{URN: urnStr, Location: fileAbsPath, Data: string(bytes)}
		schemaMap[urnStr] = &entitySchema
	}
	return schemaMap, nil
}
