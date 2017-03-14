package database

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

const (
	CounterFirmware = "firmware" // Measurement for firmware statistics
	CounterModel    = "model"    // Measurement for model statistics
)

// DB interface to use for implementation in e.g. influxdb
type DB interface {
	// AddNode data for a single node
	AddNode(nodeID string, node *runtime.Node)
	AddCounterFirmware(m runtime.CounterMap)
	AddCounterModel(m runtime.CounterMap)
	AddGlobal(stats *runtime.GlobalStats, time time.Time)

	DeleteNode()

	Close()
}

// New function with config to get DB connection interface
type New func(config *runtime.Config) DB
