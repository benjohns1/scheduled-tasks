package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // add postgres DB driver
)

// Logger interface needed for postgres log messages
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type scannable interface {
	Scan(dest ...interface{}) error
}

const dbTimeFormat = time.RFC3339Nano

// DBConn contains DB connection data
type DBConn struct {
	Host              string
	Port              int
	User              string
	Password          string
	Name              string
	MaxRetryAttempts  int
	RetrySleepSeconds int
	DB                *sql.DB
	l                 Logger
	AppName           string
}

// Close closes the wrapped DB connection
func (conn *DBConn) Close() error {
	if conn.DB == nil {
		return nil
	}
	return conn.DB.Close()
}

// NewDBConn creates struct with default DB connection info, and overrides with environment variables if set
func NewDBConn(l Logger, appName string) DBConn {

	// Defaults
	conn := DBConn{
		Host:              "localhost",
		Name:              "taskapp",
		Password:          "postgresDefault",
		Port:              5432,
		User:              "postgresUser",
		MaxRetryAttempts:  20,
		RetrySleepSeconds: 3,
		l:                 l,
		AppName:           appName,
	}

	// Override from env vars
	if host, exists := os.LookupEnv("POSTGRES_HOST"); exists {
		conn.Host = host
	}
	if name, exists := os.LookupEnv("POSTGRES_DB"); exists {
		conn.Name = name
	}
	if pass, exists := os.LookupEnv("POSTGRES_PASSWORD"); exists {
		conn.Password = pass
	}
	if user, exists := os.LookupEnv("POSTGRES_USER"); exists {
		conn.User = user
	}
	if port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT")); err == nil {
		conn.Port = port
	}
	if maxRetryAttempts, err := strconv.Atoi(os.Getenv("DBCONN_MAXRETRYATTEMPTS")); err == nil {
		conn.MaxRetryAttempts = maxRetryAttempts
	}
	if retrySleepSeconds, err := strconv.Atoi(os.Getenv("DBCONN_RETRYSLEEPSECONDS")); err == nil {
		conn.RetrySleepSeconds = retrySleepSeconds
	}

	return conn
}

// Connect opens and ping-checks a DB connection
func (conn *DBConn) Connect() (err error) {
	if conn.DB == nil {
		conn.l.Printf("connecting to db %s as %s...", conn.Name, conn.User)
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s application_name=%s sslmode=disable", conn.Host, conn.Port, conn.User, conn.Password, conn.Name, conn.AppName))
		if err != nil {
			err = fmt.Errorf("error opening db: %v", err)
			return err
		}
		conn.DB = db
	} else {
		conn.l.Printf("already connected to db %s", conn.Name)
	}

	// Ping & retry if needed
	for attempts := 0; attempts < conn.MaxRetryAttempts; attempts++ {
		err = conn.DB.Ping()
		if err == nil {
			break
		}
		conn.l.Printf("attempt %d/%d couldn't ping db: %v", attempts+1, conn.MaxRetryAttempts, err)
		time.Sleep(time.Duration(conn.RetrySleepSeconds) * time.Second)
	}

	return err
}

// Setup sets up initial DB schema
func (conn *DBConn) Setup() (setup bool, err error) {

	_, err = conn.DB.Exec(`SELECT 1 FROM task LIMIT 1; SELECT 1 FROM schedule LIMIT 1; SELECT 1 FROM recurring_task;`)
	if err == nil {
		setup = false
		return // no error, table was already setup
	}

	setup = true
	_, err = conn.DB.Exec(`
		CREATE TABLE user_account (
			id uuid PRIMARY KEY,
			displayname character varying(100)
			);
		CREATE TABLE user_external (
			user_id uuid REFERENCES user_account(id) ON DELETE CASCADE ON UPDATE CASCADE,
			provider character varying(100) NOT NULL,
			external_id character varying (100) NOT NULL,
			PRIMARY KEY(provider, external_id)
			);
		CREATE TABLE task (
			id SERIAL PRIMARY KEY,
			name character varying(100) NOT NULL,
			description character varying(500) NOT NULL,
			completed_time TIMESTAMPTZ,
			cleared_time TIMESTAMPTZ,
			created_time TIMESTAMPTZ,
			created_by uuid REFERENCES user_account(id)
			);
		CREATE TABLE schedule (
			id SERIAL PRIMARY KEY,
			paused boolean NOT NULL,
			last_checked TIMESTAMPTZ,
			removed_time TIMESTAMPTZ,
			created_by uuid REFERENCES user_account(id),
			frequency_offset integer NOT NULL,
			frequency_interval integer NOT NULL,
			frequency_time_period smallint NOT NULL,
			frequency_at_minutes smallint[],
			frequency_at_hours smallint[],
			frequency_on_days_of_week smallint[],
			frequency_on_days_of_month smallint[]
			);
		CREATE TABLE recurring_task (
			id SERIAL PRIMARY KEY,
			schedule_id integer REFERENCES schedule(id) ON DELETE CASCADE ON UPDATE CASCADE,
			name character varying(100) NOT NULL,
			description character varying(500) NOT NULL
			);
		SET timezone = 'GMT'`)

	return
}
