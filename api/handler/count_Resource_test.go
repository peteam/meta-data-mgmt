package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/middleware"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCountResource_Success(t *testing.T) {

	// Initialize
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	schemaRepo := repository.NewSchemaRepository()
	schemaService := service.NewSchemaService(schemaRepo)

	schemaMap, _ := schemaService.ReloadAllSchemas("../../conf/schemas")

	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(CountResource(mockSer, schemaMap))
	r := mux.NewRouter()
	r.Handle("/count/urn/resource/{catalog}/{resourceType}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	var count int
	mockSer.EXPECT().CountResource(gomock.Any()).Return(count, nil)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/count/urn/resource/catalog/tvseries", nil)
	req.Header.Set("X-Authorization", "test")
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var body *response.Response
	json.NewDecoder(resp.Body).Decode(&body)
	require.Equal(t, "MW-DATASERVICE-HTTP-04", body.Header.Source)
	require.Equal(t, "0", body.Header.Code)
	require.Equal(t, "Success", body.Header.Message)
	require.Equal(t, &count, body.Data.Count)

	// Clean up
	resp.Body.Close()
	ts.Close()
}

func TestCountResource_FailDB(t *testing.T) {

	// Initialize
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	schemaRepo := repository.NewSchemaRepository()
	schemaService := service.NewSchemaService(schemaRepo)

	schemaMap, _ := schemaService.ReloadAllSchemas("../../conf/schemas")

	handler := middleware.ParseHeader(CountResource(mockSer, schemaMap))
	r := mux.NewRouter()
	r.Handle("/count/urn/resource/{catalog}/{resourceType}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	var count int
	mockSer.EXPECT().CountResource(gomock.Any()).Return(count, entity.ErrDatabaseFailure)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/count/urn/resource/catalog/tvseries", nil)
	req.Header.Set("X-Authorization", "test")
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var body *response.Response
	json.NewDecoder(resp.Body).Decode(&body)
	require.Equal(t, "MW-DATASERVICE-HTTP-04", body.Header.Source)
	require.Equal(t, "-1", body.Header.Code)
	require.Equal(t, "Failure", body.Header.Message)
	require.Equal(t, "40502", body.Header.Errors[0].Code)

	// Clean up
	resp.Body.Close()
	ts.Close()
}
