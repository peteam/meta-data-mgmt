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

var countMetaDataMgmtHeader = map[string]string{
	"Source":          config.Viper.GetString("application.source.api.metadatamgmt.count"),
	"Success.Code":    "0",
	"Success.Message": "Success",
	"Failure.Code":    "-1",
	"Failure.Message": "Failure",
}

var countMetaDataMgmtErrorCode = map[error]string{
	entity.ErrInvalidInputResourceType: "40102",
	entity.ErrDatabaseFailure:          "40502",
	entity.ErrDefault:                  "40503",
}

var countMetaDataMgmtErrorDesc = map[error]string{
	entity.ErrInvalidInputResourceType: "Invalid Request URI",
	entity.ErrDatabaseFailure:          "Subsystem failure",
	entity.ErrDefault:                  "Unknown failure",
}

/*CountResource returns the count of Resources
 *
 */
func CountResource(service service.MetaDataMgmtService, schemaMap map[string]*entity.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.CountResource() handler...")

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
		errors := validateCountResourceType(resourceType, validURNs)
		if errors != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildCountMetaDataMgmtFailureRespBody(entity.ErrInvalidInputResourceType))
			return
		}

		// #3 - Invoke service for usecase execution
		var item = 0
		item, err := service.CountResource(resourceType)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildCountMetaDataMgmtFailureRespBody(entity.ErrDatabaseFailure))
			return
		}

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildCountMetaDataMgmtSuccessRespBody(item))
		return
	})
}

func validateCountResourceType(resourceType string, validURNs []string) error {
	logger.Logger.Debug("Entering handler.validateCountResourceType() ...")

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
		logger.Logger.Debug(countMetaDataMgmtErrorCode[entity.ErrInvalidInputResourceType] + " - " + countMetaDataMgmtErrorDesc[entity.ErrInvalidInputResourceType])

	}
	return inputError
}

func buildCountMetaDataMgmtSuccessRespBody(count int) []byte {
	logger.Logger.Debug("Entering handler.buildCountMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("Count: ", count)

	res := response.Response{}
	res.Header = &response.Header{
		Source:     countMetaDataMgmtHeader["Source"],
		Code:       countMetaDataMgmtHeader["Success.Code"],
		Message:    countMetaDataMgmtHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}
	res.Data = &response.Data{
		Count: &count,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}

func buildCountMetaDataMgmtFailureRespBody(validationError error) []byte {
	logger.Logger.Debug("Entering handler.buildCountDataMgmtFailureRespBody() ...")

	var errors []response.Error

	errors = append(errors,
		response.Error{
			Code:        countMetaDataMgmtErrorCode[validationError],
			Description: countMetaDataMgmtErrorDesc[validationError],
		})

	res := response.Response{}
	res.Header = &response.Header{
		Source:     countMetaDataMgmtHeader["Source"],
		Code:       countMetaDataMgmtHeader["Failure.Code"],
		Message:    countMetaDataMgmtHeader["Failure.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Errors:     errors,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}
