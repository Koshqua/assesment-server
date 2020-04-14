package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//User ...
type User struct {
	ID       string `json:"_id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	LastName string `json:"last_name"  bson:"last_name"`
	Country  string `json:"country" bson:"country"`
	City     string `json:"city" bson:"city"`
	Gender   string `json:"gender" bson:"gender"`
	//mongo-go-driver doesn't support custom types at the moment
	//and so far I didn't find  an approach for parsing custom layout of date in json without
	//creating custom type with unmarshalJSON method, so it will satisfy Unmarshaller interface
	//so, I decided to leave the date as a plain string
	BirthDate string `json:"birth_date" bson:"birth_date"`
}

type UsersHandler struct {
	collection *mongo.Collection
}
type CreateUserHandler struct {
	collection *mongo.Collection
}
type ShowUserHandler struct {
	collection *mongo.Collection
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Disconnect(context.TODO())
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	collection := client.Database("assesment").Collection("users")
	fmt.Println("Db is connected")
	router := httprouter.New()
	router.Handler("GET", "/users", UsersHandler{collection})
	router.Handler("POST", "/users", CreateUserHandler{collection})
	router.Handler("GET", "/users/:id", ShowUserHandler{collection})
	log.Fatalln(http.ListenAndServe(":3000", router))
}

func (uh UsersHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	limKey := req.URL.Query().Get("limit")
	pageKey := req.URL.Query().Get("page")
	if limKey == "" || pageKey == "" {
		limKey = "50"
		pageKey = "1"
	}
	limit, err := strconv.Atoi(limKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetSkip(int64(limit * page)).SetLimit(int64(limit))
	c, err := uh.collection.Find(context.Background(), bson.M{}, opts)
	var users []User
	err = c.All(context.TODO(), &users)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(users[0].BirthDate)
	j, err := json.Marshal(users)
	if err != nil {
		http.Error(res, fmt.Sprintf("Couldnt parse json %s", err), http.StatusBadRequest)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(j)
}

func (cuh CreateUserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get("Content-Type")
	var user User
	if ct == "application/json" {
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, fmt.Sprintf("Couldn't read req body with error %s", err), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bs, &user)
		if err != nil {
			http.Error(res, fmt.Sprintf("Couldn't parse json with error %s", err), http.StatusBadRequest)
			return
		}
	}
	var u User
	err := cuh.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u)
	if u.LastName == "" {
		doc := bson.M{"email": user.Email, "last_name": user.LastName, "country": user.Country, "city": user.City, "gender": user.Gender, "birth_date": user.BirthDate}
		insID, err := cuh.collection.InsertOne(context.TODO(), doc)
		if err != nil {
			http.Error(res, "Error occured while inserting user", http.StatusBadRequest)
		}
		res.WriteHeader(200)
		res.Write([]byte(fmt.Sprintf("User %s was added to db", insID)))
		return
	} else {
		res.Write([]byte("User is already existing"))
	}
	if err != nil {
		http.Error(res, fmt.Sprintf("Error ocured while looking for user %s", err), http.StatusBadRequest)
	}
}

func (suh ShowUserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())
	id := params.ByName("id")
	fmt.Println(id)
	objID, _ := primitive.ObjectIDFromHex(id)
	var u User
	err := suh.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&u)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	j, err := json.Marshal(u)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(j)
}
