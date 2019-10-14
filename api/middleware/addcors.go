package middleware

import (
	"net/http"
	"strings"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
		 * 	Add service specific CORS response headers here.
		 *	The below is an reference.
		 */
		 if !strings.Contains(r.URL.Path, "healthz") &&  !strings.Contains(r.URL.Path, "readyz") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Authorization, X-Authorization, Content-Type")
			w.Header().Set("Content-Type", "application/json")
		}
		if r.Method == "OPTIONS" {
			return
		}

		/*
		 *	Bring all path specific response headers logic here (if required).
		 * 	For example:
		 *	if strings.Contains(r.URL.Path, "swagger") {
		 *		w.Header().Set("Content-Type", "text/plain")
		 *	}
		 */

		next.ServeHTTP(w, r)
	})
}