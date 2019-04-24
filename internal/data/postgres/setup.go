package postgres

import (
	"database/sql"
)

// TestSetup tears down DB before setting it up, returns teardown function
func TestSetup(db *sql.DB) (teardown func() (sql.Result, error), err error) {

	teardown = func() (sql.Result, error) {
		return db.Exec("DROP TABLE task")
	}

	_, err = db.Exec("SELECT 1 FROM task LIMIT 1")
	if err == nil {
		// Table exists, teardown first
		teardown()
	}

	_, err = setup(db)

	return
}

// Setup sets up initial DB schema
func setup(db *sql.DB) (setup bool, err error) {

	_, err = db.Exec(`SELECT 1 FROM task LIMIT 1`)
	if err == nil {
		setup = false
		return // no error, table was setup
	}

	setup = true
	_, err = db.Exec(`CREATE TABLE task (
		id SERIAL PRIMARY KEY,
		name character varying(100) NOT NULL,
		description character varying(500) NOT NULL,
		completed_time TIMESTAMPTZ,
		cleared_time TIMESTAMPTZ
		);
		SET timezone = 'GMT'`)

	return
}
