package meshviewer

import (
	"log"
	"time"

	"github.com/FreifunkBremen/yanic/jsontime"
	"github.com/FreifunkBremen/yanic/state"
)

// NodesV1 struct, to support legacy meshviewer (which are in master branch)
//  i.e. https://github.com/ffnord/meshviewer/tree/master
type NodesV1 struct {
	Version   int              `json:"version"`
	Timestamp jsontime.Time    `json:"timestamp"` // Timestamp of the generation
	List      map[string]*Node `json:"nodes"`     // the current nodemap, indexed by node ID
}

// NodesV2 struct, to support new version of meshviewer (which are in legacy develop branch or newer)
//  i.e. https://github.com/ffnord/meshviewer/tree/dev or https://github.com/ffrgb/meshviewer/tree/develop
type NodesV2 struct {
	Version   int           `json:"version"`
	Timestamp jsontime.Time `json:"timestamp"` // Timestamp of the generation
	List      []*Node       `json:"nodes"`     // the current nodemap, as array
}

// GetNodesV1 transform data to legacy meshviewer
func GetNodesV1(nodes *state.Nodes) *NodesV1 {
	meshviewerNodes := &NodesV1{
		Version:   1,
		List:      make(map[string]*Node),
		Timestamp: jsontime.Now(),
	}

	for nodeID := range nodes.List {
		nodeOrigin := nodes.List[nodeID]

		if nodeOrigin.Statistics == nil {
			continue
		}

		node := &Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags:     nodeOrigin.Flags,
			Nodeinfo:  nodeOrigin.Nodeinfo,
		}
		node.Statistics = NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List[nodeID] = node
	}
	return meshviewerNodes
}

// GetNodesV2 transform data to modern meshviewers
func GetNodesV2(nodes *state.Nodes) *NodesV2 {
	meshviewerNodes := &NodesV2{
		Version:   2,
		Timestamp: jsontime.Now(),
	}

	for nodeID := range nodes.List {
		nodeOrigin := nodes.List[nodeID]
		if nodeOrigin.Statistics == nil {
			continue
		}
		node := &Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags:     nodeOrigin.Flags,
			Nodeinfo:  nodeOrigin.Nodeinfo,
		}
		node.Statistics = NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List = append(meshviewerNodes.List, node)
	}
	return meshviewerNodes
}

// Start all services to manage Nodes
func Start(config *state.Config, nodes *state.Nodes) {
	go worker(config, nodes)
}

// Periodically saves the cached DB to json file
func worker(config *state.Config, nodes *state.Nodes) {
	c := time.Tick(config.Nodes.SaveInterval.Duration)

	for range c {
		saveMeshviewer(config, nodes)
	}
}

func saveMeshviewer(config *state.Config, nodes *state.Nodes) {
	// Locking foo
	nodes.RLock()
	defer nodes.RUnlock()
	if path := config.Meshviewer.NodesPath; path != "" {
		version := config.Meshviewer.Version
		switch version {
		case 1:
			state.SaveJSON(GetNodesV1(nodes), path)
		case 2:
			state.SaveJSON(GetNodesV2(nodes), path)
		default:
			log.Panicf("invalid nodes version: %d", version)
		}
	}

	if path := config.Meshviewer.GraphPath; path != "" {
		state.SaveJSON(BuildGraph(nodes), path)
	}
}
