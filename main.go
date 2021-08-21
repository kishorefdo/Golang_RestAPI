//Develop a RESTful API with Golang and Mongodb NoSQL Database
//two dependencies 1.mux 2.mongodb go sdk
//use mongodb on docker container
package main

import (
    "context"
  	"encoding/json"
  	"fmt"
  	"log"
    "io/ioutil"
    "strings"
  	"net/http"
  	"time"
  	"github.com/gorilla/mux"
  	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
  	"go.mongodb.org/mongo-driver/bson/primitive"
  	"go.mongodb.org/mongo-driver/mongo"
    // "gopkg.in/juju/charmstore.v5/internal/v5"
)

//Person data Structure with JSON and BSON annotations
//data model definition
type Person struct{
  ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //mongodb's own objecctid, bson is official annotation
  Name string `json:"firstname,omitempty" bson:"firstname,omitempty"`
  Description  string `json:"description,omitempty" bson:"description,omitempty"`
  Address   string `json:"address,omitempty" bson:"address,omitempty"`
  BirthDate time.Time `json:"birth_date" bson:"birth_date"`
  CreatedAt time.Time     `json:"created_at" time_format:"sql_datetime" time_location:"UTC" bson:"created_at" time_format:"sql_datetime" time_location:"UTC"`
}
func (p *Person) Parse(s string) error {
	fields := map[string]string{}

	dec := json.NewDecoder(strings.NewReader(s))
	if err := dec.Decode(&fields); err != nil {
		return fmt.Errorf("decode person: %v", err)
	}

	p.Name = fields["firstname"]
  p.Description = fields["description"]
  p.Address = fields["address"]

	born, err := time.Parse("2006/01/02", fields["birth_date"])
	if err != nil {
		return fmt.Errorf("invalid date: %v", err)
	}
	p.BirthDate = born

	return nil
}

var client *mongo.Client
var person []Person

//API endpoints for HTTP interaction
//endpoint to insert data with POST method
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
  //recieve client data
  bodyBytes, err := ioutil.ReadAll(request.Body)
  if err != nil {
        log.Fatal(err)
    }
  bodyString := string(bodyBytes)
  response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)  //json payload to object
  person.CreatedAt = time.Now().UTC()
  if err := person.Parse(bodyString); err != nil {
		log.Fatalf("parse person: %v", err)
	}
  //database name is GolangDevelopment and collection is people
  ctx, _:=context.WithTimeout(context.Background(),10*time.Second)
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //connecting at this port
	client, _ = mongo.Connect(ctx, clientOptions)
	collection := client.Database("GolangDevelopment").Collection("people")
  // fmt.Println("here in main")
	// ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)   //inserting
	json.NewEncoder(response).Encode(result)     //returning result which is an objectID
  // response.WriteHeader(200)
  // response.Write([]byte("Item Created"))
  // respondWithJSON(response, http.StatusOK, person)
}

//returns all documents within our collection
//cursors in Mongodb used as in RDBMS
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json") //response in json format
	var people []Person
  ctx, _:=context.WithTimeout(context.Background(),40*time.Second)
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //connecting at this port
	client, _ = mongo.Connect(ctx, clientOptions)
	collection := client.Database("GolangDevelopment").Collection("people")
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second) //context definition
	cursor, err := collection.Find(ctx, bson.M{})  //return everything
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
  //decoding each iteration and adding it to a slice of the Person data type
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)  //error handling
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

//obtain particular document from the Database
func GetPersonEndpoint(collection *mongo.Collection) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		vars := mux.Vars(request)
		id := vars["id"]
		objectID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.D{{"_id", objectID}}
		var person Person
    //find one function used to filter on our id
		err := collection.FindOne(context.TODO(), filter).Decode(&person)
		if err != nil {
			// printJSON(response, errorJSON{Statuscode: http.StatusNotFound, Err: err.Error()})
		} else {
			response.WriteHeader(http.StatusFound)
      json.NewEncoder(response).Encode(person) //when no errors, returns the person
		}
	}

}


func DeletePersonEndpoint(collection *mongo.Collection) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]
		objectID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.D{{"_id", objectID}}
		var person Person
		collection.FindOneAndDelete(context.TODO(), filter).Decode(&person)
		// printJSON(response, person)
	}
}
func UpdatePersonEndpoint(collection *mongo.Collection) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
    ctx, _:=context.WithTimeout(context.Background(),10*time.Second) //timeout with context
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //connecting at this port
  	client, _ = mongo.Connect(ctx, clientOptions)
    collection := client.Database("GolangDevelopment").Collection("people")
		response.Header().Set("content-type", "application/json")
		vars := mux.Vars(request)
		id := vars["id"]
		objectID, _ := primitive.ObjectIDFromHex(id)
    // fmt.Println(objectID)
		filter := bson.D{{"_id", objectID}}
    // fmt.Println(filter)
		var person Person
		if err := json.NewDecoder(request.Body).Decode(&person); err != nil {
			panic(err)
		}
		result := collection.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": person}, options.FindOneAndUpdate().SetReturnDocument(1))
    decoded := Person{}
		if err := result.Decode(&decoded); err != nil {
			// panic(err)
		}
    // fmt.Println("here in main")
	}

}
//API configuration and Database connection
func main(){
  fmt.Println("Starting the application...")
  ctx, _:=context.WithTimeout(context.Background(),10*time.Second) //timeout with context
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") //connecting at this port
	client, _ = mongo.Connect(ctx, clientOptions)              //connection establishing
  collection := client.Database("GolangDevelopment").Collection("people")
  router := mux.NewRouter()                           //router definition
  router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
  router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
  router.HandleFunc("/person/{id}", GetPersonEndpoint(collection)).Methods("GET")
  router.HandleFunc("/person/{id}", DeletePersonEndpoint(collection)).Methods("DELETE")
  router.HandleFunc("/person/{id}", UpdatePersonEndpoint(collection)).Methods("PUT")
  http.ListenAndServe(":12345",router)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//encode payload to json
	response, _ := json.Marshal(payload)

	// set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
