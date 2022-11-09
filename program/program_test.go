package program

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestOptions_Run(t *testing.T) {
//	var program Options
//
//	exitValue := -1
//	fakeExit := func(x int) {
//		exitValue = x
//	}
//	patch := monkey.Patch(os.Exit, fakeExit)
//	defer patch.Unpatch()
//
//	out := capturer.CaptureStdout(func() {
//
//		_, err := program.Parse([]string{"--version"})
//
//		assert.NoError(t, err)
//
//		// version output is done as part of parsing, so we don't need to run the program
//	})
//
//	assert.Equal(t, exitValue, 0)
//	assert.Equal(t, "unknown\n", out)
//}

func TestOptions_ReadConfig(t *testing.T) {
	program := Options{KubeConfig: "testdata/config"}

	config, err := program.ReadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, 6, len(config.Clusters))

	c := config.Clusters["docker-desktop"]
	assert.Equal(t, "https://kubernetes.docker.internal:6443", c.Server)
}
