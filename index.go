package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conn = connectdb()

type Todo struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title            string             `json:"title,omitempty" bson:"title,omitempty"`
	Todo_description string             `json:"todo_description,omitempty" bson:"todo_description,omitempty"`
	Status           string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedDate      time.Time          `json:"createdDate"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todo", createTodo).Methods("POST")
	r.HandleFunc("/todo", readTodo).Methods("GET")
	r.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8001", r))

}
func readTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	todos := []Todo{}
	cur, err := conn.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
		// return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var todo Todo
		err := cur.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}

		todos = append(todos, todo)
	}

	json.NewEncoder(w).Encode(todos)

}
func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo

	_ = json.NewDecoder(r.Body).Decode(&todo)
	result, err := conn.InsertOne(context.TODO(), todo)

	if err != nil {
		log.Fatal(err)
		// return
	}

	json.NewEncoder(w).Encode(result)
}
func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	todo := Todo{}

	filter := bson.M{"_id": id}

	_ = json.NewDecoder(r.Body).Decode(&todo)

	update := bson.D{
		{"$set", bson.D{
			{"title", todo.Title},
		}},
	}

	err := conn.FindOneAndUpdate(context.TODO(), filter, update).Decode(&todo)

	if err != nil {
		log.Fatal(err)
		// return
	}

	todo.Id = id

	json.NewEncoder(w).Encode(todo)
}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	deleteResult, err := conn.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err, w)
		// return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
func connectdb() *mongo.Collection {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("todo_list").Collection("todo")
	return collection
}
