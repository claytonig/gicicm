package db

import (
	"database/sql"
	"fmt"

	"gicicm/config"
	"gicicm/logger"

	_ "github.com/lib/pq" //dialect to be used
	"go.uber.org/zap"
)

// Database is an adapter layer for the database layer.
type Database interface{}

// NewDatabaseAdapter - returns a new instance of the database adapter.
func NewDatabaseAdapter(c *config.Config) *sql.DB {

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Pass, c.Database.DBName)

	dbconn, err := sql.Open(c.Database.DBType, connectionString)

	if err != nil {
		logger.Log().Fatal("Unable to connect to database", zap.String("connectionString", connectionString), zap.Error(err))
	}

	return dbconn
}
