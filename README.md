### Golang Sample service
The following features are showcased in this Quickstart sample service
- HTTP Muxing
- Chaining of Middleware Handlers
- Logger
- Config Reader
- Couchbase integration
- Segregation of layers (Handler -> Service -> Repository)

This sample service is built using **Clean Architecture** pattern by **Robert C. Martin** as reference.
Reference: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

### Illustration of the flows
- API client -> HTTP Server -> HTTP Muxer -> Middleware handlers -> Business handlers -> Service layer -> Repository layer
- API client -> HTTP Server -> HTTP Muxer -> Liveness handler
- API client -> HTTP Server -> HTTP Muxer -> Readiness handler

### HTTP Muxing
Initialize an Router
````
    r := mux.NewRouter()
````
Add Liveness handler,
````
	r.Handle("/readyz",
		handler.Readyz(
			service)).Methods("GET")
````
Add Readiness handler,
````
	r.Handle("/healthz",
		handler.Healthz(
			service)).Methods("GET")
````
Add business handlers custom to the service. The sample below is to HTTP PUT: Insert something,
````
	 r.Handle("/foo", 
		handler.AddFoo(
			service)).Methods("PUT", "OPTIONS")
````
### Chaining of Handlers
If the service requires to perform some operations before the flow reaches the business handlers, the middleware handlers can be chained as below.
````
	 r.Handle("/foo", 
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.AddFoo(
					service)))).Methods("PUT", "OPTIONS")
````
AccessLog - This is a custom handler to log all incoming requests
Parseheader - This is a custom handler to perform validations on Header attributes. For example: Bearer token validation

Any number of handlers can be chained together as per service requirement.

### Sample middleware handler #1
The below is an custom handler to log all incoming requests,
````
package middleware

import (
	"net/http"
	"cicd-github.quickplay.com/org-name-here/repo-name-here/pkg/logger"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.AccessLogger.Info(r)
		next.ServeHTTP(w, r)
	})
}
````
### Sample middleware handler #2
The below is an custom handler to parse one (or more) header parameters, parse them and assign to context that can be read again in business handlers,
````
package middleware

import (
	"net/http"
	"context"
	"cicd-github.quickplay.com/org-name-here/repo-name-here/pkg/logger"
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
		if (len(attr1) > 0) {
			logger.Logger.Debug("Attr1 found in header: " + attr1)
			ctx := context.WithValue(r.Context(), "Attribute1", attr1)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
````

### CORS handler
The below is an custom handler to add CORS headers to all outgoing requests,
````
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
````

### Liveness Handler
````
func Healthz(service service.FooService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*
		 *	Implement service specific implementation for liveness probe here.
		 */
		w.WriteHeader(http.StatusOK)
		return
	})
}
````

### Readiness Handler
````
func Readyz(service service.FooService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*
		 *	Implement service specific implementation for readiness probe here.
		 * 	Include all checks for the service dependency components here.
		 *	For example: Couchbase DB repo initialization.
		 */
		w.WriteHeader(http.StatusOK)
		return
	})
}
````

### Sample Business handler
The below is an sample business handler which performs the following steps,
- Parse the incoming request body
- Validate the request body
- Invoke the service layer to execute business logic
- Build and return the success or failure response

````
package handler

import (
	"net/http"
	"time"
	"encoding/json"

	"cicd-github.quickplay.com/org-name-here/repo-name-here/pkg/service"
	"cicd-github.quickplay.com/org-name-here/repo-name-here/pkg/entity"
	"cicd-github.quickplay.com/org-name-here/repo-name-here/api/response"
	"cicd-github.quickplay.com/org-name-here/repo-name-here/pkg/logger"

	"gopkg.in/go-playground/validator.v9"
)

/*
	Define Struct for the request body 
	for field validation 
*/
type FooRequest struct {
	Attr1		string		`json:"attr1" validate:"required"`
	Attr2		string		`json:"attr2" validate:"required"`
}

var header = map[string]string {
	"Source": "XY-FOOSERVICE-HTTP-01",
	"Success.Code": "0",
	"Success.Message": "Success",
	"Failure.Code": "-1",
	"Failure.Message": "Failure",
}

var code = map[error]string {
	// Define error codes of AddFoo API
	entity.ErrInvalidInputAttr1: "40101",
	entity.ErrInvalidInputAttr2: "40102",
	entity.ErrDatabaseFailure: "40103",
	entity.ErrDefault: "40104",
}

var desc = map[error]string {
	// Define error descriptions of AddFoo API
	entity.ErrInvalidInputAttr1: "Invalid request parameter - attr1",
	entity.ErrInvalidInputAttr2: "Invalid request parameter - attr2",
	entity.ErrDatabaseFailure: "Subsystem failure",
	entity.ErrDefault: "Unknown failure",
}

/*
	This handler performs the following steps,
	1. Parse request body
	2. Validate input
	3. Invoke service to execute business
	4. Build success or failure response
*/
func AddFoo(service service.FooService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("Entering handler.AddFoo() ...")

		// #1 - Parse request body
		var input *FooRequest
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// #2 - Validate input
		attr1 := r.Context().Value("Attribute1")
		if (attr1 != nil) {
			input.Attr1 = attr1.(string)
		}
		errors := validateAddFoo(input)
		if errors != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(buildAddFooFailureRespBody(errors))
			return
		}

		// #3 - Invoke service for usecase execution
		err := service.AddFoo(
			&entity.Foo{
				Attr1: attr1.(string),
				Attr2: input.Attr2,
			})
		if err != nil {
			w.WriteHeader(http.StatusOK)
			var errors []error
			w.Write(buildAddFooFailureRespBody(append(errors, err)))
			return
		}

		// #4 - Build and return success response
		w.WriteHeader(http.StatusOK)
		w.Write(buildAddFooSuccessRespBody())
		return
	})
}

func validateAddFoo(input *FooRequest) ([]error) {
	logger.Logger.Debug("Entering handler.validateAddFoo() ...")

	v := validator.New()
	var err = v.Struct(input)
	if err != nil {
		var errors []error
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
				case "Attr1":
					errors = append(errors, entity.ErrInvalidInputAttr1)
					logger.Logger.Debug(code[entity.ErrInvalidInputAttr1] + " - " + desc[entity.ErrInvalidInputAttr1])
				case "Attr2":
					errors = append(errors, entity.ErrInvalidInputAttr2)
					logger.Logger.Debug(code[entity.ErrInvalidInputAttr2] + " - " + desc[entity.ErrInvalidInputAttr2])
			}
		}
		return errors
	}
	return nil
}

func buildAddFooSuccessRespBody() ([]byte) {
	logger.Logger.Debug("Entering handler.buildAddFooSuccessRespBody() ...")

	res := response.Response{}
	res.Header = &response.Header{
		Source: header["Source"],
		Code: header["Success.Code"],
		Message: header["Success.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
	}
	resStr, _ := json.Marshal(res)
	return resStr
}

func buildAddFooFailureRespBody(errlist []error) ([]byte) {
	logger.Logger.Debug("Entering handler.buildAddFooFailureRespBody() ...")

	var errors []response.Error
	for _, err := range errlist {
		errors = append(errors, 
			response.Error{
				Code: code[err],
				Description: desc[err],
		})
	}

	res := response.Response{}
	res.Header = &response.Header{
		Source: header["Source"],
		Code: header["Failure.Code"],
		Message: header["Failure.Message"],
		SystemTime: (time.Now().UnixNano() / 1e6),
		Errors: errors,
	}
	resStr, _ := json.Marshal(res)
	return resStr
}
````

### Define a Service layer
A service layer can constructed in 4 simple steps,
- Step 1: Define a Service Interface

````
type FooService interface {
	AddFoo(item *entity.Foo) error
	GetFoo(Attr1 string, Attr2 string) (*entity.Foo, error)
	DeleteFoo(Attr1 string, Attr2 string) error
	ListFoos(Attr1 string, pageNumber string, pageSize string,
		sortBy string, sortOrder string) ([]*entity.Foo, error)
}
````
- Step 2: Define a Service struct

````
type Service struct {
	repo repository.DbRepository
}
````
- Step 3: Add initialize function to inject dependencies

````
func NewService(r repository.DbRepository) *Service {
	logger.BootstrapLogger.Debug("Entering Service.NewService() ...")
	return &Service{
		repo: r,
	}
}
````
- Step 4: Add use cases as methods

````
func (s *Service) AddFoo(item *entity.Foo) error {
	logger.Logger.Debug("Entering Service.AddFoo() ...")
	// Include use case logic here
}

func (s *Service) GetFoo(Attr1 string, Attr2 string) (*entity.Foo, error) {
	logger.Logger.Debug("Entering Service.GetFoo() ...")
	// Include use case logic here
}

func (s *Service) DeleteFoo(Attr1 string, Attr2 string) error {
	logger.Logger.Debug("Entering Service.DeleteFoo() ...")
	// Include use case logic here
}

func (s *Service) ListFoos(Attr1 string, pageNumber string, pageSize string,
	sortBy string, sortOrder string) ([]*entity.Foo, error) {
	logger.Logger.Debug("Entering Service.ListFoos() ...")
	// Include use case logic here
}
````

### Define a Repository layer
A repository layer can constructed in 3 simple steps,
- Step 1: Define a Repository Interface

````
type DbRepository interface {
	Insert(item *entity.Foo) error
	Upsert(item *entity.Foo) error
	Retrieve(Attr1 string, Attr2 string) (*entity.Foo, error)
	Remove(Attr1 string, Attr2 string) error
	List(Attr1 string, pageNumber string, pageSize string,
		sortBy string, sortOrder string) ([]*entity.Foo, error)
	Count(Attr1 string) (int, error)
}
````
- Step 2: Define a Repository struct

````
type CbRepository struct {
	Cluster *gocb.Cluster
	Bucket  *gocb.Bucket
}
````
- Step 3: Add initialize function to connect database or inject dependencies

````
func NewCbRepository() *CbRepository {
	logger.BootstrapLogger.Debug("Entering Repository.NewCbRepository() ...")
	// Initialize the Cluster and Bucket
}
````
- Add DB operation as methods

````
func (r *CbRepository) Insert(item *entity.Foo) error {
	logger.Logger.Debug("Entering CbRepository.Insert() ...")
	// Include DB operation logic here
}

func (r *CbRepository) Upsert(item *entity.Foo) error {
	logger.Logger.Debug("Entering CbRepository.Upsert() ...")
	// Include DB operation logic here
}

func (r *CbRepository) Retrieve(Attr1 string, Attr2 string) (*entity.Foo, error) {
	logger.Logger.Debug("Entering CbRepository.Retrieve() ...")
	// Include DB operation logic here
}

func (r *CbRepository) Remove(Attr1 string, Attr2 string) error {
	logger.Logger.Debug("Entering CbRepository.Remove() ...")
	// Include DB operation logic here
}

func (r *CbRepository) List(Attr1 string, pageNumber string, pageSize string,
	sortBy string, sortOrder string) ([]*entity.Foo, error) {
	logger.Logger.Debug("Entering CbRepository.List() ...")
	// Include DB operation logic here
}

func (r *CbRepository) Count(Attr1 string) (int, error) {
	logger.Logger.Debug("Entering CbRepository.Count() ...")
	// Include DB operation logic here
}
````
### Putting all above together in main
````
func main() {
	logger.BootstrapLogger.Info("Service starting...")
	initService()
}
````

````
func initService() {
	logger.BootstrapLogger.Debug("Entering initService...")

	// Dependency Injection
	repository := initDatabase()
	service := service.NewService(repository)

	// Define the handlers
	r := mux.NewRouter()
	addHandlers(r, service)
	r.Use(middleware.Cors)

	// Configure and start the server
	logger.Logger.Fatal(initHttpServer(r).ListenAndServe())
}
````

````
func initDatabase() repository.DbRepository {
	logger.BootstrapLogger.Debug("Entering initDatabase...")
	return repository.NewCbRepository()
}
````

````
func addHandlers(r *mux.Router, service service.FooService) {

	// Add readiness probe handler
	r.Handle("/healthz",
		handler.Healthz(
			service)).Methods("GET")

	// Add liveness probe handler
	r.Handle("/readyz",
		handler.Readyz(
			service)).Methods("GET")

	 r.Handle("/foo", 
		middleware.AccessLog(
			middleware.ParseHeader(
				handler.AddFoo(
					service)))).Methods("PUT", "OPTIONS")
}
````
````
func initHttpServer(r *mux.Router) *http.Server {
	srv := &http.Server{Handler: r}
	srv.Addr = ":8080"
	srv.ReadTimeout = 15 * time.Second
	srv.WriteTimeout = time.Duration(writeTimeout) * time.Second
	srv.MaxHeaderBytes = 4096
	return srv
````
