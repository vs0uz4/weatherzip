package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUptimeService(t *testing.T) {
	uptimeService := NewUptimeService()

	assert.NotNil(t, uptimeService, "UptimeService should not be nil")
}

func TestGetUptime(t *testing.T) {
	uptimeService := NewUptimeService()
	time.Sleep(100 * time.Millisecond)

	uptime := uptimeService.GetUptime()

	assert.NotEmpty(t, uptime, "Uptime should not be empty")
}
