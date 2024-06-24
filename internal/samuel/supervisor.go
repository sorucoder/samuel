package samuel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
	"github.com/sorucoder/samuel/internal/email"
)

type supervisorModel struct {
	UserUUID    uuid.UUID `db:"user_uuid"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Title       string    `db:"title"`
	Email       string    `db:"email"`
	Phone       string    `db:"phone"`
	CompanyUUID uuid.UUID `db:"company_uuid"`
}

func getSupervisorModelByUserUUID(context context.Context, transaction *database.Transaction, supervisorUserUUID uuid.UUID) (*supervisorModel, error) {
	model := new(supervisorModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `supervisors` WHERE `user_uuid` = ?", supervisorUserUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func getSupervisorModelByEmail(context context.Context, transaction *database.Transaction, supervisorEmail string) (*supervisorModel, error) {
	model := new(supervisorModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `supervisors` WHERE `email` = ?", supervisorEmail)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func (model *supervisorModel) address() *email.Address {
	return email.NewAddress(model.FirstName, model.LastName, model.Email)
}

type Supervisor struct {
	model   *supervisorModel
	company *Company
	valid   bool
}

var (
	ErrSupervisorInvalid error = errors.New("supervisor invalid")
)

func getSupervisorByEmail(context context.Context, transaction *database.Transaction, supervisorEmail string) (*Supervisor, error) {
	supervisor := new(Supervisor)

	var errGetModel error
	supervisor.model, errGetModel = getSupervisorModelByEmail(context, transaction, supervisorEmail)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetCompany error
	supervisor.company, errGetCompany = getCompanyBySupervisor(context, transaction, supervisor)
	if errGetCompany != nil {
		return nil, errGetCompany
	}

	supervisor.valid = true

	return supervisor, nil
}

func getSupervisorByUser(context context.Context, transaction *database.Transaction, user *User) (*Supervisor, error) {
	supervisor := new(Supervisor)

	var errGetModel error
	supervisor.model, errGetModel = getSupervisorModelByUserUUID(context, transaction, user.model.UUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetCompany error
	supervisor.company, errGetCompany = getCompanyBySupervisor(context, transaction, supervisor)
	if errGetCompany != nil {
		return nil, errGetCompany
	}

	supervisor.valid = true

	return supervisor, nil
}

func GetSupervisorByUser(context context.Context, user *User) (*Supervisor, error) {
	if !user.valid {
		panic(ErrUserInvalid)
	}
	if !user.model.is("supervisor") {
		panic(ErrUserNotSupervisor)
	}

	errPing := database.Ping(context)
	if errPing != nil {
		return nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, errBegin
	}

	supervisor, errGetSupervisor := getSupervisorByUser(context, transaction, user)
	if errGetSupervisor != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, errors.Join(errGetSupervisor, errRollback)
		}

		return nil, errGetSupervisor
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, errCommit
	}

	return supervisor, nil
}

func (supervisor *Supervisor) MarshalJSON() ([]byte, error) {
	if !supervisor.valid {
		panic(ErrSupervisorInvalid)
	}

	return json.Marshal(map[string]any{
		"firstName": supervisor.model.FirstName,
		"lastName":  supervisor.model.LastName,
		"title":     supervisor.model.Title,
		"email":     supervisor.model.Email,
		"phone":     supervisor.model.Phone,
		"company":   supervisor.company,
	})
}
