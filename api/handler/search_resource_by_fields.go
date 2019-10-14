package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
)

var searchByFieldsHeader = map[string]string{
	"Source":          config.Viper.GetString("application.source.api.metadatamgmt.searchByFields"),
	"Success.Code":    "0",
	"Success.Message": "Success",
	"Failure.Code":    "-1",
	"Failure.Message": "Failure",
}

var searchByFieldsErrorCode = map[error]string{
	entity.ErrMissingRequiredField: "40102",
	entity.ErrItemNotFound:         "40106",
	entity.ErrInvalidJSON:          "40105",
	entity.ErrDatabaseFailure:      "40502",
	entity.ErrDefault:              "40503",
}

var searchByFieldsErrorDesc = map[error]string{
	entity.ErrMissingRequiredField: config.Viper.GetString(
		"service.error.description.missingRequiredField"),
	entity.ErrDatabaseFailure: config.Viper.GetString(
		"service.error.description.databaseFailure"),
	entity.ErrInvalidJSON: config.Viper.GetString(
		"service.error.description.invalidJSON"),
	entity.ErrItemNotFound: config.Viper.GetString(
		"service.error.description.itemNotFound"),
	entity.ErrDefault: config.Viper.GetString(
		"service.error.description.default"),
}

/* AddResource add resource into Couchbase then return
 *
 */
func SearchResourceByFields(service service.MetaDataMgmtService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.searchResourceByFields() handler...")

		// #3 - Invoke service for usecase execution
		var searchEntity *entity.SearchByFieldsBodyEntity

		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildSearchResourceByFieldsFailureRespBody(entity.ErrInvalidJSON, ""))
			return
		}

		err = json.Unmarshal(bodyByte, &searchEntity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildSearchResourceByFieldsFailureRespBody(entity.ErrInvalidJSON, err.Error()))
			return
		}

		if searchEntity.Match != entity.SearchFieldAll && searchEntity.Match != entity.SearchFieldAny {

			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildSearchResourceByFieldsFailureRespBody(entity.ErrMissingRequiredField, "Match is ALL or ANY only"))
			return

		}
		// Step 1: Accept json value and parameter
		params := mux.Vars(r)
		//"/resource/lookup/{catalog}/{contentType}",
		catalog := params["catalog"]
		contentType := params["contentType"]
		urn := fmt.Sprintf("urn:resource:%s:%s", catalog, contentType)
		if urn == "" {
			w.WriteHeader(http.StatusOK)
			w.Write(buildSearchResourceByFieldsFailureRespBody(entity.ErrMissingRequiredField, "URN is required"))
			return
		}
		searchEntity.URN = urn

		pageSizeStr := params["pageSize"]
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			pageSize = 10
		}

		pageNumberStr := params["pageNumber"]
		pageNumber, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			pageNumber = 1
		}

		err = json.Unmarshal(bodyByte, &searchEntity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildSearchResourceByFieldsFailureRespBody(entity.ErrInvalidJSON, err.Error()))
			return
		}
		searchResult, err := service.SearchResourceByFields(pageSize, pageNumber, searchEntity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(buildSearchResourceByFieldsFailureRespBody(entity.ErrDatabaseFailure, err.Error()))
			return
		}
		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildSearchResourceByFieldsSuccessRespBody(searchResult))
		return
	})
}

func buildSearchResourceByFieldsSuccessRespBody(result gocb.SearchResults) []byte {
	logger.Logger.Debug("Entering handler.buildAddMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("Count: ", result.TotalHits())

	res := response.SearchResponse{}
	res.Header = &response.Header{
		Source:     addMetaDataMgmtHeader["Source"],
		Code:       addMetaDataMgmtHeader["Success.Code"],
		Message:    addMetaDataMgmtHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}
	// var data []byte
	// allHit := make([]gocb.SearchResultHit, 0)
	var searchResults []string
	for _, hit := range result.Hits() {
		searchResults = append(searchResults, hit.Id)
	}

	var searchFacets []*gocb.SearchResultFacet

	for _, row := range result.Facets() {
		searchFacets = append(searchFacets, &row)
	}
	res.SearchResponseData = &response.SearchResponseData{
		SearchResults: searchResults,
		SearchFacets:  searchFacets,
	}

	resStr, _ := json.Marshal(res)
	return resStr
}

func buildSearchResourceByFieldsFailureRespBody(validationError error, msg string) []byte {
	logger.Logger.Debug("Entering handler.buildAddDataMgmtFailureRespBody() ...")

	var errors []response.Error
	if msg != "" {
		msg = " - " + msg
	}
	errors = append(errors,
		response.Error{
			Code:        addMetaDataMgmtErrorCode[validationError],
			Description: addMetaDataMgmtErrorDesc[validationError] + msg,
		})

	res := response.Response{}
	res.Header = &response.Header{
		Source:     addMetaDataMgmtHeader["Source"],
		Code:       addMetaDataMgmtHeader["Failure.Code"],
		Message:    addMetaDataMgmtHeader["Failure.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Errors:     errors,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}
