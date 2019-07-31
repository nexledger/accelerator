package server

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const configFIle = "accelerator_test.yaml"

func TestConfig(t *testing.T) {
	configTestFilePath := filepath.Join("testdata", configFIle)

	conf, err := loadConfig(configTestFilePath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	client, err := conf.BatchClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
