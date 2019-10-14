package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

var getcontentHeader = map[string]string{
	"Source":          "MW-DATASERVICE-HTTP-06", // config.Viper.GetString("application.source.api.metadatamgmt.get"),
	"Success.Code":    "0",
	"Success.Message": "Success",
	"Failure.Code":    "-1",
	"Failure.Message": "Failure",
}

var contentMetaDataMgmtErrorCode = map[error]string{
	entity.ErrInvalidInputResourceType: "40102",
	entity.ErrDatabaseFailure:          "40502",
	entity.ErrDefault:                  "40503",
	entity.ErrInvalidPageNumber:        "40504",
	entity.ErrInvalidPageSize:          "40505",
}

var contentMetaDataMgmtErrorDesc = map[error]string{
	entity.ErrInvalidInputResourceType: "Invalid Request URI",
	entity.ErrDatabaseFailure:          "Subsystem failure",
	entity.ErrDefault:                  "Unknown failure",
	entity.ErrInvalidPageNumber:        "PageNumber is lesser than default value",
	entity.ErrInvalidPageSize:          "PageSize should be greater or equal to 1 ",
}

/*GetContentOptional returns the content record of requested contentParameters
 *
 */
func GetContent(service service.MetaDataMgmtService, schemaMap map[string]*entity.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.GetMovies() handler...")
		// #1 - Parse query parameters
		vars := mux.Vars(r)
		var contentType string
		var catalog string
		var resourceType string
		resourceId := vars["id"]
		if vars != nil {

			contentType = vars["contentType"]
			catalog = vars["catalog"]
			resourceType = "urn:resource:" + catalog + ":" + contentType
		}
		var validURNs []string
		for _, subres := range schemaMap {

			validURNs = append(validURNs, subres.URN)
		}

		varsqp := r.URL.Query()
		// catalogTypeURN := vars.Get("catalogTypeURN")
		pageNumber := varsqp.Get("pageNumber")
		pageSize := varsqp.Get("pageSize")
		entityStatus := varsqp.Get("entityStatus")
		providerName := varsqp.Get("providerName")
		catalogType := varsqp.Get("catalogType")
		fmt.Println(pageNumber, pageSize, entityStatus)

		// #2 - Validate input
		errors := validateResourceType(resourceType, validURNs)
		if errors != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildContentMetaDataMgmtFailureRespBody(entity.ErrInvalidInputResourceType))
			return
		}

		// #3 - Invoke service for usecase execution
		var item *entity.Content
		item, err := service.GetContentOptionalService(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
		fmt.Printf("%+v\n", item)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildContentMetaDataMgmtFailureRespBody(entity.ErrItemNotFound))
			return
		}

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildContentMetaDataMgmtSuccessRespBody(item))
		return
	})
}

func validateResourceType(resourceType string, validURNs []string) error {
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

/*GetContent returns the content record of requested contentParameters
 *
 */
func GetMultiContent(service service.MetaDataMgmtService, schemaMap map[string]*entity.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.GetMovies() handler...")
		// #1 - Parse query parameters
		vars := mux.Vars(r)
		var contentType string
		var catalog string
		var resourceType string
		//resourceId := vars["id"]
		if vars != nil {

			contentType = vars["contentType"]
			catalog = vars["catalog"]
			resourceType = "urn:resource:" + catalog + ":" + contentType
		}
		var validURNs []string
		for _, subres := range schemaMap {

			validURNs = append(validURNs, subres.URN)
		}

		varsqp := r.URL.Query()
		// catalogTypeURN := vars.Get("catalogTypeURN")
		pageNumber := varsqp.Get("pageNumber")
		pageSize := varsqp.Get("pageSize")
		entityStatus := varsqp.Get("entityStatus")
		providerName := varsqp.Get("providerName")
		catalogType := varsqp.Get("catalogType")
		resourceId := varsqp.Get("ids")
		size, _ := strconv.Atoi(pageSize)
		pageNo, _ := strconv.Atoi(pageNumber)
		if pageNo < 1 {
			w.WriteHeader(http.StatusOK)
			w.Write(buildContentMetaDataMgmtFailureRespBody(entity.ErrInvalidPageNumber))
			return
		}
		if size <= 0 {
			w.WriteHeader(http.StatusOK)
			w.Write(buildContentMetaDataMgmtFailureRespBody(entity.ErrInvalidPageSize))
			return
		}
		var idsStringArray string
		if resourceId != "" {
			ids := strings.Split(resourceId, ",")
			idsStringArray = "['"
			for i := 0; i < len(ids); i++ {
				idsStringArray += ids[i]
				idsStringArray += "','"
			}
			newString1 := (strings.TrimSuffix(idsStringArray, "'"))
			newString2 := (strings.TrimSuffix(newString1, ","))
			newString2 += "]"
			idsStringArray = newString2
		}

		// #2 - Validate input
		errors := validateResourceType(resourceType, validURNs)
		if errors != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildContentMetaDataMgmtFailureRespBody(entity.ErrInvalidInputResourceType))
			return
		}

		// #3 - Invoke service for usecase execution
		var item *entity.MultiContent
		item, err := service.GetMultiContentOptionalService(idsStringArray, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
		totalcount, err := service.TotalContentCount(idsStringArray, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
		fmt.Printf("%+v\n", item)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildContentMetaDataMgmtFailureRespBody(entity.ErrItemNotFound))
			return
		}

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildContentMetaDataMgmtSuccessMultiRespBody(item, size, totalcount, pageNo))
		return
	})
}

func buildContentMetaDataMgmtSuccessRespBody(item *entity.Content) []byte {
	logger.Logger.Debug("Entering handler.buildCountMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("Count: ", item)

	res := response.ContentResponse{}
	res.Header = &response.Header{
		Source:     getcontentHeader["Source"],
		Code:       getcontentHeader["Success.Code"],
		Message:    getcontentHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}
	res.Data = item.DocResponse
	resStr, _ := json.Marshal(res)
	return resStr
}

func buildContentMetaDataMgmtSuccessMultiRespBody(item *entity.MultiContent, size int, count int, pageNo int) []byte {
	logger.Logger.Debug("Entering handler.buildCountMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("Count: ", item)
	//size, _ := strconv.Atoi(pageSize)
	//pageNo, _ := strconv.Atoi(pageNumber)
	start := (pageNo - 1) * size
	if size == 0 {
		size = count
	}
	res := response.ContentResponse{}
	res.Header = &response.Header{
		Source:     getcontentHeader["Source"],
		Code:       getcontentHeader["Success.Code"],
		Message:    getcontentHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Start:      &start,
		Rows:       &size,
		Count:      &count,
	}
	res.Data = item.DocResponse
	resStr, _ := json.Marshal(res)
	return resStr
}

func buildContentMetaDataMgmtFailureRespBody(validationError error) []byte {
	logger.Logger.Debug("Entering handler.buildCountDataMgmtFailureRespBody() ...")

	var errors []response.Error

	errors = append(errors,
		response.Error{
			Code:        contentMetaDataMgmtErrorCode[validationError],
			Description: contentMetaDataMgmtErrorDesc[validationError],
		})

	res := response.ContentResponse{}
	res.Header = &response.Header{
		Source:     getcontentHeader["Source"],
		Code:       getcontentHeader["Failure.Code"],
		Message:    getcontentHeader["Failure.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Errors:     errors,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}
