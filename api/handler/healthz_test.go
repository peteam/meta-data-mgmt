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

func TestHealthz_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(Healthz(mockSer))
	r := mux.NewRouter()
	r.Handle("/healthz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Healthz().Return(nil)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/healthz", nil)
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestHealthz_HealthzFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(Healthz(mockSer))
	r := mux.NewRouter()
	r.Handle("/healthz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Healthz().Return(entity.ErrHealthzFailure)

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/healthz", nil)
	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

//Generic test case to handle any panic / unknown error thrown from handler
func TestHealthz_GeneralFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSer := service.NewMockMetaDataMgmtService(mockCtrl)

	handler := middleware.ParseHeader(Healthz(mockSer))
	r := mux.NewRouter()
	r.Handle("/healthz", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	//Mocking
	mockSer.EXPECT().Healthz().Return(errors.New("General error"))

	// Execute request
	req, _ := http.NewRequest("GET", ts.URL+"/healthz", nil)

	resp, _ := http.DefaultClient.Do(req)

	// Asserts
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}
