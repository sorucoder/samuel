package database

import "database/sql"

type Result struct {
	raw sql.Result
}

func newResult(raw sql.Result) *Result {
	return &Result{
		raw: raw,
	}
}

func (result Result) LastInsertID() int64 {
	lastInsertID, errNotSupported := result.raw.LastInsertId()
	if errNotSupported != nil {
		panic(errNotSupported)
	}
	return lastInsertID
}

func (result Result) RowsAffected() int64 {
	rowsAffected, errNotSupported := result.raw.RowsAffected()
	if errNotSupported != nil {
		panic(errNotSupported)
	}
	return rowsAffected
}
