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
			response.HandleError(w, r.URL.Path, "empty body", 400, err)
			return
		}

		if err != nil {
			response.HandleError(w, r.URL.Path, "invalid request payload", 400, err)
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			_ = response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		existingStudent, err := storage.GetStudentByEmail(student.Email)
		if err == nil && existingStudent.Id != 0 {
			response.HandleError(
				w,
				r.URL.Path,
				fmt.Sprintf("student with email %s already exists", student.Email),
				409,
				nil,
			)
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.HandleError(w, r.URL.Path, "failed to create student", 500, err)
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
			response.HandleError(w, r.URL.Path, "invalid student ID", 400, err)
			return
		}
		studentDetail, err := storage.GetStudentById(intId)
		if err != nil {
			response.HandleError(w, r.URL.Path, "failed to get student by id", 500, err)
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
				response.HandleError(w, r.URL.Path, "failed to list students", 400, err)
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
			response.HandleError(w, r.URL.Path, "invalid student ID", 400, err)
			return
		}

				err = storage.DeleteStudentById(intId)
				if err != nil {
					response.HandleError(w, r.URL.Path, "failed to delete student by id", 500, err)
					return
				}
			slog.Info("Student deleted successfully", slog.String("id", id))
			response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
	}
}
