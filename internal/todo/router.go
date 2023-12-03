package todo

import (
	"github.com/kgoralski/go-crud-template/cmd/middleware"

	"github.com/go-chi/chi"
)

// Router structs represents Banks Handlers
type Router struct {
	r *chi.Mux
}

// Routes , all todos routes
func AddTodoRoutes(r *chi.Mux, controller *TodoController) {
	r.Get("/rest/todos/", middleware.CommonHeaders(controller.getAll()))
	r.Get("/rest/todos/{id:[0-9]+}", middleware.CommonHeaders(controller.getByID()))
	r.Post("/rest/todos/", middleware.CommonHeaders(controller.create()))
	r.Delete("/rest/todos/{id:[0-9]+}", middleware.CommonHeaders(controller.deleteByID()))
	r.Put("/rest/todos/{id:[0-9]+}", middleware.CommonHeaders(controller.update()))
	r.Delete("/rest/todos/", middleware.CommonHeaders(controller.deleteAll()))
}
