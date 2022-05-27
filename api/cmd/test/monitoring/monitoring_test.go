package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VamaSingapore/vama-api/cmd/test"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

type MonitoringTests struct {
	app *fiber.App
	db  *pgxpool.Pool
}

func (m *MonitoringTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Monitoring tests")
	testApp := test.StartTestServer()
	m.app = testApp.App
	m.db = testApp.Db
}
func (m *MonitoringTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Monitoring tests")
	m.app.Shutdown()
	m.db.Close()
}
func (m *MonitoringTests) BeforeEach(t *testing.T) {}
func (m *MonitoringTests) AfterEach(t *testing.T)  {}

// Tests that /monitoring/health/db ping returns 200
func (m *MonitoringTests) SubTestPingDB(t *testing.T) {
	resp, pingErr := m.app.Test(httptest.NewRequest("GET", "/monitoring/health/db", nil))
	require.Nil(t, pingErr)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMonitoring(t *testing.T) {
	gtest.RunSubTests(t, &MonitoringTests{})
}
