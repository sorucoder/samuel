package samuel

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
	"github.com/sorucoder/samuel/internal/email"
)

type passwordChangeModel struct {
	Token          uuid.UUID `db:"token"`
	SupervisorUUID uuid.UUID `db:"supervisor_uuid"`
	ExpiresOn      time.Time `db:"expires_on"`
}

func insertPasswordChangeModel(context context.Context, transaction *database.Transaction, passwordChangeSupervisorUUID uuid.UUID) (*passwordChangeModel, error) {
	_, errInsert := transaction.Execute(context, "INSERT INTO `password_changes` (`supervisor_uuid`) VALUE (?)", passwordChangeSupervisorUUID)
	if errInsert != nil {
		return nil, errInsert
	}

	model := new(passwordChangeModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `password_changes` WHERE `supervisor_uuid` = ?", passwordChangeSupervisorUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func getPasswordChangeModelByToken(context context.Context, transaction *database.Transaction, passwordChangeToken uuid.UUID) (*passwordChangeModel, error) {
	model := new(passwordChangeModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `password_changes` WHERE `token` = ?", passwordChangeToken)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func getPasswordChangeModelBySupervisorUUID(context context.Context, transaction *database.Transaction, passwoerdChangeSupervisorUUID uuid.UUID) (*passwordChangeModel, error) {
	model := new(passwordChangeModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `password_changes` WHERE `supervisor_uuid` = ?", passwoerdChangeSupervisorUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func (model *passwordChangeModel) expired() bool {
	return model.ExpiresOn.Before(time.Now())
}

func (model *passwordChangeModel) delete(context context.Context, transaction *database.Transaction) error {
	_, errDelete := transaction.Execute(context, "DELETE FROM `password_changes` WHERE `token` = ?", model.Token)
	if errDelete != nil {
		return errDelete
	}

	return nil
}

type PasswordChange struct {
	model *passwordChangeModel
	valid bool
}

var (
	ErrPasswordChangeInvalid error = errors.New("password change invalid")
	ErrPasswordChangeExists  error = errors.New("password change exists")
)

func newPasswordChange(context context.Context, transaction *database.Transaction, supervisor *Supervisor) (*PasswordChange, error) {
	passwordChange := new(PasswordChange)

	var errInsertModel error
	passwordChange.model, errInsertModel = insertPasswordChangeModel(context, transaction, supervisor.model.UserUUID)
	if errInsertModel != nil {
		return nil, errInsertModel
	}

	passwordChange.valid = true

	return passwordChange, nil
}

func getPasswordChangeByToken(context context.Context, transaction *database.Transaction, passwordChangeToken uuid.UUID) (*PasswordChange, error) {
	passwordChange := new(PasswordChange)

	var errGetModel error
	passwordChange.model, errGetModel = getPasswordChangeModelByToken(context, transaction, passwordChangeToken)
	if errGetModel != nil {
		return nil, errGetModel
	}

	passwordChange.valid = false

	return passwordChange, nil
}

func getPasswordChangeBySupervisor(context context.Context, transaction *database.Transaction, supervisor *Supervisor) (*PasswordChange, error) {
	passwordChange := new(PasswordChange)

	var errGetModel error
	passwordChange.model, errGetModel = getPasswordChangeModelBySupervisorUUID(context, transaction, supervisor.model.UserUUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	passwordChange.valid = true

	return passwordChange, nil
}

func beginPasswordChange(context context.Context, transaction *database.Transaction, supervisor *Supervisor) (*PasswordChange, error) {
	passwordChange, errGet := getPasswordChangeBySupervisor(context, transaction, supervisor)
	if errGet == nil {
		expired := passwordChange.expired()
		if !expired {
			return passwordChange, ErrPasswordChangeExists
		}

		errEnd := passwordChange.end(context, transaction)
		if errEnd != nil {
			return nil, errEnd
		}
	}

	passwordChange, errNew := newPasswordChange(context, transaction, supervisor)
	if errNew != nil {
		return nil, errNew
	}

	return passwordChange, nil
}

func CreatePasswordChange(context context.Context, supervisorEmail string) error {
	errPing := database.Ping(context)
	if errPing != nil {
		return errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return errBegin
	}

	supervisor, errGetSupervisor := getSupervisorByEmail(context, transaction, supervisorEmail)
	if errGetSupervisor != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errGetSupervisor, errRollback)
		}

		return errGetSupervisor
	}

	user, errGetUser := getUserByUUID(context, transaction, supervisor.model.UserUUID)
	if errGetUser != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errGetSupervisor, errRollback)
		}

		return errGetUser
	}

	passwordChange, errBeginPasswordChange := beginPasswordChange(context, transaction, supervisor)
	if errBeginPasswordChange != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errBeginPasswordChange, errRollback)
		}

		return errBeginPasswordChange
	}

	errSend := email.Send(context, supervisor.model.address(), email.PasswordChangeRequestTemplate, map[string]any{
		"firstName":         supervisor.model.FirstName,
		"passwordChangeURL": email.GenerateLink("password_change", passwordChange.model.Token),
	})
	if errSend != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errSend, errRollback)
		}

		return errSend
	}

	errRecord := recordAudit(context, transaction, "Requested password change.", user)
	if errRecord != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errRecord, errRollback)
		}

		return errRecord
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func FulfillPasswordChange(context context.Context, passwordChangeToken uuid.UUID, newPassword string) error {
	errPing := database.Ping(context)
	if errPing != nil {
		return errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return errBegin
	}

	passwordChange, errGetPasswordChange := getPasswordChangeByToken(context, transaction, passwordChangeToken)
	if errGetPasswordChange != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errGetPasswordChange, errRollback)
		}

		return errGetPasswordChange
	}

	user, errGetUser := getUserByPasswordChange(context, transaction, passwordChange)
	if errGetUser != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errGetUser, errRollback)
		}

		return errGetUser
	}

	errChangePassword := user.changePassword(context, transaction, newPassword)
	if errChangePassword != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errChangePassword, errRollback)
		}

		return errChangePassword
	}

	errEndPasswordChange := passwordChange.end(context, transaction)
	if errEndPasswordChange != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errEndPasswordChange, errRollback)
		}

		return errEndPasswordChange
	}

	errRecord := recordAudit(context, transaction, "Fulfilled password change.", user)
	if errRecord != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errRecord, errRollback)
		}

		return errRecord
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (passwordChange *PasswordChange) expired() bool {
	return passwordChange.model.expired()
}

func (passwordChange *PasswordChange) end(context context.Context, transaction *database.Transaction) error {
	errDeleteModel := passwordChange.model.delete(context, transaction)
	if errDeleteModel != nil {
		return errDeleteModel
	}

	passwordChange.valid = false

	return nil
}
