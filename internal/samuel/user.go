package samuel

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
	"github.com/sorucoder/samuel/internal/ldap"
	"golang.org/x/crypto/bcrypt"
)

type userModel struct {
	UUID         uuid.UUID `db:"uuid"`
	Identity     string    `db:"identity"`
	PasswordHash []byte    `db:"password_hash"`
	RoleID       string    `db:"role_id"`
	CreatedOn    time.Time `db:"created_on"`
}

func getUserModelByUUID(context context.Context, transaction *database.Transaction, userUUID uuid.UUID) (*userModel, error) {
	model := new(userModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `users` WHERE `uuid` = ?", userUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func getUserModelByIdentity(context context.Context, transaction *database.Transaction, userIdentity string) (*userModel, error) {
	model := new(userModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `users` WHERE `identity` = ?", userIdentity)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func (model *userModel) is(roleID string) bool {
	return model.RoleID == roleID
}

func (model *userModel) updatePasswordHash(context context.Context, transaction *database.Transaction, newUserPasswordHash []byte) error {
	_, errUpdate := transaction.Execute(context, "UPDATE `users` SET `password_hash` = ? WHERE `uuid` = ?", newUserPasswordHash, model.UUID)
	if errUpdate != nil {
		return errUpdate
	}

	model.PasswordHash = newUserPasswordHash

	return nil
}

type User struct {
	model *userModel
	role  *Role
	valid bool
}

var (
	ErrUserInvalid          error = errors.New("user invalid")
	ErrUserUsesLDAP         error = errors.New("user uses ldap")
	ErrUserNotAdministrator error = errors.New("user not administrator")
	ErrUserNotInstructor    error = errors.New("user not instructor")
	ErrUserNotSupervisor    error = errors.New("user not supervisor")
	ErrUserNotStudent       error = errors.New("user not student")
)

func getUserByUUID(context context.Context, transaction *database.Transaction, userUUID uuid.UUID) (*User, error) {
	user := new(User)

	var errGetModel error
	user.model, errGetModel = getUserModelByUUID(context, transaction, userUUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetRole error
	user.role, errGetRole = getRoleByID(context, transaction, user.model.RoleID)
	if errGetRole != nil {
		return nil, errGetRole
	}

	user.valid = true

	return user, nil
}

func getUserByIdentity(context context.Context, transaction *database.Transaction, userIdentity string) (*User, error) {
	user := new(User)

	var errGetModel error
	user.model, errGetModel = getUserModelByIdentity(context, transaction, userIdentity)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetRole error
	user.role, errGetRole = getRoleByID(context, transaction, user.model.RoleID)
	if errGetRole != nil {
		return nil, errGetRole
	}

	user.valid = true

	return user, nil
}

func getUserBySession(context context.Context, transaction *database.Transaction, session *Session) (*User, error) {
	user := new(User)

	var errGetModel error
	user.model, errGetModel = getUserModelByUUID(context, transaction, session.model.UserUUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetRole error
	user.role, errGetRole = getRoleByID(context, transaction, user.model.RoleID)
	if errGetRole != nil {
		return nil, errGetRole
	}

	user.valid = true

	return user, nil
}

func getUserByPasswordChange(context context.Context, transaction *database.Transaction, passwordChange *PasswordChange) (*User, error) {
	user := new(User)

	var errGetModel error
	user.model, errGetModel = getUserModelByUUID(context, transaction, passwordChange.model.SupervisorUUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	var errGetRole error
	user.role, errGetRole = getRoleByID(context, transaction, user.model.RoleID)
	if errGetRole != nil {
		return nil, errGetRole
	}

	user.valid = true

	return user, nil
}

func LoginUser(context context.Context, userIdentity string, userPassword string) (*User, *Session, error) {
	errPing := database.Ping(context)
	if errPing != nil {
		return nil, nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, nil, errBegin
	}

	user, errGetUser := getUserByIdentity(context, transaction, userIdentity)
	if errGetUser != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, nil, errors.Join(errGetUser, errRollback)
		}

		return nil, nil, errGetUser
	}

	switch user.model.RoleID {
	case "administrator", "instructor", "student":
		errAuthenticate := ldap.Authenticate(context, userIdentity, userPassword)
		if errAuthenticate != nil {
			return nil, nil, errAuthenticate
		}
	case "supervisor":
		errAuthenticate := bcrypt.CompareHashAndPassword(user.model.PasswordHash, []byte(userPassword))
		if errAuthenticate != nil {
			return nil, nil, errAuthenticate
		}
	}

	session, errStartSession := beginSession(context, transaction, user)
	if errStartSession != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, nil, errors.Join(errStartSession, errRollback)
		}

		return nil, nil, errStartSession
	}

	errRecord := recordAudit(context, transaction, "Logged in.", user)
	if errRecord != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, nil, errors.Join(errRecord, errRollback)
		}

		return nil, nil, errRecord
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, nil, errCommit
	}

	return user, session, nil
}

func AuthenticateSession(context context.Context, sessionToken uuid.UUID) (*User, *Session, error) {
	errPing := database.Ping(context)
	if errPing != nil {
		return nil, nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, nil, errBegin
	}

	session, errContinueSession := continueSession(context, transaction, sessionToken)
	if errContinueSession != nil {
		if errors.Is(errContinueSession, ErrSessionExpired) {
			errCommit := transaction.Commit()
			if errCommit != nil {
				return nil, nil, errors.Join(errContinueSession, errCommit)
			}

			return nil, nil, errContinueSession
		} else {
			errRollback := transaction.Rollback()
			if errRollback != nil {
				return nil, nil, errors.Join(errContinueSession, errRollback)
			}

			return nil, nil, errContinueSession
		}
	}

	user, errGetUser := getUserBySession(context, transaction, session)
	if errGetUser != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, nil, errors.Join(errGetUser, errRollback)
		}

		return nil, nil, errGetUser
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, nil, errCommit
	}

	return user, session, nil
}

func LogoutUser(context context.Context, user *User, session *Session) error {
	if !user.valid {
		panic(ErrUserInvalid)
	}
	if !session.valid {
		panic(ErrSessionInvalid)
	}

	errPing := database.Ping(context)
	if errPing != nil {
		return errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return errBegin
	}

	errEnd := session.end(context, transaction)
	if errEnd != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return errors.Join(errEnd, errRollback)
		}

		return errEnd
	}

	errRecord := recordAudit(context, transaction, "Logged out.", user)
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

func (user *User) Role() *Role {
	if !user.valid {
		panic(ErrUserInvalid)
	}

	return user.role
}

func (user *User) Is(roleID string) bool {
	if !user.valid {
		panic(ErrUserInvalid)
	}

	return user.model.is(roleID)
}

func (user *User) changePassword(context context.Context, transaction *database.Transaction, newUserPassword string) error {
	if !user.model.is("supervisor") {
		return ErrUserUsesLDAP
	}

	newUserPasswordHash, errHash := bcrypt.GenerateFromPassword([]byte(newUserPassword), 0)
	if errHash != nil {
		return errHash
	}

	errUpdateModel := user.model.updatePasswordHash(context, transaction, newUserPasswordHash)
	if errUpdateModel != nil {
		return errUpdateModel
	}

	return nil
}

func (user *User) MarshalJSON() ([]byte, error) {
	if !user.valid {
		panic(ErrUserInvalid)
	}

	return json.Marshal(map[string]any{
		"uuid":      user.model.UUID,
		"role":      user.role,
		"createdOn": user.model.CreatedOn,
	})
}
