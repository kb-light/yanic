package database

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

// DB interface to use for implementation in e.g. influxdb
type DB interface {
	// AddNode data for a single node
	AddNode(nodeID string, node *runtime.Node)
	AddStatistics(stats *runtime.GlobalStats, time time.Time)

	DeleteNode(deleteAfter time.Duration)

	Close()
}

// New function with config to get DB connection interface
type New func(config map[string]interface{}) DB
