package database

import (
	"time"

	"github.com/FreifunkBremen/yanic/state"
)

const (
	CounterMeasurementFirmware = "firmware" // Measurement for firmware statistics
	CounterMeasurementModel    = "model"    // Measurement for model statistics
)

// DB interface to use for implementation in e.g. influxdb
type DB interface {
	// AddNode data for a single node
	AddNode(nodeID string, node *state.Node)
	AddCounterMap(name string, m state.CounterMap)
	AddGlobal(stats *state.GlobalStats, time time.Time)

	DeleteNode()

	Close()
}

// New function with config to get DB connection interface
type New func(config *state.Config) DB
