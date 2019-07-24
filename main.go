package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/sha3"
)

// TimeFormat used as Time.format()
const TimeFormat = "2006-01-02 15:04:05.000000 Z0700 UTC"

// Hash Simple Structure
type Hash struct {
	Token   string    `json:"token,omitempty"`
	Hash    string    `json:"hash"`
	Created time.Time `json:"created_at"`
}

// MarshalJSON format Time of Hash ex.:2019-07-21 23:22:24.664425 -0300 UTC
func (u *Hash) MarshalJSON() ([]byte, error) {
	type Alias Hash
	return json.Marshal(&struct {
		*Alias
		Created string `json:"created_at"`
	}{
		Alias:   (*Alias)(u),
		Created: u.Created.Format(TimeFormat),
	})
}

// Hashes simulate database
var Hashes []Hash

// populate data
func init() {
	time3, _ := time.Parse(TimeFormat, "2012-10-31 16:13:58.292387 +0000 UTC")
	time4, _ := time.Parse(TimeFormat, "2012-10-31 16:13:58.292387 +0000 UTC")
	Hashes = []Hash{
		{"Texto a ser cifrado 3", "aae889e4b2d49258c278632345044426982d8adcd701443ab9b8ede1f23a9032", time3},
		{"Texto a ser cifrado 4", "062e47e7d4475416ceb204edc882d07e615a1f656487a4666ae08c3c706a1eb1", time4},
	}
}

func main() {
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	hashRequests(myRouter)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func hashRequests(myRouter *mux.Router) {
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/hashes", allHashes).Methods("GET")
	myRouter.HandleFunc("/hashes/{id}", getHash).Methods("GET")
	myRouter.HandleFunc("/hash", createHash).Methods("POST")
}

func allHashes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: allHashes")
	if len(Hashes) > 0 {
		json.NewEncoder(w).Encode(Hashes)
	}
}

func getHash(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getHash")
	vars := mux.Vars(r)
	key := vars["id"]

	for _, hash := range Hashes {
		if hash.Hash == key {
			json.NewEncoder(w).Encode(hash)
			return
		}
	}

	http.Error(w, "There is no Hash with id: "+key, http.StatusNoContent) // Hash Not Found
}

func createHash(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createHash")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var sendHash Hash
	json.Unmarshal([]byte(body), &sendHash)

	if HashAlreadyAdded(sendHash.Token) {
		http.Error(w, "Token Already Added", http.StatusConflict)
		return
	}

	populateHash(&sendHash)
	addHash(sendHash)

	json.NewEncoder(w).Encode(
		&Hash{Token: "", Hash: sendHash.Hash, Created: sendHash.Created},
	)
}

func addHash(hash Hash) {
	Hashes = append(Hashes, hash)
}

func getHashToken(token string) string {
	hash := sha3.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// HashAlreadyAdded verify if a new token is in
func HashAlreadyAdded(token string) bool {
	for _, hash := range Hashes {
		if hash.Token == token {
			return true
		}
	}
	return false
}

func populateHash(hash *Hash) {
	hash.Hash = getHashToken(hash.Token)
	hash.Created = time.Now()
}
