package handler

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tidwall/sjson"
	"github.com/yalp/jsonpath"
)

type PostResource struct {
	URN string `json:"urn"`
}

var addMetaDataMgmtHeader = map[string]string{
	"Source":          config.Viper.GetString("application.source.api.metadatamgmt.add"),
	"Success.Code":    "0",
	"Success.Message": "Success",
	"Failure.Code":    "-1",
	"Failure.Message": "Failure",
}

var addMetaDataMgmtErrorCode = map[error]string{
	entity.ErrMissingRequiredField: "40102",
	entity.ErrItemNotFound:         "40106",
	entity.ErrInvalidJSON:          "40105",
	entity.ErrDatabaseFailure:      "40502",
	entity.ErrDefault:              "40503",
}

var addMetaDataMgmtErrorDesc = map[error]string{
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
func AddResource(service service.MetaDataMgmtService, schemaService service.SchemaLocalService, schemaMap map[string]*entity.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.AddResource() handler...")

		// #3 - Invoke service for usecase execution
		var item *entity.Resource

		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrInvalidJSON, ""))
			return
		}
		// Step 1: Accept json value and parameter
		params := mux.Vars(r)
		catalog := params["catalog"]
		contentType := params["contentType"]
		urn := fmt.Sprintf("urn:resource:%s:%s", catalog, contentType)
		if urn == "" {
			w.WriteHeader(http.StatusOK)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrMissingRequiredField, "URN is required"))
			return
		}

		jsonSchema := schemaMap[urn]
		if jsonSchema == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrItemNotFound, "Json Schema does not existed"))
			return
		}

		var resource entity.Resource
		err = json.Unmarshal(bodyByte, &resource)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrInvalidJSON, err.Error()))
			return
		}
		var resourceJSONMap interface{}
		json.Unmarshal(bodyByte, &resourceJSONMap)

		bodyByte, err = sjson.SetBytes(bodyByte, "urn", urn)
		//Step 2: Add/Update ID
		if resource.ID == "" {
			resourceID := uuid.New().String()
			bodyByte, err = sjson.SetBytes(bodyByte, "id", resourceID)
		}
		bodyByte = IdentityAdding(bodyByte)
		//Step 3: Update Name Value
		contentName, err := BuildResourceName(resource.URN, resourceJSONMap)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrDatabaseFailure, err.Error()))
			return
		}

		bodyByte, err = sjson.SetBytes(bodyByte, "name", contentName)

		//Step 4: Validate
		result, err := schemaService.Validate(string(bodyByte), jsonSchema.Data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrInvalidJSON,
				fmt.Sprintf("Please double check json: %v", result.ValidateResult.Errors())))
			return
		}
		// Re-parse resource for extract latest values such as name
		err = json.Unmarshal(bodyByte, &resource)
		resource.Body = string(bodyByte)
		if err != nil {
			w.WriteHeader(http.StatusNotModified)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrDatabaseFailure, err.Error()))
			return
		}

		cbKey, err := BuildCouchbaseKey(resource.URN, contentName, resourceJSONMap)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrDatabaseFailure, err.Error()))
			return
		}
		//Step 4: Check if ID is set or not
		resource.Key = cbKey

		// Find by db key and id if exist -> return - NO add

		item, err = service.AddResource(&resource)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildAddMetaDataMgmtFailureRespBody(entity.ErrDatabaseFailure, err.Error()))
			return
		}

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildAddMetaDataMgmtSuccessRespBody(item))
		return
	})
}

func buildAddMetaDataMgmtSuccessRespBody(item *entity.Resource) []byte {
	logger.Logger.Debug("Entering handler.buildAddMetaDataMgmtSuccessRespBody() ...")
	logger.Logger.Debug("Count: ", item)

	res := response.ResponseResponse{}
	res.Header = &response.Header{
		Source:     addMetaDataMgmtHeader["Source"],
		Code:       addMetaDataMgmtHeader["Success.Code"],
		Message:    addMetaDataMgmtHeader["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}
	var data interface{}
	json.Unmarshal([]byte(item.Body), &data)
	res.Data = &data
	resStr, _ := json.Marshal(res)
	return resStr
}

func buildAddMetaDataMgmtFailureRespBody(validationError error, msg string) []byte {
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

func BuildCouchbaseKey(urn string, name string, resourceJSONMap interface{}) (string, error) {
	if urn == "" {
		return "", entity.ErrInvalidInputAttr2
	}

	providerNameObj, err := jsonpath.Read(resourceJSONMap, "$.providerName")
	if err != nil {
		return "", entity.ErrInvalidInputAttr2
	}
	providerName := fmt.Sprintf("%v", providerNameObj)

	ss := strings.Split(urn, ":")
	resourceType := ss[len(ss)-1]
	return fmt.Sprintf("%s_%s_%s", providerName, resourceType, name), nil

}
func BuildResourceName(urn string, resourseJSONMap interface{}) (string, error) {
	if urn == "" {
		return "", entity.ErrInvalidInputAttr1
	}

	// https://quickplay.atlassian.net/browse/PE-429
	// “name” of the content should be the title name in case of movie or standalone content.
	if strings.HasSuffix(urn, "tvepisode") {
		// Name of the content should be <seriesName>_ <season_num>_ <epside_num> for tvepisodes
		seriesNameObj, err := jsonpath.Read(resourseJSONMap, "$.show_details.series_title")
		if err != nil {
			return "", err
		}
		seriesName := strings.Join(strings.Fields(fmt.Sprintf("%v", seriesNameObj)), "")

		seasonNoObj, err := jsonpath.Read(resourseJSONMap, "$.show_details.season_number")
		if err != nil {
			return "", err
		}
		seasonNo := fmt.Sprintf("%v", seasonNoObj)

		episodeNoObj, err := jsonpath.Read(resourseJSONMap, "$.show_details.episode_number")
		if err != nil {
			return "", err
		}
		episodeNo := fmt.Sprintf("%v", episodeNoObj)

		return strings.ToLower(fmt.Sprintf("%s_s%s_e%s", seriesName, seasonNo, episodeNo)), nil
	} else if strings.HasSuffix(urn, "tvseason") {
		// Name of the content should be <seriesName>_ <season_num>_ <epside_num> for tvepisodes
		seriesNameObj, err := jsonpath.Read(resourseJSONMap, "$.show_details.series_title")
		if err != nil {
			return "", err
		}
		seriesName := strings.Join(strings.Fields(fmt.Sprintf("%v", seriesNameObj)), "")

		seasonNoObj, err := jsonpath.Read(resourseJSONMap, "$.show_details.season_number")
		if err != nil {
			return "", err
		}
		seasonNo := fmt.Sprintf("%v", seasonNoObj)

		return strings.ToLower(fmt.Sprintf("%s_%s", seriesName, seasonNo)), nil
	} else if strings.HasSuffix(urn, "linearchannel") {
		//For live channel name would be just channel name and for live schedules name
		//field would be <channel name>_ <epoch start time>
		channelNameObj, err := jsonpath.Read(resourseJSONMap, "$.name")
		if err != nil {
			return "", err
		}
		channelName := strings.Join(strings.Fields(fmt.Sprintf("%v", channelNameObj)), "")

		return strings.ToLower(fmt.Sprintf("%s_%d", channelName, time.Now().Unix())), nil
	} else {
		titleNameObj, err := jsonpath.Read(resourseJSONMap, "$.titleMetadata.title")
		if err != nil {
			return "", err
		}

		//Remove all space
		titleName := strings.Join(strings.Fields(fmt.Sprintf("%v", titleNameObj)), "")
		return strings.ToLower(titleName), nil
	}

}

func generateHashFromMap(m map[string]interface{}) string {
	bytes, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}

	return generateHashFromString(string(bytes))
}

func generateHashFromString(s string) string {

	h := sha1.New()

	h.Write([]byte(s))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

type JsonIdModifier func(map[string]interface{})

func AttributeModifier(valMap map[string]interface{}) {
	if valMap["id"] == nil && valMap["urn"] != nil && valMap["urn"] != "" {
		valMap["id"] = generateHashFromMap(valMap)
	}
}

// IdentityAdding traveral to all json node then add id attribute if it is not existed and urn is found
func IdentityAdding(input []byte) []byte {
	// Creating the maps for JSON
	m := map[string]interface{}{}

	// Parsing/Unmarshalling JSON encoding/json
	err := json.Unmarshal(input, &m)

	if err != nil {
		panic(err)
	}

	parseMap(m, AttributeModifier)

	bytes, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}

	return bytes
}
func parseArray(anArray []interface{}, fn JsonIdModifier) {
	for _, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			valMap := val.(map[string]interface{})
			parseMap(valMap, fn)
			fn(valMap)
		case []interface{}:
			parseArray(val.([]interface{}), fn)
		default:
			break
		}
	}
}

func parseMap(aMap map[string]interface{}, fn JsonIdModifier) {
	for _, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:

			valMap := val.(map[string]interface{})
			parseMap(valMap, fn)
			fn(valMap)
		case []interface{}:
			parseArray(val.([]interface{}), fn)
		default:
			break
		}
	}
}
