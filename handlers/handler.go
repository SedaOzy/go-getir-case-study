package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SedaOzy/go-getir-case-study/configuration"
	customlogger "github.com/SedaOzy/go-getir-case-study/customloggers"
)

// Generates the router used in the HTTP Server
func InitRouter(config *configuration.Config) *http.ServeMux {
	// Create router and define routes and return that router
	router := http.NewServeMux()

	router.HandleFunc(MongoDbHandlerRouteUrl, MongoDbHandler(config))
	router.HandleFunc(InMemoryHandlerRouteUrl, InMemoryHandler(config))

	return router
}

func Run(config *configuration.Config, router *http.ServeMux) {
	// Define server options
	server := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  config.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Server.Timeout.Write * time.Second,
		IdleTimeout:  config.Server.Timeout.Idle * time.Second,
	}

	// Alert the user that the server is starting
	customlogger.Infof("Server is starting on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			// Normal interrupt operation, ignore
		} else {
			customlogger.Errorf("Server failed to start due to err: %v", err)
		}
	}

	customlogger.Infof("Server is shutting down")
}

func ErrorResponseHandlerf(w *http.ResponseWriter, code int, message string) {
	customlogger.Errorf(message)
	http.Error(*w, message, code)
}

func ErrorResponseHandler(w *http.ResponseWriter, code int, message string, err error) {
	customlogger.Error(err, message)
	http.Error(*w, message, code)
}

func GetRequestBody(r *http.Request, w *http.ResponseWriter, v interface{}) bool {
	// Read body
	if r.Body == nil {
		ErrorResponseHandlerf(w, http.StatusBadRequest, "Please send a request body")
		return false
	}

	// body, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if err != nil {
	// 	ErrorResponseHandler(&w, http.StatusUnprocessableEntity, "Body has not been read!", err)
	// 	return
	// }

	// if err := json.Unmarshal(body, &keyValuePair); err != nil {
	// 	ErrorResponseHandler(&w, http.StatusUnprocessableEntity, "Body has not been deserialized to KeyValuePair!", err)
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		ErrorResponseHandler(w, http.StatusUnprocessableEntity, "Body has not been deserialized", err)
		return false
	}

	return true
}

func PrepareResponse(w *http.ResponseWriter, code int, v interface{}) {
	if err := json.NewEncoder(*w).Encode(v); err != nil {
		ErrorResponseHandler(w, http.StatusUnprocessableEntity, "Body has not been serialized!", err)
		return
	}

	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(code)
}
