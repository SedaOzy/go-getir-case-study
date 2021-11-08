package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SedaOzy/go-getir-case-study/configuration"
	"github.com/SedaOzy/go-getir-case-study/handlers"
)

var Router http.ServeMux

func Init() {
	cfg, isSuccessful := configuration.Init()
	if !isSuccessful {
		return
	}

	Router = *handlers.InitRouter(cfg)
}

func TestMain(m *testing.M) {
	Init()
	code := m.Run()
	os.Exit(code)
}

// Verifies that if request does not contain query string id.
func TestInMemoryIfRequestNotContainQueryString(t *testing.T) {
	req, _ := http.NewRequest("GET", handlers.InMemoryHandlerRouteUrl, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// Verifies that if key does not exist.
func TestInMemoryIfKeyDoesNotExist(t *testing.T) {
	requestUrl := fmt.Sprintf("%s?id=%s", handlers.InMemoryHandlerRouteUrl, "sample")
	req, _ := http.NewRequest("GET", requestUrl, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Verifies that if Post request is successfull.
func TestInMemoryIfPostRequestSuccessful(t *testing.T) {
	keyValuePair := handlers.KeyValuePair{
		Key:   "sample",
		Value: "deneme",
	}
	jsonBytes, _ := json.Marshal(keyValuePair)
	req, _ := http.NewRequest("POST", handlers.InMemoryHandlerRouteUrl, bytes.NewBuffer(jsonBytes))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	requestUrl := fmt.Sprintf("%s?id=%s", handlers.InMemoryHandlerRouteUrl, keyValuePair.Key)
	reqGet, _ := http.NewRequest("GET", requestUrl, nil)
	responseGet := executeRequest(reqGet)

	checkResponseCode(t, http.StatusOK, responseGet.Code)
}

// Verifies that if body is empty.
func TestInMemoryIfBodyIsEmpty(t *testing.T) {
	req, _ := http.NewRequest("POST", handlers.InMemoryHandlerRouteUrl, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// Verifies that if body is not valid.
func TestInMemoryIfBodyIsNotValid(t *testing.T) {
	keyValuePair := handlers.KeyValuePair{}
	jsonBytes, _ := json.Marshal(keyValuePair)
	req, _ := http.NewRequest("POST", handlers.InMemoryHandlerRouteUrl, bytes.NewBuffer(jsonBytes))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestMongoDbIfPostRequestSuccessful(t *testing.T) {
	filter := handlers.MongoDbFilterRequest{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}

	jsonBytes, _ := json.Marshal(filter)
	req, _ := http.NewRequest("POST", handlers.MongoDbHandlerRouteUrl, bytes.NewBuffer(jsonBytes))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
