package counter

import (
	"fmt"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// The string value for the enum MUST be the same as the cassandra table for
// the counter.
// The counter table MUST have its ID column called "id" and its counter
// column called "count".
type counterTable string

type CounterTable interface {
	Increment(ID string, value int) error
	GetValue(ID string) (count int, err error)
}

// Increment a cassandra counter.
//
// Returns ErrInvalidID for invalid IDs and ErrUnexpected for any other errors.
func (c counterTable) Increment(ID string, value int) error {
	if ID == "" {
		return errors.ErrInvalidID
	}

	q := fmt.Sprintf("UPDATE %s SET count = count + ? WHERE id = ?", c)
	err := service.Cassandra().Query(q, value, ID).Exec()

	if err != nil {
		log.LogError("counter.Increment", "Could not increment the counter", err)
		return errors.ErrUnexpected
	} else {
		return nil
	}
}

// Get the value for a cassandra counter.
//
// Returns ErrInvalidID for invalid IDs and ErrUnexpected for any other errors.
func (c counterTable) GetValue(ID string) (count int, err error) {
	if ID == "" {
		return 0, errors.ErrInvalidID
	}

	q := fmt.Sprintf("SELECT count FROM %s WHERE id = ?", c)
	err = service.Cassandra().Query(q, ID).Scan(&count)

	// If it doesn't have a row in this table, it's because its counter is 0.
	if err == gocql.ErrNotFound {
		return 0, nil
	}

	if err != nil {
		log.LogError("counter.GetValue", "Could not get the counter", err)
		return 0, errors.ErrUnexpected
	} else {
		return count, nil
	}
}
