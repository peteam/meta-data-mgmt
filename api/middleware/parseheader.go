package middleware

import (
	"context"
	"net/http"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
)

func ParseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Parsing request headers")

		/*
			Add here your business logic to parse Header
			For examples:
			1. Parse and decrypt OVAT to read user identifer.
			2. Validate Bearer token

			Note: Add parsed information (if required) to context
			to use in downstream layers
		*/
		attr1 := r.Header.Get("X-Authorization")
		if len(attr1) > 0 {
			logger.Logger.Debug("Attr1 found in header: " + attr1)
			ctx := context.WithValue(r.Context(), "Attribute1", attr1)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
