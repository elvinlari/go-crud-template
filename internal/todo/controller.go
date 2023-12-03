package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TodoController struct {
	storage *TodoStore
}

func NewTodoController(storage *TodoStore) *TodoController {
	return &TodoController{
		storage: storage,
	}
}

// @Summary Create one todo.
// @Description creates one todo.
// @Tags todos
// @Accept */*
// @Produce json
// @Router /todos [post]
func (t *TodoController) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			handleErrors(w, err)
			return
		}
		id, err := t.storage.create(todo)
		if err != nil {
			handleErrors(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(id); err != nil {
			handleErrors(w, err)
			return
		}
	}
}

// @Summary Update one todo.
// @Description updates one todo.
// @Tags todos
// @Accept */*
// @Produce json
// @Router /todos/1 [post]
func (t *TodoController) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			handleErrors(w, errors.Wrap(err, http.StatusText(http.StatusBadRequest)))
			return
		}
		var todo Todo
		if errDecode := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			handleErrors(w, errDecode)
			return
		}
		updatedTodo, err := t.storage.update(Todo{ID: id, Name: todo.Name})
		if err != nil {
			handleErrors(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
			handleErrors(w, err)
			return
		}
	}
}

// @Summary Get one todo.
// @Description gets one todo.
// @Tags todo
// @Accept */*
// @Produce json
// @Router /todos/1 [get]
func (t *TodoController) getByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			handleErrors(w, errors.Wrap(err, http.StatusText(http.StatusBadRequest)))
			return
		}
		b, err := t.storage.get(id)
		if err != nil {
			handleErrors(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(b); err != nil {
			handleErrors(w, err)
			return
		}
	}
}

// @Summary Get all todos.
// @Description fetch every todo available.
// @Tags todos
// @Accept */*
// @Produce json
// @Success 200 {object} []Todo
// @Router /todos [get]
func (t *TodoController) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := t.storage.getAll()
		if err != nil {
			handleErrors(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			handleErrors(w, err)
			return
		}
	}

}

// @Summary Delete one todo.
// @Description delete one todo.
// @Tags todo
// @Accept */*
// @Produce json
// @Router /todos/1 [delete]
func (t *TodoController) deleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			handleErrors(w, errors.Wrap(err, http.StatusText(http.StatusBadRequest)))
			return
		}

		if err = t.storage.delete(id); err != nil {
			handleErrors(w, err)
			return
		}
	}
}

// @Summary Delete all todo.
// @Description delete all todos.
// @Tags todo
// @Accept */*
// @Produce json
// @Router /todos/ [delete]
func (t *TodoController) deleteAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.storage.deleteAll(); err != nil {
			handleErrors(w, err)
			return
		}
	}
}

// ErrorResponse represents json error structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// JSONError is converting error to JSON response
func JSONError(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ErrorResponse{error}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleErrors , DB errors to Rest mapper
func handleErrors(w http.ResponseWriter, err error) {
	const logFormat = "fatal: %+v\n"
	if strings.Contains(err.Error(), "connection refused") {
		log.Warnf(logFormat, err)
		JSONError(w, "DB_CONNECTION_FAIL", http.StatusServiceUnavailable)
		return
	}
	if err.Error() == http.StatusText(400) {
		log.Warnf(logFormat, err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch err.(type) {
	case ErrDbQuery:
		log.Warnf(logFormat, err.(ErrDbQuery).Err)
		JSONError(w, err.Error(), http.StatusConflict)
	case ErrDbNotSupported:
		log.Warnf(logFormat, err.(ErrDbNotSupported).Err)
		JSONError(w, err.Error(), http.StatusConflict)
	case ErrEntityNotFound:
		log.Warnf(logFormat, err.(ErrEntityNotFound).Err)
		JSONError(w, err.Error(), http.StatusNotFound)
	default:
		log.Warnf(logFormat, err)
		JSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return
}
