package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/handler"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/api/middleware"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/service"

	"github.com/gorilla/mux"
)

func init() {
	// Bring project specific custom initializations here
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func main() {
	logger.BootstrapLogger.Info("Service starting...")
	initService()
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func initService() {

	defaultApplicationName := "Data Management Service"
	config.Viper.SetDefault("application.name", defaultApplicationName)
	logger.BootstrapLogger.Debug("Entering initService...")
	logger.BootstrapLogger.Info("Starting " + config.Viper.GetString("application.name") +
		" with profile=" + config.Viper.GetString("profile") + " properties")
	logger.BootstrapLogger.Debug(config.Viper.AllSettings())
	logger.BootstrapLogger.Debug(config.Viper.GetString("schema.json.location"))
	// Dependency Injection
	repository := initDatabase()
	fooservice := service.NewService(repository)

	// JSON Schema Dependency Injection
	schemaRepository := initSchemaRepository()
	schemaService := service.NewSchemaService(schemaRepository)
	//schemaService.ReloadAllSchemas(config.Viper.GetString("schema.json.location"))
	fmt.Print(schemaService)
	schemaMap := make(map[string]*entity.Schema)
	var err error
	schemaMap, err = schemaService.ReloadAllSchemas(config.Viper.GetString("schema.json.location"))
	if err != nil {
		//logger.BootstrapLogger.Error(err)
		log.Fatal(err)
	}

	// Define the handlers
	r := mux.NewRouter()
	addHandlers(r, fooservice, schemaService, schemaMap)
	r.Use(middleware.Cors)

	// Configure and start the server
	logger.Logger.Fatal(initHttpServer(r).ListenAndServe())
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func initDatabase() repository.DbRepository {
	logger.BootstrapLogger.Debug("Entering initDatabase...")

	if config.Viper.GetString(
		"database.couchbase.enable") == "true" {
		logger.BootstrapLogger.Debug("About to initialize Couchbase repo...")
		return repository.NewCbRepository()

	} else {
		// Throw panic
		logger.BootstrapLogger.Error("Incorrect database configuration settings. Can't proceed!")
		panic(entity.ErrInvalidConfig)
	}
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func initSchemaRepository() repository.SchemaRepository {
	logger.BootstrapLogger.Debug("Entering initSchemaRepository...")

	return repository.NewSchemaRepository()
}

//-- SYSTEM GENEREATED: DO NOT MODIFY ----
func initHttpServer(r *mux.Router) *http.Server {
	srv := &http.Server{Handler: r}
	srv.Addr = ":" + config.Viper.GetString("http.port")

	srv.ReadTimeout = 15 * time.Second
	readTimeoutStr := config.Viper.GetString("http.readTimeout")
	if len(readTimeoutStr) > 0 {
		readTimeout, err := strconv.Atoi(readTimeoutStr)
		if err != nil {
			logger.BootstrapLogger.Error("Invalid ReadTimeout in config-env.yml, default to 15 seconds")
		} else {
			srv.ReadTimeout = time.Duration(readTimeout) * time.Second
		}
	}
	logger.BootstrapLogger.Info("ReadTimeout seconds set to HTTP server")

	srv.WriteTimeout = 15 * time.Second
	writeTimeoutStr := config.Viper.GetString("http.writeTimeout")
	if len(writeTimeoutStr) > 0 {
		writeTimeout, err := strconv.Atoi(writeTimeoutStr)
		if err != nil {
			logger.BootstrapLogger.Error("Invalid WriteTimeout in config-env.yml, default to 15 seconds")
		} else {
			srv.WriteTimeout = time.Duration(writeTimeout) * time.Second
		}
	}
	logger.BootstrapLogger.Info("WriteTimeout seconds set to HTTP server")

	srv.MaxHeaderBytes = 4096
	maxHeaderBytesStr := config.Viper.GetString("http.maxHeaderBytes")
	if len(maxHeaderBytesStr) > 0 {
		maxHeaderBytes, err := strconv.Atoi(maxHeaderBytesStr)
		if err != nil {
			logger.BootstrapLogger.Error("Invalid MaxHeaderBytes in config-env.yml, default to 4096 bytes")
		} else {
			srv.MaxHeaderBytes = maxHeaderBytes
		}
	}
	logger.BootstrapLogger.Info("MaxHeaderBytes seconds set to HTTP server")

	return srv
}

//-- SYSTEM GENEREATED: Testing Code ----
//-- This is a primitive test handler that should be removed by the developer
func simplehandler(service service.MetaDataMgmtService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := "simple response..."
		fmt.Fprintf(w, "Hello from:  "+title+"\n")
	})
}
func addHandlers(r *mux.Router, service service.MetaDataMgmtService, schemaService service.SchemaLocalService, schemaMap map[string]*entity.Schema) {

	// --------- WARNING: DO NOT DELETE THIS HANDLER -------------
	// Add readiness probe handler
	r.Handle("/healthz",
		handler.Healthz(
			service)).Methods("GET")

	// --------- WARNING: DO NOT DELETE THIS HANDLER -------------
	// Add liveness probe handler
	r.Handle("/readyz",
		handler.Readyz(
			service)).Methods("GET")

	//-- This is a primitive test handler that should be removed by the developer
	r.Handle("/", simplehandler(service)).Methods("GET")

	//-- Retrieve count of resources for a given resource type
	// -- Sample API "http://localhost:8080/count/urn/resource/vod/movie"
	r.Handle("/count/urn/resource/{catalog}/{resourceType}",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.CountResource(
					service, schemaMap)))).Methods("GET", "OPTIONS")

	//-- list distinct resource types
	// -- Sample API "http://localhost:8080/list/resourceTypes"

	r.Handle("/list/resourceTypes",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.ListResource(
					service)))).Methods("GET", "OPTIONS")

	//-- Retrieve schema of resources for a given resource type
	// -- Sample API "http://localhost:8080/schema/urn/resource/vod/movie"
	r.Handle("/schema/urn/resource/{catalog}/{resourceType}",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.RetrieveSchema(
					service, schemaMap)))).Methods("GET", "OPTIONS")

	/*
	 * Get Content by contentType
	 */
	r.Handle("/resource/urn/resource/{catalog}/{contentType}",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.GetMultiContent(
					service, schemaMap)))).Methods("GET", "OPTIONS")

	/*
	 * Get Content by ID
	 */
	r.Handle("/resource/urn/resource/{catalog}/{contentType}/{id}",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.GetContent(
					service, schemaMap)))).Methods("GET", "OPTIONS")

	r.Handle("/resource/{catalog}/{contentType}",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.AddResource(
					service, schemaService, schemaMap)))).Methods("POST", "OPTIONS")

	r.Handle("/resource/lookup/{catalog}/{contentType}",
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.SearchResourceByFields(
					service)))).Methods("POST")

}
