package todolist

import (
	"net/http"
)

func (h *TodoListHandlers) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement list logic
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "not implemented yet"}`))
}

func (h *TodoListHandlers) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get logic
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "not implemented yet"}`))
}

func (h *TodoListHandlers) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement create logic
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "not implemented yet"}`))
}

func (h *TodoListHandlers) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement update logic
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "not implemented yet"}`))
}

func (h *TodoListHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement delete logic
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "not implemented yet"}`))
}
