package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const configFIle = "../../testdata/accelerator_test.yaml"

func TestConfig(t *testing.T) {
	conf, err := loadConfig(configFIle)
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	client, err := conf.BatchClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
