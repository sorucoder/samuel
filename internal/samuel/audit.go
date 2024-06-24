package samuel

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/database"
)

type auditModel struct {
	ID          uint64    `db:"id"`
	Description string    `db:"description"`
	UserUUID    uuid.UUID `db:"user_uuid"`
	Timestamp   time.Time `db:"timestamp"`
}

func insertAuditModel(context context.Context, transaction *database.Transaction, auditDescription string, auditUserUUID uuid.UUID) (*auditModel, error) {
	result, errInsert := transaction.Execute(context, "INSERT INTO `audit` (`description`, `user_uuid`) VALUE (?, ?)", auditDescription, auditUserUUID)
	if errInsert != nil {
		return nil, errInsert
	}

	model := new(auditModel)

	errGet := transaction.Get(context, model, "SELECT * FROM `audit` WHERE `id` = ?", result.LastInsertID())
	if errGet != nil {
		return nil, errGet
	}

	return model, nil
}

func selectAuditModelsByDate(context context.Context, transaction *database.Transaction, auditDate time.Time, number int, limit int, sort string, descending bool) ([]*auditModel, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT * FROM `audit`")
	switch sort {
	case "description":
		if descending {
			queryBuilder.WriteString(" WHERE DATE(`timestamp`) = DATE(?) ORDER BY `description` DESC, `timestamp` DESC")
		} else {
			queryBuilder.WriteString(" WHERE DATE(`timestamp`) = DATE(?) ORDER BY `description` ASC, `timestamp` DESC")
		}
	case "user":
		if descending {
			queryBuilder.WriteString(" JOIN `users` ON `audit`.`user_uuid` = `user`.`uuid` JOIN `roles` ON `users`.`role_id` = `roles`.`id` WHERE DATE(`timestamp`) = DATE(?) ORDER BY `role`.`priority` DESC, `timestamp` DESC")
		} else {
			queryBuilder.WriteString(" JOIN `users` ON `audit`.`user_uuid` = `user`.`uuid` JOIN `roles` ON `users`.`role_id` = `roles`.`id` WHERE DATE(`timestamp`) = DATE(?) ORDER BY `role`.`priority` ASC, `timestamp` DESC")
		}
	default:
		queryBuilder.WriteString(" WHERE DATE(`timestamp`) = DATE(?) ORDER BY `timestamp` DESC")
	}
	queryBuilder.WriteString(" LIMIT ? OFFSET ?")

	offset := number * limit

	models := make([]*auditModel, 0, limit)

	errSelect := transaction.Select(context, &models, queryBuilder.String(), auditDate, limit, offset)
	if errSelect != nil {
		return nil, errSelect
	}

	return models, nil
}

func countAuditModelsOnDate(context context.Context, transaction *database.Transaction, auditDate time.Time) (int64, error) {
	var count int64

	errGet := transaction.Get(context, &count, "SELECT COUNT(*) FROM `audit` WHERE DATE(`timestamp`) = DATE(?)", auditDate)
	if errGet != nil {
		return 0, errGet
	}

	return count, nil
}

type Audit struct {
	model *auditModel
	valid bool
}

var (
	ErrAuditInvalid error = errors.New("audit invalid")
)

func recordAudit(context context.Context, transaction *database.Transaction, auditDescription string, actor *User) error {
	_, errInsertModel := insertAuditModel(context, transaction, auditDescription, actor.model.UUID)
	if errInsertModel != nil {
		return errInsertModel
	}

	return nil
}

func getAuditBatchByDate(context context.Context, transaction *database.Transaction, auditDate time.Time, page int, count int, sort string, descending bool) (*Batch[*Audit], error) {
	audits := make([]*Audit, 0, count)

	auditModelCount, errCountModels := countAuditModelsOnDate(context, transaction, auditDate)
	if errCountModels != nil {
		return nil, errCountModels
	}

	auditModels, errGetModels := selectAuditModelsByDate(context, transaction, auditDate, page, count, sort, descending)
	if errGetModels != nil {
		return nil, errGetModels
	}
	for _, auditModel := range auditModels {
		audits = append(audits, &Audit{
			model: auditModel,
			valid: true,
		})
	}

	return newBatch(page, count, auditModelCount, "audits", audits...), nil
}

func GetAuditBatchByDate(context context.Context, onDate time.Time, batchNumber int, batchSize int, sortColumn string, sortDescending bool) (*Batch[*Audit], error) {
	errPing := database.Ping(context)
	if errPing != nil {
		return nil, errPing
	}

	transaction, errBegin := database.Begin(context)
	if errBegin != nil {
		return nil, errBegin
	}

	auditBatch, errGetAuditBatch := getAuditBatchByDate(context, transaction, onDate, batchNumber, batchSize, sortColumn, sortDescending)
	if errGetAuditBatch != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			return nil, errors.Join(errGetAuditBatch, errRollback)
		}

		return nil, errGetAuditBatch
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		return nil, errCommit
	}

	return auditBatch, nil
}

func (audit *Audit) MarshalJSON() ([]byte, error) {
	if !audit.valid {
		panic(ErrAuditInvalid)
	}

	return json.Marshal(map[string]any{
		"description": audit.model.Description,
		"userUUID":    audit.model.UserUUID,
		"timestamp":   audit.model.Timestamp,
	})
}
