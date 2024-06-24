package samuel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/sorucoder/samuel/internal/database"
)

type roleModel struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Priority uint8  `db:"priority"`
}

func getRoleModelByID(context context.Context, transaction *database.Transaction, roleID string) (*roleModel, error) {
	model := new(roleModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `roles` WHERE `id` = ?", roleID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Role struct {
	model *roleModel
	valid bool
}

var (
	ErrInvalidRole error = errors.New("invalid role")
)

func getRoleByID(context context.Context, transaction *database.Transaction, roleID string) (*Role, error) {
	role := new(Role)

	var errGetModel error
	role.model, errGetModel = getRoleModelByID(context, transaction, roleID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	role.valid = true

	return role, nil
}

func (role *Role) ID() string {
	if !role.valid {
		panic(ErrInvalidRole)
	}

	return role.model.ID
}

func (role *Role) MarshalJSON() ([]byte, error) {
	if !role.valid {
		panic(ErrInvalidRole)
	}

	return json.Marshal(map[string]any{
		"id":   role.model.ID,
		"name": role.model.Name,
	})
}
