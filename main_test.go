package main

import (
  "testing"
  // "fmt"
  // "os"
  // "encoding/json"
  "github.com/stretchr/testify/assert"
  "net/http"
  "time"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/mongo"
  "context"
  "bytes"
  "net/http/httptest"
  "github.com/gorilla/mux"
)

func Router() *mux.Router {
    ctx, _:=context.WithTimeout(context.Background(),10*time.Second) //timeout with context
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //connecting at this port
  	client, _ = mongo.Connect(ctx, clientOptions)              //connection establishing
    collection := client.Database("GolangDevelopment").Collection("people")
    router := mux.NewRouter()
    router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
    router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
    router.HandleFunc("/person/611ec3570fee2615f4e0846b", GetPersonEndpoint(collection)).Methods("GET")
    router.HandleFunc("/person/61210acb695fb08fde578219", DeletePersonEndpoint(collection)).Methods("DELETE")
    // router.HandleFunc("/person/611fe584a1e4cc636c7e9836", UpdatePersonEndpoint(collection)).Methods("PUT")
    return router
}
func TestGetPeopleEndpoint(t *testing.T){
  request, _ := http.NewRequest("GET", "http://localhost:12345/people", nil)
  response := httptest.NewRecorder()
  Router().ServeHTTP(response, request)
  assert.Equal(t, 200, response.Code, "OK response is expected")
}
func TestGetPersonEndpoint(t *testing.T){
  request, _ := http.NewRequest("GET", "http://localhost:12345/person/611ec3570fee2615f4e0846b", nil)
  response := httptest.NewRecorder()
  Router().ServeHTTP(response, request)
  assert.Equal(t, 200, response.Code, "OK response is expected")
}
func TestDeletePersonEndpoint(t *testing.T){
  request, _ := http.NewRequest("DELETE", "http://localhost:12345/person/61210acb695fb08fde578219", nil)
  response := httptest.NewRecorder()
  Router().ServeHTTP(response, request)
  assert.Equal(t, 200, response.Code, "OK response is expected")
}
func TestUpdatePersonEndpoint(t *testing.T){
  ctx, _:=context.WithTimeout(context.Background(),10*time.Second) //timeout with context
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //connecting at this port
  client, _ = mongo.Connect(ctx, clientOptions)              //connection establishing
  collection := client.Database("GolangDevelopment").Collection("people")
  var jsonStr = []byte(`{
          "id": "611fe584a1e4cc636c7e9836",
          "firstname": "philo",
          "description": "programmer"
      }`)
  request, _ := http.NewRequest("PUT", "http://localhost:12345/person/{id}", bytes.NewBuffer(jsonStr))
  response := httptest.NewRecorder()
  handler := http.HandlerFunc(UpdatePersonEndpoint(collection))
	handler.ServeHTTP(response, request)
  // Router().ServeHTTP(response, request)
  assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestCreatePersonEndpoint(t *testing.T) {
    var jsonStr = []byte(`{
            "firstname": "maj",
            "description": "programmer",
            "address": "tirunelveli",
            "birth_date": "2005/11/10"
        }`)
    request, _ := http.NewRequest("POST", "http://localhost:12345/person", bytes.NewBuffer(jsonStr))
    response := httptest.NewRecorder()
    // handler := http.HandlerFunc(CreatePersonEndpoint)
	  // handler.ServeHTTP(response, request)
    // fmt.Println("here")
    Router().ServeHTTP(response, request)
    assert.Equal(t, 200, response.Code, "OK response is expected")
    // fmt.Println(response.Body)
    // if status := response.Code; status != http.StatusOK {
    // 		t.Errorf("handler returned wrong status code: got %v want %v",
    // 			status, http.StatusOK)
    // 	}
    // 	expected := `{
    //         "firstname": "maj",
    //         "description": "programmer",
    //         "address": "tirunelveli",
    //         "birth_date": "2005-11-10T00:00:00Z",
    //         "created_at": "2021-08-19T20:47:19.294Z"
    //     }`
    // 	if response.Body.String() != expected {
    // 		t.Errorf("handler returned unexpected body: got %v want %v",
    // 			response.Body.String(), expected)
    // 	}
}

// func TestCreatePersonEndpoint(t *testing.T) {
//
// 	var jsonStr = []byte(`{
//         "_id": "611ec3570fee2615f4e0846b",
//         "firstname": "kishore",
//         "description": "programmer",
//         "address": "tirunelveli",
//         "birth_date": "2005/11/10",
//         "created_at": "2021-08-19T20:47:19.294Z"
//     }`)
//
// 	req, err := http.NewRequest("POST", "http://localhost:12345/person", bytes.NewBuffer(jsonStr))
// 	if err != nil {
// 		// t.Fatal(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(CreatePersonEndpoint)
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// 	expected := `{
//         "_id": "611ec3570fee2615f4e0846b",
//         "firstname": "kishore",
//         "description": "programmer",
//         "address": "tirunelveli",
//         "birth_date": "2005-11-10T00:00:00Z",
//         "created_at": "2021-08-19T20:47:19.294Z"
//     }`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
