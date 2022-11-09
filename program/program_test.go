package program

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions_ReadConfig(t *testing.T) {
	program := Options{KubeConfig: "testdata/config"}

	config, err := program.ReadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, 6, len(config.Clusters))

	c := config.Clusters["docker-desktop"]
	assert.Equal(t, "https://kubernetes.docker.internal:6443", c.Server)
}
