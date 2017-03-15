package influxdb

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"

	"github.com/FreifunkBremen/yanic/state"
)

const (
	MeasurementNode   = "node"   // Measurement for per-node statistics
	MeasurementGlobal = "global" // Measurement for summarized global statistics
	batchMaxSize      = 500
	batchTimeout      = 5 * time.Second
)

type DB struct {
	config *state.Config
	client client.Client
	points chan *client.Point
	wg     sync.WaitGroup
	quit   chan struct{}
}

func New(config *state.Config) *DB {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Influxdb.Address,
		Username: config.Influxdb.Username,
		Password: config.Influxdb.Password,
	})

	if err != nil {
		panic(err)
	}

	db := &DB{
		config: config,
		client: c,
		points: make(chan *client.Point, 1000),
		quit:   make(chan struct{}),
	}

	db.wg.Add(1)
	go db.addWorker()
	go db.deleteWorker()

	return db
}

func (db *DB) DeleteNode() {
	query := fmt.Sprintf("delete from %s where time < now() - %ds", MeasurementNode, db.config.Influxdb.DeleteAfter.Duration/time.Second)
	db.client.Query(client.NewQuery(query, db.config.Influxdb.Database, "m"))
}

func (db *DB) AddPoint(name string, tags models.Tags, fields models.Fields, time time.Time) {
	point, err := client.NewPoint(name, tags.Map(), fields, time)
	if err != nil {
		panic(err)
	}
	db.points <- point
}

// Saves the values of a CounterMap in the database.
// The key are used as 'value' tag.
// The value is used as 'counter' field.
func (db *DB) AddCounterMap(name string, m state.CounterMap) {
	now := time.Now()
	for key, count := range m {
		db.AddPoint(
			name,
			models.Tags{
				models.Tag{Key: []byte("value"), Value: []byte(key)},
			},
			models.Fields{"count": count},
			now,
		)
	}
}

// AddGlobal implementation of database
func (db *DB) AddGlobal(stats *state.GlobalStats, time time.Time) {
	db.AddPoint(MeasurementGlobal, nil, GlobalStatsFields(stats), time)
}

// AddNode implementation of database
func (db *DB) AddNode(nodeID string, node *state.Node) {
	tags, fields := NodeToInflux(node)
	db.AddPoint(MeasurementNode, tags, fields, time.Now())
}

// Close all connection and clean up
func (db *DB) Close() {
	close(db.quit)
	close(db.points)
	db.wg.Wait()
	db.client.Close()
}

// prunes node-specific data periodically
func (db *DB) deleteWorker() {
	ticker := time.NewTicker(db.config.Influxdb.DeleteInterval.Duration)
	for {
		select {
		case <-ticker.C:
			db.DeleteNode()
		case <-db.quit:
			ticker.Stop()
			return
		}
	}
}

// stores data points in batches into the influxdb
func (db *DB) addWorker() {
	bpConfig := client.BatchPointsConfig{
		Database:  db.config.Influxdb.Database,
		Precision: "m",
	}

	var bp client.BatchPoints
	var err error
	var writeNow, closed bool
	timer := time.NewTimer(batchTimeout)

	for !closed {
		// wait for new points
		select {
		case point, ok := <-db.points:
			if ok {
				if bp == nil {
					// create new batch
					timer.Reset(batchTimeout)
					if bp, err = client.NewBatchPoints(bpConfig); err != nil {
						log.Fatal(err)
					}
				}
				bp.AddPoint(point)
			} else {
				closed = true
			}
		case <-timer.C:
			if bp == nil {
				timer.Reset(batchTimeout)
			} else {
				writeNow = true
			}
		}

		// write batch now?
		if bp != nil && (writeNow || closed || len(bp.Points()) >= batchMaxSize) {
			log.Println("saving", len(bp.Points()), "points")

			if err = db.client.Write(bp); err != nil {
				log.Fatal(err)
			}
			writeNow = false
			bp = nil
		}
	}
	timer.Stop()
	db.wg.Done()
}
