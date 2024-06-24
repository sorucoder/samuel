package samuel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
)

type administratorModel struct {
	UserUUID  uuid.UUID `db:"user_uuid"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
}

func getAdministratorModelByUserUUID(context context.Context, transaction *database.Transaction, administratorUserUUID uuid.UUID) (*administratorModel, error) {
	model := new(administratorModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `administrators` WHERE `user_uuid` = ?", administratorUserUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Administrator struct {
	model *administratorModel
	user  *User
	valid bool
}

var (
	ErrAdministratorInvalid error = errors.New("administrator invalid")
)

func getAdministratorByUser(context context.Context, transaction *database.Transaction, user *User) (*Administrator, error) {
	administrator := new(Administrator)

	var errGetModel error
	administrator.model, errGetModel = getAdministratorModelByUserUUID(context, transaction, user.model.UUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	administrator.user = user

	administrator.valid = true

	return administrator, nil
}

func GetAdministratorByUser(context context.Context, user *User) (*Administrator, error) {
	if !user.valid {
		panic(ErrUserInvalid)
	}
	if !user.model.is("administrator") {
		panic(ErrUserNotAdministrator)
	}

	errPing := database.Ping(context)
	if errPing != nil {
		return nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, errBegin
	}

	administrator, errGetAdministrator := getAdministratorByUser(context, transaction, user)
	if errGetAdministrator != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, errors.Join(errGetAdministrator, errRollback)
		}

		return nil, errGetAdministrator
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, errCommit
	}

	return administrator, nil
}

func (adminstrator *Administrator) MarshalJSON() ([]byte, error) {
	if !adminstrator.valid {
		panic(ErrAdministratorInvalid)
	}

	return json.Marshal(map[string]any{
		"firstName": adminstrator.model.FirstName,
		"lastName":  adminstrator.model.LastName,
		"email":     adminstrator.model.Email,
		"phone":     adminstrator.model.Phone,
	})
}
