package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"

	"gopkg.in/go-playground/validator.v9"
)

var schemaMetaDataMgmtHeader = map[string]string{
	"Source":          config.Viper.GetString("application.source.api.metadatamgmt.schema"),
	"Success.Code":    "0",
	"Success.Message": "Success",
	"Failure.Code":    "-1",
	"Failure.Message": "Failure",
}

var schemaMetaDataMgmtErrorCode = map[error]string{
	entity.ErrInvalidInputResourceType: "40102",
	entity.ErrDatabaseFailure:          "40502",
	entity.ErrDefault:                  "40503",
}

var schemaMetaDataMgmtErrorDesc = map[error]string{
	entity.ErrInvalidInputResourceType: "Invalid Request URI",
	entity.ErrDatabaseFailure:          "Subsystem failure",
	entity.ErrDefault:                  "Unknown failure",
}

/*RetrieveSchema returns the schema of Resources
 *
 */
func RetrieveSchema(service service.MetaDataMgmtService, schemaMap map[string]*entity.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.RetrieveSchema() handler...")

		// #1 - Parse path parameters & request headers
		var resource string
		var catalogType string
		var resourceType string
		vars := mux.Vars(r)
		if vars != nil {

			catalogType = vars["catalog"]
			resource = vars["resourceType"]
			resourceType = "urn:resource:" + catalogType + ":" + resource
		}
		var validURNs []string
		for _, subres := range schemaMap {

			validURNs = append(validURNs, subres.URN)
		}
		fmt.Print(validURNs[0])

		// #2 - Validate input
		errors := validateSchemaResourceType(resourceType, validURNs)
		if errors != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildSchemaMetaDataMgmtFailureRespBody(entity.ErrInvalidInputResourceType))
			return
		}

		// #3 - Fetch the Schema
		var smap = schemaMap[resourceType].Data

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildSchemaMetaDataMgmtSuccessRespBody(smap))
		return
	})
}

func validateSchemaResourceType(resourceType string, validURNs []string) error {
	logger.Logger.Debug("Entering handler.validateSchemaResourceType() ...")

	v := validator.New()
	var inputError error

	var err = v.Var(resourceType, "required")
	var res = false

	if err == nil {
		for _, subres := range validURNs {
			logger.Logger.Info("Resource Type Parsed:" + resourceType + "Compare against schema resourceType:" + subres)
			if strings.Compare(resourceType, subres) == 0 {
				logger.Logger.Debug("Resource Matched :" + subres)
				res = true
				break
			}

		}

	}

	if err != nil || res == false {
		inputError = entity.ErrInvalidInputResourceType
		logger.Logger.Debug(schemaMetaDataMgmtErrorCode[entity.ErrInvalidInputResourceType] + " - " + schemaMetaDataMgmtErrorDesc[entity.ErrInvalidInputResourceType])

	}
	return inputError
}

func buildSchemaMetaDataMgmtSuccessRespBody(schema string) []byte {
	logger.Logger.Debug("Entering handler.buildSchemaMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("Schema: ", schema)

	res := response.ContentResponse{}
	res.Header = &response.Header{
		Source:     schemaMetaDataMgmtHeader["Source"],
		Code:       schemaMetaDataMgmtHeader["Success.Code"],
		Message:    schemaMetaDataMgmtHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(schema), &result)

	res.Data = result

	resStr, _ := json.Marshal(res)

	return resStr
}

func buildSchemaMetaDataMgmtFailureRespBody(validationError error) []byte {
	logger.Logger.Debug("Entering handler.buildSchemaDataMgmtFailureRespBody() ...")

	var errors []response.Error

	errors = append(errors,
		response.Error{
			Code:        schemaMetaDataMgmtErrorCode[validationError],
			Description: schemaMetaDataMgmtErrorDesc[validationError],
		})

	res := response.Response{}
	res.Header = &response.Header{
		Source:     schemaMetaDataMgmtHeader["Source"],
		Code:       schemaMetaDataMgmtHeader["Failure.Code"],
		Message:    schemaMetaDataMgmtHeader["Failure.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Errors:     errors,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}
