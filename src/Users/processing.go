package Users

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/gocql/gocql"
	"net/http"
	"strconv"
)

// FormToUser -- fills a User struct with submitted form data
// params:
// r - request reader to fetch form data or url params (unused here)
// returns:
// User struct if successful
// array of strings of errors if any occur during processing
func FormToUser(r *http.Request) (User, []string) {
	var user User
	var errStr, ageStr string
	var errs []string
	var err error

	user.FirstName, errStr = processFormField(r, "firstname")
	errs = appendError(errs, errStr)
	user.LastName, errStr = processFormField(r, "lastname")
	errs = appendError(errs, errStr)
	user.Email, errStr = processFormField(r, "email")
	errs = appendError(errs, errStr)
	user.City, errStr = processFormField(r, "city")
	errs = appendError(errs, errStr)

	ageStr, errStr = processFormField(r, "age")
	if len(errStr) != 0 {
		errs = append(errs, errStr)
	} else {
		user.Age, err = strconv.Atoi(ageStr)
		if err != nil {
			errs = append(errs, "Parameter 'age' not an integer")
		}
	}
	return user, errs
}

func appendError(errs []string, errStr string) []string {
	if len(errStr) > 0 {
		errs = append(errs, errStr)
	}
	return errs
}

func processFormField(r *http.Request, field string) (string, string) {
	fieldData := r.PostFormValue(field)
	if len(fieldData) == 0 {
		return "", "Missing '" + field + "' parameter, cannot continue"
	}
	return fieldData, ""
}

func Enrich2(uuids []gocql.UUID) map[string]*User2 {
	if len(uuids) > 0 {
		users := map[string]*User2{}

		m := map[string]interface{}{}
		iterable := Cassandra.Session.Query("SELECT id, username, name FROM users WHERE id IN ?", uuids).Iter()
		for iterable.MapScan(m) {
			user_id := m["id"].(gocql.UUID)
			users[user_id.String()] = &User2{
				ID:       user_id,
				Username: m["username"].(string),
				Name:     m["name"].(string),
			}
			m = map[string]interface{}{}
		}

		return users
	}
	return map[string]*User2{}
}
