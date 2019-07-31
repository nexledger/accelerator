package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	server, err := New(configFIle)
	assert.NoError(t, err)

	failure := server.Serve()

	server.Stop()
	time.Sleep(time.Millisecond * 100)

	select {
	case <-failure:
	default:
		t.Fatalf("Should have closed")
	}
}
