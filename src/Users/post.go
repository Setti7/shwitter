package Users

import (
    "encoding/json"
    "fmt"
    "github.com/Setti7/shwitter/Cassandra"
    "github.com/gocql/gocql"
    "net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
    var errs []string
    var gocqlUuid gocql.UUID
    // FormToUser() is included in Users/processing.go
    // we will describe this later
    user, errs := FormToUser(r)
    // have we created a user correctly
    var created bool = false
    // if we had no errors from FormToUser, we will
    // attempt to save our data to Cassandra
    if len(errs) == 0 {
        fmt.Println("creating a new user")
        // generate a unique UUID for this user
        gocqlUuid = gocql.TimeUUID()
        // write data to Cassandra
        if err := Cassandra.Session.Query(`
      INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
            gocqlUuid, user.FirstName, user.LastName, user.Email, user.City, user.Age).Exec(); err != nil {
            errs = append(errs, err.Error())
        } else {
            created = true
        }
    }
    // depending on whether we created the user, return the
    // resource ID in a JSON payload, or return our errors
    if created {
        fmt.Println("user_id", gocqlUuid)
        json.NewEncoder(w).Encode(NewUserResponse{ID: gocqlUuid})
    } else {
        fmt.Println("errors", errs)
        json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
    }
}
