package tests

import (
	web_bootstrap "github.com/gouef/web-bootstrap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := web_bootstrap.LoadConfig("./config/config.yml")

	assert.NoError(t, err)
	assert.Equal(t, "./views/templates", cfg.Renderer.Dir)
	assert.Equal(t, "lalala", cfg.Renderer.Custom["test"])
}
