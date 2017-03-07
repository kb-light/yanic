package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/debugdatabase"
	"github.com/FreifunkBremen/yanic/influxdb"
	"github.com/FreifunkBremen/yanic/meshviewer"
	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/rrd"
	"github.com/FreifunkBremen/yanic/state"
	"github.com/FreifunkBremen/yanic/webserver"
)

var (
	configFile string
	config     *state.Config
	collector  *respond.Collector
	db         database.DB
	nodes      *state.Nodes
)

func main() {
	var importPath string
	var timestamps bool
	flag.StringVar(&importPath, "import", "", "import global statistics from the given RRD file, requires influxdb")
	flag.StringVar(&configFile, "config", "config.toml", "path of configuration file (default:config.yaml)")
	flag.BoolVar(&timestamps, "timestamps", true, "print timestamps in output")
	flag.Parse()

	if !timestamps {
		log.SetFlags(0)
	}

	config = state.ReadConfigFile(configFile)

	if INFLUXDB_BOOTSTRAP && config.Influxdb.Enable {
		db = influxdb.New(config)
		defer db.Close()

	}

	if DEBUGDATABASE_BOOTSTRAP && config.Debug.Enable {
		db = debugdatabase.New(config)
		defer db.Close()

	}

	if db != nil && importPath != "" {
		importRRD(importPath)
		return
	}

	nodes = state.NewNodes(config)
	nodes.Start()
	meshviewer.Start(config, nodes)

	if config.Respondd.Enable {
		collector = respond.NewCollector(db, nodes, config.Respondd.Interface)
		collector.Start(config.Respondd.CollectInterval.Duration)
		defer collector.Close()
	}

	if config.Webserver.Enable {
		log.Println("starting webserver on", config.Webserver.Bind)
		srv := webserver.New(config.Webserver.Bind, config.Webserver.Webroot)
		go srv.Close()
	}

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)
}

func importRRD(path string) {
	log.Println("importing RRD from", path)
	for ds := range rrd.Read(path) {
		db.AddGlobal(
			&state.GlobalStats{
				Nodes:   uint32(ds.Nodes),
				Clients: uint32(ds.Clients),
			},
			ds.Time,
		)
	}
}
