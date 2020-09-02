package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/gorm"
)

//Account holds structural data for users
type Account struct {
	gorm.Model
	Username       string
	HashedPassword string
}

//Start kicks off the api server
//The Api server controls CRUD
//operation for important game
//models.
func Start() {
	db = databaseInit()
	log.Println("API Server started")
	//Handles
	http.HandleFunc("/api", createAccount)
}
func createAccount(w http.ResponseWriter, r *http.Request) {
	reqBodyJSON, _ := ioutil.ReadAll(r.Body)
	var reqBody = Account{}
	err := json.Unmarshal(reqBodyJSON, &reqBody)
	if err != nil {
		log.Println("Create account json unmarshal", err)
		fmt.Fprintln(w, err)
	}
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
