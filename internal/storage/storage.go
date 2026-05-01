package storage

import "github.com/anukool23/usermanagement-lms-go/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	ListStudents() ([]types.Student, error)
	GetStudentByEmail(email string) (types.Student, error)
	DeleteStudentById(id int64) error
}
