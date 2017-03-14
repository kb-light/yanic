// +build exampledatabase

package main

import (
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/database/exampledatabase"
	"github.com/FreifunkBremen/yanic/runtime"
)

func connectDB(config *runtime.Config) (db database.DB) {
	db = exampledatabase.New(config.Database.Connection)
	database.Start(db, config)
	return
}
