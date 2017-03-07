package debugdatabase

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/FreifunkBremen/yanic/state"
)

type DB struct {
	config *state.Config
	file   *os.File
}

func New(config *state.Config) *DB {
	file, err := os.OpenFile(config.Debug.File, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println("File could not opened: ", config.Debug.File)
		return nil
	}
	return &DB{config: config, file: file}
}
func (db *DB) AddNode(nodeID string, node *state.Node) {
	db.log("AddNode: [", nodeID, "] clients: ", node.Statistics.Clients.Total)
}

func (db *DB) AddGlobal(stats *state.GlobalStats, time time.Time) {
	db.log("AddGlobal: [", time.String(), "] nodes: ", stats.Nodes, ", clients: ", stats.Clients)
}

func (db *DB) AddCounterMap(name string, m state.CounterMap) {
	db.log("AddCounterMap: [", name, "] count: ", len(m))
}

func (db *DB) DeleteNode() {
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
