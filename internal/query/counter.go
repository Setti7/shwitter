package query

import (
	"fmt"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// Increment a cassandra counter.
//
// Returns ErrInvalidID for invalid IDs and ErrUnexpected for any other errors.
func IncrementCounterValue(ID string, c entity.CounterTable, value int) error {
	if ID == "" {
		return ErrInvalidID
	}

	q := fmt.Sprintf("UPDATE %s SET count = count + ? WHERE id = ?", c)
	err := service.Cassandra().Query(q, value, ID).Exec()

	if err != nil {
		log.LogError("query.IncrementCounterValue", "Could not increment the counter", err)
		return ErrUnexpected
	} else {
		return nil
	}
}

// Get the value for a cassandra counter.
//
// Returns ErrInvalidID for invalid IDs and ErrUnexpected for any other errors.
func GetCounterValue(ID string, c entity.CounterTable) (count int, err error) {
	if ID == "" {
		return 0, ErrInvalidID
	}

	q := fmt.Sprintf("SELECT count FROM %s WHERE id = ?", c)
	err = service.Cassandra().Query(q, ID).Scan(&count)

	// If it doesn't have a row in this table, it's because its counter is 0.
	if err == gocql.ErrNotFound {
		return 0, nil
	}

	if err != nil {
		log.LogError("query.GetCounterValue", "Could not get the counter", err)
		return 0, ErrUnexpected
	} else {
		return count, nil
	}
}