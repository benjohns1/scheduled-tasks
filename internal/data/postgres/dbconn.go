package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // add postgres DB driver
)

// DBConn contains DB connection data
type DBConn struct {
	Host              string
	Port              int
	User              string
	Password          string
	Name              string
	MaxRetryAttempts  int
	RetrySleepSeconds int
}

// NewDBConn creates struct with default DB connection info, and overrides with environment variables if set
func NewDBConn() DBConn {

	// Defaults
	conn := DBConn{
		Host:              "localhost",
		Name:              "taskapp",
		Password:          "postgresDefault",
		Port:              5432,
		User:              "postgresUser",
		MaxRetryAttempts:  20,
		RetrySleepSeconds: 3,
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
func connect(l Logger, conn DBConn) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conn.Host, conn.Port, conn.User, conn.Password, conn.Name))
	if err != nil {
		err = fmt.Errorf("error opening db: %v", err)
		l.Println(err)
		return
	}

	// Ping & retry if needed
	for attempts := 0; attempts < conn.MaxRetryAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			break
		}
		l.Printf("attempt %d/%d couldn't ping db: %v", attempts+1, conn.MaxRetryAttempts, err)
		time.Sleep(time.Duration(conn.RetrySleepSeconds) * time.Second)
	}

	return
}
