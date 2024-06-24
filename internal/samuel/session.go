package samuel

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
)

type sessionModel struct {
	Token     uuid.UUID `db:"token"`
	UserUUID  uuid.UUID `db:"user_uuid"`
	StartedOn time.Time `db:"started_on"`
	ExpiresOn time.Time `db:"expires_on"`
}

func insertSessionModel(context context.Context, transaction *database.Transaction, sessionUserUUID uuid.UUID) (*sessionModel, error) {
	_, errInsert := transaction.Execute(context, "INSERT INTO `sessions` (`user_uuid`) VALUE (?)", sessionUserUUID)
	if errInsert != nil {
		return nil, errInsert
	}

	model := new(sessionModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `sessions` WHERE `user_uuid` = ?", sessionUserUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func getSessionModelByToken(context context.Context, transaction *database.Transaction, sessionToken uuid.UUID) (*sessionModel, error) {
	model := new(sessionModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `sessions` WHERE `token` = ?", sessionToken)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func getSessionModelByUserUUID(context context.Context, transaction *database.Transaction, sessionUserUUID uuid.UUID) (*sessionModel, error) {
	model := new(sessionModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `sessions` WHERE `user_uuid` = ?", sessionUserUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func (model *sessionModel) expired() bool {
	return model.ExpiresOn.Before(time.Now())
}

func (model *sessionModel) updateExpiresOn(context context.Context, transaction *database.Transaction) error {
	_, errUpdate := transaction.Execute(context, "UPDATE `sessions` SET `expires_on` = DATE_ADD(NOW(), INTERVAL 20 MINUTE) WHERE `token` = ?", model.Token)
	if errUpdate != nil {
		return errUpdate
	}

	errGetExpiresOn := transaction.Get(context, &model.ExpiresOn, "SELECT `expires_on` FROM `sessions` WHERE `token` = ?", model.Token)
	if errGetExpiresOn != nil {
		return errGetExpiresOn
	}

	return nil
}

func (model *sessionModel) delete(context context.Context, transaction *database.Transaction) error {
	_, errDelete := transaction.Execute(context, "DELETE FROM `sessions` WHERE `token` = ?", model.Token)
	if errDelete != nil {
		return errDelete
	}

	return nil
}

type Session struct {
	model *sessionModel
	valid bool
}

var (
	ErrSessionInvalid error = errors.New("session invalid")
	ErrSessionExpired error = errors.New("session expired")
)

func newSession(context context.Context, transaction *database.Transaction, sessionUser *User) (*Session, error) {
	session := new(Session)

	var errInsertModel error
	session.model, errInsertModel = insertSessionModel(context, transaction, sessionUser.model.UUID)
	if errInsertModel != nil {
		return nil, errInsertModel
	}

	session.valid = true

	return session, nil
}

func getSessionByToken(context context.Context, transaction *database.Transaction, sessionToken uuid.UUID) (*Session, error) {
	session := new(Session)

	var errGetModel error
	session.model, errGetModel = getSessionModelByToken(context, transaction, sessionToken)
	if errGetModel != nil {
		return nil, errGetModel
	}

	session.valid = true

	return session, nil
}

func getSessionByUser(context context.Context, transaction *database.Transaction, sessionUser *User) (*Session, error) {
	session := new(Session)

	var errGetModel error
	session.model, errGetModel = getSessionModelByUserUUID(context, transaction, sessionUser.model.UUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	session.valid = true

	return session, nil
}

func beginSession(context context.Context, transaction *database.Transaction, sessionUser *User) (*Session, error) {
	session, errGet := getSessionByUser(context, transaction, sessionUser)
	if errGet == nil {
		if !session.expired() {
			errRefresh := session.refresh(context, transaction)
			if errRefresh != nil {
				return nil, errRefresh
			}

			return session, nil
		} else {
			errEnd := session.end(context, transaction)
			if errEnd != nil {
				return nil, errEnd
			}
		}
	}

	session, errNew := newSession(context, transaction, sessionUser)
	if errNew != nil {
		return nil, errNew
	}

	return session, nil
}

func continueSession(context context.Context, transaction *database.Transaction, sessionToken uuid.UUID) (*Session, error) {
	session, errGet := getSessionByToken(context, transaction, sessionToken)
	if errGet != nil {
		return nil, errGet
	}

	if session.expired() {
		errEnd := session.end(context, transaction)
		if errEnd != nil {
			return nil, errors.Join(ErrSessionExpired, errEnd)
		}

		return nil, ErrSessionExpired
	}

	errRefresh := session.refresh(context, transaction)
	if errRefresh != nil {
		return nil, errRefresh
	}

	return session, nil
}

func (session *Session) expired() bool {
	return session.model.expired()
}

func (session *Session) refresh(context context.Context, transaction *database.Transaction) error {
	errUpdateModel := session.model.updateExpiresOn(context, transaction)
	if errUpdateModel != nil {
		return errUpdateModel
	}

	return nil
}

func (session *Session) end(context context.Context, transaction *database.Transaction) error {
	errDeleteModel := session.model.delete(context, transaction)
	if errDeleteModel != nil {
		return errDeleteModel
	}

	session.valid = false

	return nil
}

func (session *Session) MarshalJSON() ([]byte, error) {
	if !session.valid {
		panic(ErrSessionInvalid)
	}

	return json.Marshal(map[string]any{
		"token":     session.model.Token,
		"startedOn": session.model.StartedOn,
		"expiresOn": session.model.ExpiresOn,
	})
}
