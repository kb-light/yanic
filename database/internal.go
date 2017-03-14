package database

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

var quit chan struct{}

// Start workers of database
// WARNING: Do not override this function
//  you should use New()
func Start(db DB, config *runtime.Config) {
	quit = make(chan struct{})
	go deleteWorker(db, config.Database.DeleteInterval.Duration, config.Database.DeleteAfter.Duration)
}

func Close(db DB) {
	close(quit)
	db.Close()
}

// prunes node-specific data periodically
func deleteWorker(db DB, deleteInterval time.Duration, deleteAfter time.Duration) {
	ticker := time.NewTicker(deleteInterval)
	for {
		select {
		case <-ticker.C:
			db.DeleteNode(deleteAfter)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
