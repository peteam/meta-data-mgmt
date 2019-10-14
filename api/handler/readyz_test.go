package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/middleware"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestReadyZ_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)
	handler := middleware.ParseHeader(Readyz(mockSer))
	r := mux.NewRouter()
	r.Handle("/readyz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Readyz().Return(nil)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/readyz", nil)
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestReadyZ_ReadyZFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(Readyz(mockSer))
	r := mux.NewRouter()
	r.Handle("/readyz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Readyz().Return(entity.ErrReadyzFailure)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/readyz", nil)
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestReadyZ_HealthZFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(Readyz(mockSer))
	r := mux.NewRouter()
	r.Handle("/readyz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Readyz().Return(entity.ErrHealthzFailure)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/readyz", nil)
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

//Generic test case to handle any panic / unknown error thrown from handler
func TestReadyZ_GeneralFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(Readyz(mockSer))
	r := mux.NewRouter()
	r.Handle("/readyz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Readyz().Return(errors.New("General error"))

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/readyz", nil)

	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}