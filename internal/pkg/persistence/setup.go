package persistence

import (
	"database/sql"
)

// Setup sets up initial DB schema
func Setup(db *sql.DB) (setup bool, err error) {

	_, err = db.Exec(`SELECT 1 FROM task LIMIT 1`)
	if err == nil {
		setup = false
		return // no error, table was setup
	}

	setup = true
	_, err = db.Exec(`CREATE TABLE task (
		key serial NOT NULL,
		name character varying(100) NOT NULL,
		description character varying(500) NOT NULL,
		CONSTRAINT task_pkey PRIMARY KEY (key)
		)
		WITH (OIDS=FALSE)`)

	return
}
