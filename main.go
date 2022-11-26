package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID int `json:"id"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email	string `json:"email"`
	Address Address `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type Address struct {
	Street string `json:"street"`
	City string `json:"city"`
	PostalCode string `json:"postalCode"`
	Country string `json:"country"`
}

var users = []User{}

func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	newUser.ID = 1;
	if(len(users) > 0) {
		newUser.ID = users[len(users) - 1].ID + 1;
	}

	w.Header().Set("Content-Type"," application/json")
	users = append(users, newUser)

	w.WriteHeader(http.StatusCreated)
	
	json.NewEncoder(w).Encode(newUser)

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type"," application/json")
	
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to int"))
		return
	}

	w.Header().Set("Content-Type"," application/json")

	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(404)
		w.Write([]byte("No record found for a user with the specified ID"))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	userId, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to int"))
		return
	}
	var updatedUser User

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data to update {firstname, lastname}")
	}

	json.Unmarshal(reqBody, &updatedUser)

	w.Header().Set("Content-Type"," application/json")

	for i, user := range users {
		if user.ID == userId {
			user.FirstName = updatedUser.FirstName
			user.LastName = updatedUser.LastName
			users = append(users[:i], user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(404)
		w.Write([]byte("No record found for a user with the specified ID"))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	userId, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to int"))
		return
	}
	for i, user := range users {
		if user.ID == userId {
			users = append(users[:i], users[i+1:]...)
			fmt.Fprintf(w, "The user with ID %v has been deleted successfully", userId)
			return
		}
	}

	w.WriteHeader(404)
		w.Write([]byte("No record found for a user with the specified ID"))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", addUser).Methods("POST")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	
	http.ListenAndServe(":5000", router)
}