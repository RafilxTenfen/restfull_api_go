package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/karalabe/go-ethereum/crypto/sha3"
	"golang.org/x/crypto/sha3"
)

type Hash struct {
	Token     string `json:"token"`
	Hash      string `json:"hash"`
	CreatedAt string `json:"createdAt"`
}

var Hashes []Hash

func main() {
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/hashes", allHashes)
	myRouter.HandleFunc("/hashes/{id}", getHash)
	myRouter.HandleFunc("/hash", createHash).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func allHashes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: allHashes")
	json.NewEncoder(w).Encode(Hashes)
}

func getHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, hash := range Hashes {
		if hash.Hash == key {
			json.NewEncoder(w).Encode(hash)
		}
	}

	fmt.Println("Endpoint Hit: getHash")
}

func createHash(w http.ResponseWriter, r *http.Request) {
	var resHash Hash
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(body), &resHash)

	for _, hash := range Hashes {
		if hash.Token == resHash.Token {
			strResponse := "This Token Was Already Added"
			json.NewEncoder(w).Encode(strResponse)
			return
		}
	}
	//sum := sha256.Sum256([]byte("Texto a ser cifrado"))
	h := make([]byte, 32)
	buf := []byte("Texto a ser cifrado")
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(h, buf)
	fmt.Printf("%x\n", h)
	// expected 369ee900da8fd705ea41965e3df5df6cb7cc87a682bd29cf6c5c99253e9f87d5
	//sha3.ShakeSum256(h, buf)
	//fmt.Printf("%x", h)
	//resHash.Hash = h
	json.NewEncoder(w).Encode(resHash)
}

func decodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}
