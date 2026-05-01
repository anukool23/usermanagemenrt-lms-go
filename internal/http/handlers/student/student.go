package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

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

		existingStudent, err := storage.GetStudentByEmail(student.Email)
		if err == nil && existingStudent.Id != 0 {
			response.WriteJSON(w, http.StatusConflict, response.GeneralError(fmt.Errorf("student with email %s already exists", student.Email)))
			slog.Error("Student with email already exists", slog.String("email", student.Email))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			slog.Error("Failed to create student", slog.String("error", err.Error()))
			return
		}
		slog.Info("Student created successfully", slog.Int64("id", lastId))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting student by ID", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			slog.Error("Invalid student ID", slog.String("error", err.Error()))
			return
		}
		studentDetail, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			slog.Error("Failed to get student by ID", slog.String("error", err.Error()))
			return
		}
		slog.Info("Student retrieved successfully", slog.String("id", id))
		response.WriteJSON(w, http.StatusOK, studentDetail)

	}
}

func ListStudents(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Listing all students")

		students, err := storage.ListStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		slog.Info("Students retrieved successfully", slog.Int("count", len(students)))
		response.WriteJSON(w, http.StatusOK, students)
	}
}

func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Deleting student by ID", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil{
			response.WriteJSON(w, http.StatusBadGateway, response.GeneralError(err))
		}

			err = storage.DeleteStudentById(intId)
			if err != nil {
				response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
				slog.Error("Failed to delete student by ID", slog.String("error", err.Error()))
				return
			}
			slog.Info("Student deleted successfully", slog.String("id", id))
			response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
	}
}
