package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/anukool23/usermanagement-lms-go/internal/storage"
	"github.com/anukool23/usermanagement-lms-go/internal/types"
	"github.com/anukool23/usermanagement-lms-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating new student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			slog.Error("Empty request body", slog.String("error", err.Error()))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err :=storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil{
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			slog.Error("Failed to create student", slog.String("error", err.Error()))
			return
		}
		slog.Info("Student created successfully", slog.Int64("id", lastId))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}
