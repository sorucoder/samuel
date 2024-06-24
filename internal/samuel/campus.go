package samuel

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sorucoder/samuel/internal/database"
)

type campusModel struct {
	ID      string         `db:"id"`
	Name    string         `db:"name"`
	Address string         `db:"address"`
	Unit    sql.NullString `db:"unit"`
	City    string         `db:"city"`
	State   string         `db:"state"`
	ZIP     string         `db:"zip"`
	Phone   string         `db:"phone"`
}

func getCampusModelByID(context context.Context, transaction *database.Transaction, campusID string) (*campusModel, error) {
	model := new(campusModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `campuses` WHERE `id` = ?", campusID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Campus struct {
	model *campusModel
	valid bool
}

var (
	ErrCampusInvalid error = errors.New("campus invalid")
)

func getCampusByInstructor(context context.Context, transaction *database.Transaction, instructor *Instructor) (*Campus, error) {
	campus := new(Campus)

	var errGetModel error
	campus.model, errGetModel = getCampusModelByID(context, transaction, instructor.model.CampusID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	campus.valid = true

	return campus, nil
}

func getCampusByStudent(context context.Context, transaction *database.Transaction, student *Student) (*Campus, error) {
	campus := new(Campus)

	var errGetModel error
	campus.model, errGetModel = getCampusModelByID(context, transaction, student.model.CampusID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	campus.valid = true

	return campus, nil
}

func (campus *Campus) MarshalJSON() ([]byte, error) {
	if !campus.valid {
		panic(ErrCampusInvalid)
	}

	campusMap := map[string]any{
		"name":    campus.model.Name,
		"address": campus.model.Address,
		"city":    campus.model.City,
		"state":   campus.model.State,
		"zip":     campus.model.ZIP,
		"phone":   campus.model.Phone,
	}
	if campus.model.Unit.Valid {
		campusMap["unit"] = campus.model.Unit.String
	}

	return json.Marshal(campusMap)
}
