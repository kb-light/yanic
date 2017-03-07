package state

import (
	"net"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/jsontime"
)

// Node struct
type Node struct {
	Address    net.IP           `json:"address"` // the last known IP address
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Flags      Flags            `json:"flags"`
	Statistics *data.Statistics `json:"statistics"`
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"-"`
}

// Flags status of node set by collector for the meshviewer
type Flags struct {
	Online  bool `json:"online"`
	Gateway bool `json:"gateway"`
}
