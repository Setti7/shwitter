package Messages

import (
    "encoding/json"
    "fmt"
    "github.com/Setti7/shwitter/Cassandra"
    "github.com/Setti7/shwitter/Users"
    "github.com/gocql/gocql"
    "github.com/gorilla/mux"
    "net/http"
)

// Get -- handles GET request to /messages/ to fetch all messages
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params (unused here)
func Get(w http.ResponseWriter, r *http.Request) {
    var messageList []Message
    var enrichedMessages []Message
    var userList []gocql.UUID
    m := map[string]interface{}{}

    fmt.Println("Fetching activities from Database")
    query := "SELECT id,user_id,message FROM messages"
    iterable := Cassandra.Session.Query(query).Iter()
    for iterable.MapScan(m) {
        userID := m["user_id"].(gocql.UUID)
        messageList = append(messageList, Message{
            ID:      m["id"].(gocql.UUID),
            UserID:  userID,
            Message: m["message"].(string),
        })
        userList = append(userList, userID)
        m = map[string]interface{}{}
    }

    names := Users.Enrich(userList)
    for _, message := range messageList {
        message.UserFullName = names[message.UserID.String()]
        enrichedMessages = append(enrichedMessages, message)
    }
    fmt.Println("message list after enrichment", enrichedMessages)

    json.NewEncoder(w).Encode(AllMessagesResponse{Messages: enrichedMessages})
}


// GetOne -- handles GET request to /messages/{message_uuid} to fetch one message
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func GetOne(w http.ResponseWriter, r *http.Request) {
    var message Message
    var errs []string
    var found bool = false

    vars := mux.Vars(r)
    id := vars["uuid"]

    uuid, err := gocql.ParseUUID(id)
    if err != nil {
        errs = append(errs, err.Error())
    } else {
        m := map[string]interface{}{}
        query := "SELECT id,user_id,message FROM messages WHERE id=? LIMIT 1"
        iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
        for iterable.MapScan(m) {
            found = true
            userID := m["user_id"].(gocql.UUID)
            names := Users.Enrich([]gocql.UUID{userID})
            fmt.Println("names", names)
            message = Message{
                ID:           m["id"].(gocql.UUID),
                UserID:       userID,
                UserFullName: names[userID.String()],
                Message:      m["message"].(string),
            }
        }
        if !found {
            errs = append(errs, "Message not found")
        }
    }

    if found {
        json.NewEncoder(w).Encode(GetMessageResponse{Message: message})
    } else {
        json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
    }
}