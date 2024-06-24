package samuel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
)

type instructorModel struct {
	UserUUID  uuid.UUID `db:"user_uuid"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	CampusID  string    `db:"campus_id"`
}

func getInstructorModelByUserUUID(context context.Context, transaction *database.Transaction, instructorUserUUID uuid.UUID) (*instructorModel, error) {
	model := new(instructorModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `instructors` WHERE `user_uuid` = ?", instructorUserUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Instructor struct {
	model  *instructorModel
	campus *Campus
	valid  bool
}

var (
	ErrInstructorInvalid error = errors.New("instructor invalid")
)

func getInstructorByUser(context context.Context, transaction *database.Transaction, user *User) (*Instructor, error) {
	instructor := new(Instructor)

	var errGetModel error
	instructor.model, errGetModel = getInstructorModelByUserUUID(context, transaction, user.model.UUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetCampus error
	instructor.campus, errGetCampus = getCampusByInstructor(context, transaction, instructor)
	if errGetCampus != nil {
		return nil, errGetCampus
	}

	instructor.valid = true

	return instructor, nil
}

func GetInstructorByUser(context context.Context, user *User) (*Instructor, error) {
	if !user.valid {
		panic(ErrUserInvalid)
	}
	if !user.model.is("instructor") {
		panic(ErrUserNotInstructor)
	}

	errPing := database.Ping(context)
	if errPing != nil {
		return nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, errBegin
	}

	instructor, errGetInstructor := getInstructorByUser(context, transaction, user)
	if errGetInstructor != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, errors.Join(errGetInstructor, errRollback)
		}

		return nil, errGetInstructor
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, errCommit
	}

	return instructor, nil
}

func (instructor *Instructor) MarshalJSON() ([]byte, error) {
	if !instructor.valid {
		panic(ErrInstructorInvalid)
	}

	return json.Marshal(map[string]any{
		"firstName": instructor.model.FirstName,
		"lastName":  instructor.model.LastName,
		"email":     instructor.model.Email,
		"phone":     instructor.model.Phone,
		"campus":    instructor.campus,
	})
}
