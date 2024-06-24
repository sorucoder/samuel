package samuel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/sorucoder/samuel/internal/database"
)

type programModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func getProgramModelByID(context context.Context, transaction *database.Transaction, programID string) (*programModel, error) {
	model := new(programModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `programs` WHERE `id` = ?", programID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Program struct {
	model *programModel
	valid bool
}

var (
	ErrProgramInvalid error = errors.New("invalid program")
)

func getProgramByStudent(context context.Context, transaction *database.Transaction, student *Student) (*Program, error) {
	program := new(Program)

	var errGetModel error
	program.model, errGetModel = getProgramModelByID(context, transaction, student.model.ProgramID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	program.valid = true

	return program, nil
}

func (program *Program) MarshalJSON() ([]byte, error) {
	if !program.valid {
		panic(ErrProgramInvalid)
	}

	return json.Marshal(program.model.Name)
}
