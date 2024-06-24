package samuel

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
)

type companyModel struct {
	UUID    uuid.UUID      `db:"uuid"`
	Name    string         `db:"name"`
	Address string         `db:"address"`
	Unit    sql.NullString `db:"unit"`
	City    string         `db:"city"`
	State   string         `db:"state"`
	ZIP     string         `db:"zip"`
	Phone   string         `db:"phone"`
}

func getCompanyModelByUUID(context context.Context, transaction *database.Transaction, companyUUID uuid.UUID) (*companyModel, error) {
	model := new(companyModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `companies` WHERE `uuid` = ?", companyUUID)
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

type Company struct {
	model *companyModel
	valid bool
}

var (
	ErrCompanyInvalid error = errors.New("company invalid")
)

func getCompanyBySupervisor(context context.Context, transaction *database.Transaction, supervisor *Supervisor) (*Company, error) {
	company := new(Company)

	var errGetModel error
	company.model, errGetModel = getCompanyModelByUUID(context, transaction, supervisor.model.CompanyUUID)
	if errGetModel != nil {
		return nil, errGetModel
	}

	company.valid = true

	return company, nil
}

func (company *Company) MarshalJSON() ([]byte, error) {
	if !company.valid {
		panic(ErrCompanyInvalid)
	}

	companyMap := map[string]any{
		"name":    company.model.Name,
		"address": company.model.Address,
		"city":    company.model.City,
		"state":   company.model.State,
		"zip":     company.model.ZIP,
		"phone":   company.model.Phone,
	}
	if company.model.Unit.Valid {
		companyMap["unit"] = company.model.Unit.String
	}

	return json.Marshal(companyMap)
}
