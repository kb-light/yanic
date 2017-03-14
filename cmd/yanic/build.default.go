// +build !exampledatabase

package main

import (
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/database/influxdb"
	"github.com/FreifunkBremen/yanic/runtime"
)

func connectDB(config *runtime.Config) (db database.DB) {
	db = influxdb.New(config.Database.Connection)
	database.Start(db, config)
	return
}
