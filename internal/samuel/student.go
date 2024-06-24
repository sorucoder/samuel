package samuel

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
)

type studentModel struct {
	UserUUID  uuid.UUID       `db:"user_uuid"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
	Address   string          `db:"address"`
	Unit      *sql.NullString `db:"unit"`
	City      string          `db:"city"`
	State     string          `db:"state"`
	ZIP       string          `db:"zip"`
	Email     string          `db:"email"`
	Phone     string          `db:"phone"`
	CampusID  string          `db:"campus_id"`
	ProgramID string          `db:"program_id"`
}

func getStudentModelByUserUUID(context context.Context, transaction *database.Transaction, studentUserUUID uuid.UUID) (*studentModel, error) {
	model := new(studentModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `students` WHERE `user_uuid` = ?", studentUserUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Student struct {
	model   *studentModel
	campus  *Campus
	program *Program
	valid   bool
}

var (
	ErrStudentInvalid error = errors.New("student invalid")
)

func getStudentByUser(context context.Context, transaction *database.Transaction, user *User) (*Student, error) {
	student := new(Student)

	var errGetModel error
	student.model, errGetModel = getStudentModelByUserUUID(context, transaction, user.model.UUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetCampus error
	student.campus, errGetCampus = getCampusByStudent(context, transaction, student)
	if errGetCampus != nil {
		return nil, errGetCampus
	}

	var errGetProgram error
	student.program, errGetProgram = getProgramByStudent(context, transaction, student)
	if errGetProgram != nil {
		return nil, errGetProgram
	}

	student.valid = true

	return student, nil
}

func GetStudentByUser(context context.Context, user *User) (*Student, error) {
	if !user.valid {
		panic(ErrUserInvalid)
	}
	if !user.model.is("student") {
		panic(ErrUserNotStudent)
	}

	errPing := database.Ping(context)
	if errPing != nil {
		return nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, errBegin
	}

	student, errGetStudent := getStudentByUser(context, transaction, user)
	if errGetStudent != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, errors.Join(errGetStudent, errRollback)
		}

		return nil, errGetStudent
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, errCommit
	}

	return student, nil
}

func (student *Student) MarshalJSON() ([]byte, error) {
	if !student.valid {
		panic(ErrStudentInvalid)
	}

	studentMap := map[string]any{
		"firstName": student.model.FirstName,
		"lastName":  student.model.LastName,
		"address":   student.model.Address,
		"city":      student.model.City,
		"state":     student.model.State,
		"zip":       student.model.ZIP,
		"email":     student.model.Email,
		"phone":     student.model.Phone,
		"campus":    student.campus,
		"program":   student.program,
	}
	if student.model.Unit.Valid {
		studentMap["unit"] = student.model.Unit.String
	}

	return json.Marshal(studentMap)
}
