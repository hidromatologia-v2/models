package connection

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/random"
	"github.com/stretchr/testify/assert"
)

func TestPostgresController(t *testing.T) {
	c := PostgresController()
	assert.NotNil(t, c)
	defer c.Close()
}

func TestDefaultConsumer(t *testing.T) {
	c := DefaultConsumer(t)
	assert.NotNil(t, c)
	defer c.Destroy()
}

func TestDefaultProducer(t *testing.T) {
	c := DefaultProducer(t)
	assert.NotNil(t, c)
	defer c.Destroy()
}

func TestNewConsumer(t *testing.T) {
	c := NewConsumer(t, random.String(), random.String())
	assert.NotNil(t, c)
	defer c.Destroy()
}

func TestNewProducer(t *testing.T) {
	c := NewProducer(t, random.String(), random.String())
	assert.NotNil(t, c)
	defer c.Destroy()
}
