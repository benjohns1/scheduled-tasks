package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // add postgres DB driver
)

// DBConn contains DB connection data
type DBConn struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// NewDBConn loads default DB connection info, overrides with environment variables
func NewDBConn() DBConn {

	// Defaults
	conn := DBConn{
		Host:     "localhost",
		Name:     "taskapp",
		Password: "postgresDefault",
		Port:     5432,
		User:     "postgresUser",
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

	return conn
}

// Connect opens and ping-checks a DB connection
func connect(conn DBConn) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conn.Host, conn.Port, conn.User, conn.Password, conn.Name))
	if err != nil {
		err = fmt.Errorf("error opening db: %v", err)
		log.Println(err)
	}

	// DB connection retry logic
	maxAttempts := 20
	retrySleep := 3
	for attempts := 0; attempts < maxAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("couldn't ping db: %v", err)
		time.Sleep(time.Duration(retrySleep) * time.Second)
	}
	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("error pinging db: %v", err)
	}

	return
}
