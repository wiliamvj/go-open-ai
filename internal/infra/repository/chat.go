package repository

import (
	"database/sql"

	"github.com/wiliamvj/go-open-ai/internal/infra/db"
)

type ChatRepositoryMySQL struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewChatRepositoryMySQL(dbt *sql.DB) *ChatRepositoryMySQL {
	return &ChatRepositoryMySQL{
		DB:      dbt,
		Queries: db.New(dbt),
	}
}
