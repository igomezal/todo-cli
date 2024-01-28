package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const todoDirectory = ".todo"

type status int

const (
	Pending status = iota
	Done
)

func (s status) String() string {
	return [2]string{"todo", "done"}[s]
}

type Todo struct {
	ID            int
	Todo          string
	State         status
	DateCreated   time.Time // Probar si funciona bien el time.Time
	DateCompleted sql.NullTime
}

type todoDB struct {
	db *sql.DB
}

func NewTodoDB() (*todoDB, error) {
	homeUserDir, err := os.UserHomeDir()

	if err != nil {
		return nil, errors.New("Home User Directory couldn't be used")
	}

	todoFullPathDirectory := filepath.Join(homeUserDir, todoDirectory)

	if _, err := os.Stat(todoFullPathDirectory); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(todoFullPathDirectory, 0o770)
		}
	}

	todoFullPathFile := filepath.Join(todoFullPathDirectory, "todos.db")

	var todoDB = &todoDB{}

	todoDB.db, err = sql.Open("sqlite3", "file:"+todoFullPathFile)

	if err != nil {
		return nil, err
	}

	err = todoDB.setupTodoSchema()

	if err != nil {
		return nil, err
	}

	return todoDB, nil
}

func (t *todoDB) setupTodoSchema() error {
	_, err := t.db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id               INTEGER PRIMARY KEY AUTOINCREMENT,
			todo             VARCHAR(255) NOT NULL,
			state            INTEGER NOT NULL,
			date_created     DATETIME NOT NULL,
			date_completed   DATETIME
		);
	`)

	if err != nil {
		return err
	}

	return nil
}

func (t *todoDB) Close() error {
	return t.db.Close()
}

func getTodosHelper(functionName string, db *sql.DB, predicate string, filters ...any) ([]Todo, error) {
	var todos []Todo

	rows, err := db.Query(predicate, filters...)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", functionName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Todo,
			&todo.State,
			&todo.DateCreated,
			&todo.DateCompleted,
		)

		if err != nil {
			return nil, fmt.Errorf("%q: %w", functionName, err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%q: %w", functionName, err)
	}

	return todos, nil
}

func (t *todoDB) GetTasks() ([]Todo, error) {
	return getTodosHelper("GetTasks", t.db, "SELECT * FROM todos")
}

func (t *todoDB) GetFilteredTasksByState(state status) ([]Todo, error) {
	return getTodosHelper("GetFilteredTasksByState", t.db, "SELECT * FROM todos WHERE state = ?", state)
}

func (t *todoDB) GetFilteredTasksByCreationDate(time time.Time) ([]Todo, error) {
	return getTodosHelper("GetFilteredTasksByCreationDate", t.db, "SELECT * FROM todos WHERE date(date_created) = date(?)", time)
}

func (t *todoDB) GetFilteredTasksByStateAndDate(state status, time time.Time) ([]Todo, error) {
	return getTodosHelper("GetFilteredTasksByState", t.db, "SELECT * FROM todos WHERE state = ? AND date(date_created) = date(?)", state, time)
}

func (t *todoDB) CreateTodo(title string, tags string) error {
	_, err := t.db.Exec(`
		INSERT INTO todos
			(todo, state, date_created)
		VALUES
			(?,?,?)
	`, title, Pending, time.Now())

	return err
}

func (t *todoDB) CompleteTodo(todoId int) error {
	var state status

	row := t.db.QueryRow("SELECT state FROM todos WHERE id = ?", todoId)
	err := row.Scan(&state)

	if err != nil {
		return err
	}

	if state == Done {
		return nil
	}

	_, err = t.db.Exec(`
		UPDATE todos SET state = ?, date_completed = ? WHERE id = ?
	`, Done, time.Now(), todoId)

	return err
}

func (t *todoDB) UncompleteTodo(todoId int) error {
	_, err := t.db.Exec(`
		UPDATE todos SET state = ?, date_completed = null WHERE id = ?
	`, Pending, todoId)

	return err
}

func (t *todoDB) ChangeTodoName(todoId int, newName string) error {
	_, err := t.db.Exec(`
		UPDATE todos SET todo = ? WHERE id = ?
	`, newName, todoId)

	return err
}

func (t *todoDB) DeleteTodo(todoId int) error {
	_, err := t.db.Exec(`
		DELETE FROM todos WHERE id = ?
	`, todoId)

	return err
}
