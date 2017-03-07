package influxdb

import (
	"github.com/FreifunkBremen/yanic/state"
)

// GlobalStatsFields returns fields for InfluxDB
func GlobalStatsFields(stats *state.GlobalStats) map[string]interface{} {
	return map[string]interface{}{
		"nodes":          stats.Nodes,
		"gateways":       stats.Gateways,
		"clients.total":  stats.Clients,
		"clients.wifi":   stats.ClientsWifi,
		"clients.wifi24": stats.ClientsWifi24,
		"clients.wifi5":  stats.ClientsWifi5,
	}
}
