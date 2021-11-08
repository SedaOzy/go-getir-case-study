package handlers

import (
	"net/http"
	"sync"

	"github.com/SedaOzy/go-getir-case-study/configuration"
)

const (
	InMemoryHandlerRouteUrl = "/in-memory"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var lock = sync.RWMutex{}
var mapForInMemo map[string]string = make(map[string]string)

// In Memory request will handle two verb:
// Get method handler verifies query string id and checks whether key exists in map.
// Post method handler verifies body message and add item into map.
func InMemoryHandler(config *configuration.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if len(id) == 0 {
				ErrorResponseHandlerf(&w, http.StatusBadRequest, "id parameter must be provided!")
				return
			}

			keyValuePair := read(id)
			if len(keyValuePair.Value) == 0 {
				ErrorResponseHandlerf(&w, http.StatusNotFound, "Value not found!")
				return
			}

			PrepareResponse(&w, http.StatusOK, &keyValuePair)

		case http.MethodPost:
			var keyValuePair KeyValuePair
			isSuccessful := GetRequestBody(r, &w, &keyValuePair)
			if !isSuccessful {
				return
			}

			if len(keyValuePair.Key) == 0 {
				ErrorResponseHandlerf(&w, http.StatusBadRequest, "Key should not be empty!")
				return
			}

			if len(keyValuePair.Value) == 0 {
				ErrorResponseHandlerf(&w, http.StatusBadRequest, "Value should not be empty!")
				return
			}

			write(&keyValuePair)
			PrepareResponse(&w, http.StatusOK, &keyValuePair)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Read key from map.
func read(key string) KeyValuePair {
	lock.RLock()
	defer lock.RUnlock()

	return KeyValuePair{
		Key:   key,
		Value: mapForInMemo[key],
	}
}

// Write data into map using lock
func write(item *KeyValuePair) {
	lock.Lock()
	defer lock.Unlock()
	mapForInMemo[item.Key] = item.Value
}
