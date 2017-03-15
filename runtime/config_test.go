package runtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	assert := assert.New(t)

	config := ReadConfigFile("../config_example.toml")
	assert.NotNil(config)

	assert.True(config.Respondd.Enable)
	assert.Equal("eth0", config.Respondd.Interface)
	assert.Equal(time.Minute, config.Respondd.CollectInterval.Duration)

	assert.Equal(2, config.Meshviewer.Version)
	assert.Equal("/var/www/html/meshviewer/data/nodes.json", config.Meshviewer.NodesPath)
	assert.Equal(time.Hour*24*7, config.Nodes.PruneAfter.Duration)
}
