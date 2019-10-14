package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/middleware"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/response"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestListResource_Success(t *testing.T) {

	// Initialize
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(ListResource(mockSer))
	r := mux.NewRouter()
	r.Handle("/list/resourceTypes", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	listOfResource := []*entity.ResourceType{
		{
			ContentType: "catalog",
			URN:         "movie",
		},
		{
			ContentType: "catalog",
			URN:         "tvSeries",
		},
	}

	mockSer.EXPECT().ListResource().Return(listOfResource, nil)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/list/resourceTypes", nil)
	req.Header.Set("X-Authorization", "test")
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var body *response.MultiResponse
	json.NewDecoder(resp.Body).Decode(&body)
	require.Equal(t, "MW-DATASERVICE-HTTP-03", body.Header.Source)
	require.Equal(t, "0", body.Header.Code)
	require.Equal(t, "Success", body.Header.Message)
	require.Equal(t, "Success", body.Header.Message)
	require.Equal(t, "movie", body.ResourceData.ResItem[0].URN)
	for _, val := range body.ResourceData.ResItem {
		fmt.Println("val", val.URN)
	}
	// Clean up
	resp.Body.Close()
	ts.Close()
}

func TestListResource_FailDB(t *testing.T) {

	// Initialize
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(ListResource(mockSer))
	r := mux.NewRouter()
	r.Handle("/list/resourceTypes", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	listOfResource := []*entity.ResourceType{
		{
			ContentType: "catalog",
			URN:         "movie",
		},
		{
			ContentType: "catalog",
			URN:         "tvSeries",
		},
	}

	mockSer.EXPECT().ListResource().Return(listOfResource, entity.ErrDatabaseFailure)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/list/resourceTypes", nil)
	req.Header.Set("X-Authorization", "test")
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var body *response.Response
	json.NewDecoder(resp.Body).Decode(&body)
	require.Equal(t, "MW-DATASERVICE-HTTP-03", body.Header.Source)
	require.Equal(t, "-1", body.Header.Code)
	require.Equal(t, "Failure", body.Header.Message)
	require.Equal(t, "40502", body.Header.Errors[0].Code)

	// Clean up
	resp.Body.Close()
	ts.Close()
}
