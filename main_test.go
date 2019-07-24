package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAllHashes(t *testing.T) {
	req, err := http.NewRequest("GET", "/hashes", nil)
	handleReqError(err, t)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(allHashes)
	handler.ServeHTTP(rr, req)
	handleReqStatus(rr.Code, t)

	expected := `[{"token":"Texto a ser cifrado 3","hash":"aae889e4b2d49258c278632345044426982d8adcd701443ab9b8ede1f23a9032","created_at":"2012-10-31 16:13:58.292387 Z UTC"},{"token":"Texto a ser cifrado 4","hash":"062e47e7d4475416ceb204edc882d07e615a1f656487a4666ae08c3c706a1eb1","created_at":"2012-10-31 16:13:58.292387 Z UTC"}]`
	handleTestExpected(expected, rr.Body.String(), t)
}

const HashCifrado4 = "062e47e7d4475416ceb204edc882d07e615a1f656487a4666ae08c3c706a1eb1"

func TestGetHashStatus(t *testing.T) {
	handler := http.HandlerFunc(getHash)
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, "http://localhost:8080")
	e.GET("/hashes/{id}").
		WithPath("id", HashCifrado4).
		Expect().
		Status(http.StatusOK)
}

func TestGetHashBody(t *testing.T) {
	handler := http.HandlerFunc(getHash)
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, "http://localhost:8080")
	body := e.GET("/hashes/{id}").
		WithPath("id", HashCifrado4).
		Expect().
		Body()

	expected := `{"token":"Texto a ser cifrado 4","hash":"062e47e7d4475416ceb204edc882d07e615a1f656487a4666ae08c3c706a1eb1","created_at":"2012-10-31T16:13:58.292387Z"}`
	handleTestExpected(expected, body.Raw(), t)
}

func TestGetHashNoContent(t *testing.T) {
	handler := http.HandlerFunc(getHash)
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, "http://localhost:8080")
	e.GET("/hashes/{id}").
		WithPath("id", "hashErrado").
		Expect().
		Status(http.StatusNoContent)
}

func TestCreateHash(t *testing.T) {
	jsonStr := []byte(`{
		"token": "Texto a ser cifrado 2"
	}`)

	req, err := http.NewRequest("POST", "/hash", bytes.NewBuffer(jsonStr))
	handleReqError(err, t)

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createHash)
	handler.ServeHTTP(rr, req)
	handleReqStatus(rr.Code, t)

	body, err := ioutil.ReadAll(rr.Body)
	handleReqError(err, t)

	var HashTest Hash
	json.Unmarshal([]byte(body), &HashTest)

	expected := `a37f2f5b614918c29b2d89a75810568f4926febd91be04190680fb0d9d52bb49`
	handleTestExpected(expected, HashTest.Hash, t)
}

func TestCreateHashDupplicate(t *testing.T) {
	r := rand.New(rand.NewSource(99))
	strRand := fmt.Sprintf("{\"token\": \"Texto a ser cifrado %v \" }", r.Int63())
	jsonStr := []byte(strRand)

	req, _ := http.NewRequest("POST", "/hash", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createHash)
	handler.ServeHTTP(rr, req)

	// second request
	req, err := http.NewRequest("POST", "/hash", bytes.NewBuffer(jsonStr))
	handleReqError(err, t)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(createHash)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusConflict)
	}

	expected := "Token Already Added"
	handleTestExpected(expected, rr.Body.String(), t)
}

func TestGetHashToken(t *testing.T) {
	token := "Texto a ser cifrado"
	expected := "369ee900da8fd705ea41965e3df5df6cb7cc87a682bd29cf6c5c99253e9f87d5"
	handleTestExpected(getHashToken(token), expected, t)
}

func handleReqError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func handleReqStatus(code int, t *testing.T) {
	if code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			code, http.StatusOK)
	}
}

func handleTestExpected(expected string, value string, t *testing.T) {
	if !strings.EqualFold(strings.TrimSpace(expected), strings.TrimSpace(value)) {
		t.Errorf("handler returned unexpected value: got %s want %s", value, expected)
	}
}
