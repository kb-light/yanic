package exampledatabase

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

type DB struct {
	config *config
	file   *os.File
}
type config struct {
	Path string
}

func toConfig(configMap map[string]interface{}) *config {
	return &config{
		Path: configMap["path"].(string),
	}
}

func New(configMap map[string]interface{}) *DB {
	config := toConfig(configMap)
	file, err := os.OpenFile(config.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("File could not opened: ", config.Path)
		return nil
	}
	return &DB{config: config, file: file}
}
func (db *DB) AddNode(nodeID string, node *runtime.Node) {
	db.log("AddNode: [", nodeID, "] clients: ", node.Statistics.Clients.Total)
}

func (db *DB) AddStatistics(stats *runtime.GlobalStats, time time.Time) {
	db.log("AddStatistics: [", time.String(), "] nodes: ", stats.Nodes, ", clients: ", stats.Clients)
}

func (db *DB) DeleteNode(deleteAfter time.Duration) {
	db.log("DeleteNode")
}

func (db *DB) Close() {
	db.log("Close")
	db.file.Close()
}

func (db *DB) log(v ...interface{}) {
	log.Println(v)
	db.file.WriteString(fmt.Sprintln("[", time.Now().String(), "]", v))
}
