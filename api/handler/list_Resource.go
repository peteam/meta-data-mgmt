package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
)

var listMetaDataMgmtHeader = map[string]string{
	"Source":          config.Viper.GetString("application.source.api.metadatamgmt.list"),
	"Success.Code":    "0",
	"Success.Message": "Success",
	"Failure.Code":    "-1",
	"Failure.Message": "Failure",
}

var listMetaDataMgmtErrorCode = map[error]string{
	entity.ErrInvalidInputResourceType: "40102",
	entity.ErrDatabaseFailure:          "40502",
	entity.ErrDefault:                  "40503",
}

var listMetaDataMgmtErrorDesc = map[error]string{
	entity.ErrInvalidInputResourceType: "Invalid Request URI",
	entity.ErrDatabaseFailure:          "Subsystem failure",
	entity.ErrDefault:                  "Unknown failure",
}

/*ListResource returns the count of Resources
 *
 */
func ListResource(service service.MetaDataMgmtService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.listResource() handler...")

		// #3 - Invoke service for usecase execution
		var item []*entity.ResourceType
		item, err := service.ListResource()
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildListMetaDataMgmtFailureRespBody(entity.ErrDatabaseFailure))
			return
		}

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildListMetaDataMgmtSuccessRespBody(item))
		return
	})
}

func buildListMetaDataMgmtSuccessRespBody(items []*entity.ResourceType) []byte {
	logger.Logger.Debug("Entering handler.buildListMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("List: ", items)

	res := response.MultiResponse{}
	res.Header = &response.Header{
		Source:     listMetaDataMgmtHeader["Source"],
		Code:       listMetaDataMgmtHeader["Success.Code"],
		Message:    listMetaDataMgmtHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}
	var data []response.Data
	for _, item := range items {
		data = append(data,
			response.Data{
				ResourceType: item.ContentType,
				URN:          item.URN})
	}

	res.ResourceData = &response.ResourceData{
		ResItem: data,
	}

	resStr, _ := json.Marshal(res)
	return resStr
}

func buildListMetaDataMgmtFailureRespBody(validationError error) []byte {
	logger.Logger.Debug("Entering handler.buildListDataMgmtFailureRespBody() ...")

	var errors []response.Error

	errors = append(errors,
		response.Error{
			Code:        listMetaDataMgmtErrorCode[validationError],
			Description: listMetaDataMgmtErrorDesc[validationError],
		})

	res := response.Response{}
	res.Header = &response.Header{
		Source:     listMetaDataMgmtHeader["Source"],
		Code:       listMetaDataMgmtHeader["Failure.Code"],
		Message:    listMetaDataMgmtHeader["Failure.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Errors:     errors,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}
