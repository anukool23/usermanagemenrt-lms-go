package types

type Student struct {
	Id    int
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int
}
