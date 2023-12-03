package todo

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // DB driver registration
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// how the todo is stored in the database
type Todo struct {
	ID   int    `json:"id" DB:"id"`
	Name string `json:"name" DB:"name"`
}

// TodoStore for persistence
type TodoStore struct {
	db *sqlx.DB
}

// NewStore creates new TodoStore for Todos
func NewTodoStore(db *sqlx.DB) *TodoStore {
	return &TodoStore{db: db}
}

func (s *TodoStore) getAll() ([]Todo, error) {
	var todos []Todo
	if err := s.db.Select(&todos, "SELECT * FROM todos"); err != nil {
		return nil, ErrDbQuery{Err: errors.Wrap(err, "TodoStore.getAll() error")}
	}
	if todos == nil {
		return []Todo{}, nil
	}
	return todos, nil
}

func (s *TodoStore) get(id int) (*Todo, error) {
	var bank Todo
	if err := s.db.Get(&bank, "SELECT * FROM todos WHERE id=?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEntityNotFound{Err: errors.Wrap(err, "TodoStore.get() ErrNoRows error")}
		}
		return nil, ErrDbQuery{Err: errors.Wrap(err, "TodoStore.get() error")}
	}
	return &bank, nil
}

func (s *TodoStore) create(bank Todo) (int, error) {
	result, err := s.db.Exec("INSERT into todos (name) VALUES (?)", bank.Name)
	if err != nil {
		return 0, ErrDbQuery{Err: errors.Wrap(err, "")}
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, ErrDbNotSupported{Err: errors.Wrap(err, "TodoStore.create() error")}
	}
	return int(lastID), nil
}

func (s *TodoStore) deleteAll() error {
	if _, err := s.db.Exec("TRUNCATE table todos"); err != nil {
		return ErrDbQuery{Err: errors.Wrap(err, "TodoStore.deleteAll() error")}
	}
	return nil
}

func (s *TodoStore) delete(id int) error {
	res, err := s.db.Exec("DELETE from todos where id=?", id)
	if err != nil {
		return ErrDbQuery{Err: errors.Wrap(err, "TodoStore.delete() error")}
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return ErrDbQuery{Err: errors.Wrap(err, "TodoStore.delete() RowsAffected error")}
	}
	if affect == 0 {
		return ErrEntityNotFound{Err: errors.Wrap(err, "TodoStore.delete() NotFound error")}
	}
	return nil
}

func (s *TodoStore) update(bank Todo) (*Todo, error) {
	res, err := s.db.Exec("UPDATE todos SET name=? WHERE id=?", bank.Name, bank.ID)
	if err != nil {
		return nil, ErrDbQuery{Err: errors.Wrap(err, "TodoStore.update() error")}
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, ErrDbQuery{Err: errors.Wrap(err, "TodoStore.update() RowsAffected error")}
	}
	if affect == 0 {
		return nil, ErrEntityNotFound{Err: errors.Wrap(err, "TodoStore.update() NotFound error")}
	}
	return &bank, nil
}
