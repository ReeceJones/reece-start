package utils

import (
	"database/sql"

	"gorm.io/gorm"
)

// GetGormSQLTx extracts the underlying SQL transaction from a GORM transaction
func GetGormSQLTx(tx *gorm.DB) *sql.Tx {
	sqlTx := tx.Statement.ConnPool.(*sql.Tx)
	return sqlTx
}
