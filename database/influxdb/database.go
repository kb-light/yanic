package influxdb

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"

	"github.com/FreifunkBremen/yanic/runtime"
)

const (
	MeasurementNode            = "node"     // Measurement for per-node statistics
	MeasurementGlobal          = "global"   // Measurement for summarized global statistics
	CounterMeasurementFirmware = "firmware" // Measurement for firmware statistics
	CounterMeasurementModel    = "model"    // Measurement for model statistics
	batchMaxSize               = 500
	batchTimeout               = 5 * time.Second
)

type DB struct {
	config *config
	client client.Client
	points chan *client.Point
	wg     sync.WaitGroup
}

func New(configMap map[string]interface{}) *DB {
	config := toConfig(configMap)
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
	})

	if err != nil {
		panic(err)
	}

	db := &DB{
		config: config,
		client: c,
		points: make(chan *client.Point, 1000),
	}

	db.wg.Add(1)
	go db.addWorker()

	return db
}

func (db *DB) DeleteNode(deleteAfter time.Duration) {
	query := fmt.Sprintf("delete from %s where time < now() - %ds", MeasurementNode, deleteAfter/time.Second)
	db.client.Query(client.NewQuery(query, db.config.Database, "m"))
}

func (db *DB) addPoint(name string, tags models.Tags, fields models.Fields, time time.Time) {
	point, err := client.NewPoint(name, tags.Map(), fields, time)
	if err != nil {
		panic(err)
	}
	db.points <- point
}

// Saves the values of a CounterMap in the database.
// The key are used as 'value' tag.
// The value is used as 'counter' field.
func (db *DB) addCounterMap(name string, m runtime.CounterMap) {
	now := time.Now()
	for key, count := range m {
		db.addPoint(
			name,
			models.Tags{
				models.Tag{Key: []byte("value"), Value: []byte(key)},
			},
			models.Fields{"count": count},
			now,
		)
	}
}

// AddStatistics implementation of database
func (db *DB) AddStatistics(stats *runtime.GlobalStats, time time.Time) {
	db.addPoint(MeasurementGlobal, nil, GlobalStatsFields(stats), time)
	db.addCounterMap(CounterMeasurementModel, stats.Models)
	db.addCounterMap(CounterMeasurementFirmware, stats.Firmwares)
}

// AddNode implementation of database
func (db *DB) AddNode(nodeID string, node *runtime.Node) {
	tags, fields := nodeToInflux(node)
	db.addPoint(MeasurementNode, tags, fields, time.Now())
}

// Close all connection and clean up
func (db *DB) Close() {
	close(db.points)
	db.wg.Wait()
	db.client.Close()
}

// stores data points in batches into the influxdb
func (db *DB) addWorker() {
	bpConfig := client.BatchPointsConfig{
		Database:  db.config.Database,
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
