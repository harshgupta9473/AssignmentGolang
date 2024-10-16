package config

import (
	"database/sql"

	"github.com/harshgupta9473/recruitmentManagement/database"
)

func GetDB() *sql.DB{
	return database.DB
}