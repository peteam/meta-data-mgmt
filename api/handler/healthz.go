package handler

import (
	"net/http"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"
)

func Healthz(service service.MetaDataMgmtService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if service.Healthz() != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	})
}
