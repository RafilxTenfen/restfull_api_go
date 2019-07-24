# Hash API by Tokens
Simple rest api for hash SHA3-256.

## How to run the application 

1. Clone the application with the following command - 
    ```shell
        git clone https://github.com/RafilxTenfen/restfull_api_go.git
    ``` 

2. There are number of dependencies which need to be imported before running the application. Please get the dependenices through the following commands -

    ```shell
        go get "golang.org/x/crypto/sha3"
        go get "github.com/gorilla/mux"
    ```

3. To run the application, please use the following command -

    ```shell
        go run main.go
    ```
> Note: By default the port number its being run on is **8080**.

## Hash Structure

1. Hash - Token encrypted
2. Token - Text to be encrypted
3. Created - Time that the Hash was created

```code
    Hash struct {
        Token   string    `json:"token,omitempty"`
        Hash    string    `json:"hash"`
        Created time.Time `json:"created_at"`
    }
```

## Endpoints Description

### Get All Hashes

```JSON
    URL - http://localhost:8080/hashes
    Method - GET
```

### Get Hash By ID

1. The ID is the token encrypted
2. The Result is a Hash Structure

```JSON
    URL - http://localhost:8080/hashes/:id
    Method - GET
```

### Create Entry
1. Token is a simple text to be encrypted
 
```JSON
    URL - http://localhost:8080/hash
    Method - POST
    Body - (content-type = application/json)
    {
        "token": "Text to be encrypted"
    }
```
2. The result will be similar as following - 

```JSON
    {
        "hash":"a37f2f5b614918c29b2d89a75810568f4926febd91be04190680fb0d9d52bb49",
        "created_at":"2019-07-21 23:57:26.431475 -0300 UTC"
    }
```

## Unit Test Description

To run all the unit test cases, please do the following -

1. `go run main.go`
2. Open a new terminal
3. `go test -v`

## Hope everything works. Thank you.