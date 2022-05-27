package test

import (
	"context"
	"fmt"
	"testing"

	contactrepo "github.com/VamaSingapore/vama-api/internal/entities/contact/repositories"

	"github.com/VamaSingapore/vama-api/cmd/test"
	test_util "github.com/VamaSingapore/vama-api/cmd/test/util"
	test_data "github.com/VamaSingapore/vama-api/cmd/test/util/data"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/gofiber/fiber/v2"
	"github.com/houqp/gtest"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ContactTests struct {
	app *fiber.App
	db  *pgxpool.Pool
}

func (m *ContactTests) Setup(t *testing.T) {
	fmt.Println("Starting test server for Contact tests")
	app := test.StartTestServer()
	m.app = app.App
	m.db = app.Db
	test_util.InitializeDbSchemas(m.db)
}

func (m *ContactTests) Teardown(t *testing.T) {
	fmt.Println("Killing test server for Contact tests")
	m.app.Shutdown()
	m.db.Close()
}

func (m *ContactTests) BeforeEach(t *testing.T) {
	test_util.ClearDB(m.db)
	test_data.FillBootstrapUserData(m.db)
}

func (m *ContactTests) AfterEach(t *testing.T) {
}

// Test adding user2 to user1's contacts returns a 200
// and user2's object in the response.
func (m *ContactTests) SubTestCreateContact(t *testing.T) {
	params := map[string]interface{}{
		"contactId": 2,
	}
	asUserID := 1
	resp := &response.User{}
	test_util.MakePostRequestAssert200(t, m.app, "/api/v1/contacts/me", params, &asUserID, resp)

	assert.Equal(t, 2, resp.ID)
}

// Test getting user1's contacts page 1 returns a 200
// and the correct pagination metadata
func (m *ContactTests) SubTestGetContacts(t *testing.T) {
	test_data.FillBootstrapContactData(m.db)
	params := map[string]string{
		"page": "1",
		"size": "10",
	}
	asUserID := 1
	resp := &response.GetContacts{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/contacts/me", params, asUserID, resp)

	require.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Contacts))
	assert.Equal(t, 1, resp.Paging.Page)
	assert.Equal(t, 10, resp.Paging.Size)
}

// Test getting user2's contacts page 1 returns a 200
// and the correct pagination metadata but excludes the
// added creator contact
func (m *ContactTests) SubTestGetContactsExcludeGoats(t *testing.T) {
	ctx := context.Background()
	userID := 2
	contactID := 1 // creator contact
	contactDBErr := contactrepo.CreateContactDB(ctx, userID, contactID, m.db)
	if contactDBErr != nil {
		vlog.Fatalf(ctx, "Error creating contact: %s", contactDBErr.Error())
	}

	params := map[string]string{
		"page":         "1",
		"size":         "10",
		"excludeGoats": "true",
	}
	asUserID := userID
	resp := &response.GetContacts{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/contacts/me", params, asUserID, resp)

	require.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Contacts))
	assert.Equal(t, 1, resp.Paging.Page)
	assert.Equal(t, 10, resp.Paging.Size)
}

// Test getting user2 from user1's contacts by id returns a 200
// and user2's object.
func (m *ContactTests) SubTestIsContact(t *testing.T) {
	test_data.FillBootstrapContactData(m.db)
	params := map[string]string{
		"contactID": "2",
	}
	asUserID := 1
	resp := &response.IsContact{}
	test_util.MakeGetRequestAssert200(t, m.app, "/api/v1/contacts/me/id", params, asUserID, resp)

	require.NotNil(t, resp)
	assert.Equal(t, true, resp.IsContact)
}

// Test deleting user2 from user1's contacts returns a 200
func (m *ContactTests) SubTestDeleteContact(t *testing.T) {
	test_data.FillBootstrapContactData(m.db)
	params := map[string]interface{}{
		"contactID": 2,
	}
	asUserID := 1
	resp := &response.DeleteContact{}
	test_util.MakeDeleteRequestAssert200(t, m.app, "/api/v1/contacts/me", params, asUserID, resp)

	require.NotNil(t, resp)
}

func TestContact(t *testing.T) {
	gtest.RunSubTests(t, &ContactTests{})
}
