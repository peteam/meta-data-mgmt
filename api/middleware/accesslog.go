package middleware

import (
	"net/http"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.AccessLogger.Info(r)
		next.ServeHTTP(w, r)
	})
}