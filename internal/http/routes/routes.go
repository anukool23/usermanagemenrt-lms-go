package routes

import (
	"net/http"

	"github.com/anukool23/usermanagement-lms-go/internal/http/handlers/student"
	"github.com/anukool23/usermanagement-lms-go/internal/storage"
)

func Register(mux *http.ServeMux, s storage.Storage) {
	mux.HandleFunc("POST /api/v1/student", student.New(s))
	mux.HandleFunc("GET /api/v1/student/{id}", student.GetById(s))
	mux.HandleFunc("GET /api/v1/student", student.ListStudents(s))
	mux.HandleFunc("DELETE /api/v1/student/{id}", student.DeleteById(s))
}
