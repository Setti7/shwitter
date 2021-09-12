// Source:
// https://getstream.io/blog/building-a-performant-api-using-go-and-cassandra/
package main

import (
    "encoding/json"
    "github.com/Setti7/shwitter/Cassandra"
    "github.com/Setti7/shwitter/Messages"
    "github.com/Setti7/shwitter/Users"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

type heartbeatResponse struct {
    Status string `json:"status"`
    Code   int    `json:"code"`
}

func main() {
    CassandraSession := Cassandra.Session
    defer CassandraSession.Close()
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", heartbeat)
    router.HandleFunc("/users/", Users.Post).Methods("POST")
    router.HandleFunc("/users/", Users.Get).Methods("GET")
    router.HandleFunc("/users/{uuid}", Users.GetOne).Methods("GET")
    router.HandleFunc("/messages/", Messages.Post).Methods("POST")
    router.HandleFunc("/messages/", Messages.Get).Methods("GET")
    router.HandleFunc("/messages/{uuid}", Messages.GetOne).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
