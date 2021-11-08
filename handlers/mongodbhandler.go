package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SedaOzy/go-getir-case-study/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoDbHandlerRouteUrl = "/"
)

type MongoDbFilterRequest struct {
	// StartsDate will contain the date in a “YYYY-MM-DD” format.
	StartDate string `json:"startDate"`
	// StartsDate will contain the date in a “YYYY-MM-DD” format.
	EndDate string `json:"endDate"`
	// Sum of the “count” array in the documents should be between “minCount” and “maxCount”
	MinCount int `json:"minCount"`
	// Sum of the “count” array in the documents should be between “minCount” and “maxCount”
	MaxCount int `json:"maxCount"`
}

type Record struct {
	Key        string `json:"key"`
	CreatedAt  string `json:"createdAt"`
	TotalCount int    `json:"totalCount"`
}

type Records struct {
	Code    int      `json:"code"`
	Message string   `json:"msg"`
	Records []Record `json:"records"`
}

var client mongo.Client

//  Creates client and queries collection
func GetRecordsMongoDbClient(config *configuration.Config, filter *MongoDbFilterRequest) ([]string, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDb.Url).SetConnectTimeout(config.MongoDb.ConnectTimeoutMS))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	database := client.Database(config.MongoDb.TableName)
	// collection := database.Collection(config.MongoDb.CollectionName)
	database.Aggregate()

	// Create a date object using time.Parse()
	lteDate, err := time.Parse(time.RFC3339, filter.StartDate)
	if err != nil {
		log.Fatal(err)
	}

	gteDate, err := time.Parse(time.RFC3339, filter.EndDate)
	if err != nil {
		log.Fatal(err)
	}

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	return nil, err
	// }

	return databases, nil
}

func MongoDbHandler(config *configuration.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var request MongoDbFilterRequest

			isSuccessful := GetRequestBody(r, &w, &request)
			if !isSuccessful {
				return
			}

			databases, err := GetRecordsMongoDbClient(config, &request)
			if err != nil {
				ErrorResponseHandler(&w, http.StatusFailedDependency, "Mongodb connection failed!", err)
				return
			}

			PrepareResponse(&w, http.StatusOK, databases)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
