package taskapp

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
)

func setupSQLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./taskapp.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestShouldReturnNewTaskKeys(t *testing.T) {
	testTasks := []*Task{
		{Name: "task name", Description: "task description"},
		{Name: "asdf", Description: "1234"},
	}

	for _, task := range testTasks {
		t.Run(fmt.Sprintf("task:%v", task.Name), func(t *testing.T) {

			//db, err := setupSQLiteDB()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("setting up DB: %v", err)
			}
			defer db.Close()

			stmt := mock.ExpectPrepare("^INSERT INTO task")
			var expectedKey TaskKey = 1
			rows := sqlmock.NewRows([]string{"key"}).AddRow(expectedKey)
			stmt.ExpectQuery().WithArgs(task.Name, task.Description).WillReturnRows(rows)

			key, err := AddTask(db, task)
			if err != nil {
				t.Errorf("error not expected while adding task: %v", err)
			}
			if key != expectedKey {
				t.Errorf("TaskKey %v does not equal expected %v", key, expectedKey)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled db mock expectations: %v", err)
			}
		})
	}
}
