package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDatabase(t *testing.T) {
	// Remove any existing test database
	testDBPath := "test_graph.db"
	_ = os.Remove(testDBPath)

	// Override global db variable for testing
	origDB := db
	defer func() { db = origDB }()

	// Set up test database
	err := initDatabase(testDBPath)
	assert.NoError(t, err, "Database should initialize without error")
	assert.NotNil(t, db, "Database connection should not be nil")

	// Clean up
	db.Close()
	_ = os.Remove(testDBPath)
}

func TestExecuteSQL(t *testing.T) {
	// Setup test database
	testDBPath := "test_graph.db"
	_ = os.Remove(testDBPath)
	err := initDatabase(testDBPath)
	assert.NoError(t, err)

	// Insert a test node
	_, _, err = executeSQL(`INSERT INTO nodes (id, label, properties) VALUES (100, 'Person', '{"name":"Test User"}')`)
	assert.NoError(t, err)

	// Query the test node
	columns, rows, err := executeSQL(`SELECT label, properties FROM nodes WHERE id = 100`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"label", "properties"}, columns)
	assert.Equal(t, [][]string{{"Person", `{"name":"Test User"}`}}, rows)

	// Clean up
	db.Close()
	_ = os.Remove(testDBPath)
}
