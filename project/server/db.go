package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize database and create tables
func initDatabase(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create schema for storing graph data
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Insert sample data matching the Neo4j structure
	if err := insertNeo4jData(); err != nil {
		return fmt.Errorf("failed to insert sample data: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	// Create nodes table: stores graph nodes with labels and JSON properties
	nodeTable := `
    CREATE TABLE IF NOT EXISTS nodes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        label TEXT NOT NULL,
        properties TEXT DEFAULT '{}'
    )`

	// Create relationships table: stores graph relationships with types and JSON properties
	relTable := `
    CREATE TABLE IF NOT EXISTS relationships (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        from_id INTEGER NOT NULL,
        to_id INTEGER NOT NULL,
        type TEXT NOT NULL,
        properties TEXT DEFAULT '{}',
        FOREIGN KEY(from_id) REFERENCES nodes(id),
        FOREIGN KEY(to_id) REFERENCES nodes(id)
    )`

	if _, err := db.Exec(nodeTable); err != nil {
		return fmt.Errorf("failed to create nodes table: %w", err)
	}

	if _, err := db.Exec(relTable); err != nil {
		return fmt.Errorf("failed to create relationships table: %w", err)
	}

	return nil
}

// Insert data matching the Neo4j structure
func insertNeo4jData() error {
	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM nodes").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing data: %w", err)
	}

	if count > 0 {
		log.Println("Sample data already exists, skipping insertion")
		return nil // Data already exists
	}

	// Begin transaction for data consistency
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() succeeds

	// Insert People
	people := []map[string]interface{}{
		{"name": "Alice Johnson", "age": 28, "occupation": "Engineer"},
		{"name": "Bob Smith", "age": 35, "occupation": "Teacher"},
		{"name": "Carol Williams", "age": 42, "occupation": "Doctor"},
		{"name": "David Brown", "age": 31, "occupation": "Artist"},
		{"name": "Eve Davis", "age": 29, "occupation": "Writer"},
	}

	for _, person := range people {
		props, _ := json.Marshal(person)
		_, err := tx.Exec("INSERT INTO nodes (label, properties) VALUES (?, ?)",
			"Person", string(props))
		if err != nil {
			return fmt.Errorf("failed to insert person: %w", err)
		}
	}

	// Insert Cities
	cities := []map[string]interface{}{
		{"name": "Austin", "country": "USA", "population": 978908},
		{"name": "Dallas", "country": "USA", "population": 1343573},
		{"name": "Houston", "country": "USA", "population": 2320268},
		{"name": "San Antonio", "country": "USA", "population": 1547253},
	}

	for _, city := range cities {
		props, _ := json.Marshal(city)
		_, err := tx.Exec("INSERT INTO nodes (label, properties) VALUES (?, ?)",
			"City", string(props))
		if err != nil {
			return fmt.Errorf("failed to insert city: %w", err)
		}
	}

	// Insert Movies
	movies := []map[string]interface{}{
		{"title": "The Matrix", "year": 1999, "genre": "Sci-Fi", "rating": 8.7},
		{"title": "Inception", "year": 2010, "genre": "Sci-Fi", "rating": 8.8},
		{"title": "The Godfather", "year": 1972, "genre": "Crime", "rating": 9.2},
		{"title": "Pulp Fiction", "year": 1994, "genre": "Crime", "rating": 8.9},
		{"title": "The Shawshank Redemption", "year": 1994, "genre": "Drama", "rating": 9.3},
	}

	for _, movie := range movies {
		props, _ := json.Marshal(movie)
		_, err := tx.Exec("INSERT INTO nodes (label, properties) VALUES (?, ?)",
			"Movie", string(props))
		if err != nil {
			return fmt.Errorf("failed to insert movie: %w", err)
		}
	}

	// Insert Relationships
	relationships := []struct {
		fromID, toID int
		relType      string
		properties   map[string]interface{}
	}{
		// LIVES_IN relationships (People -> Cities)
		{1, 6, "LIVES_IN", map[string]interface{}{"since": 2020}}, // Alice -> Austin
		{2, 7, "LIVES_IN", map[string]interface{}{"since": 2018}}, // Bob -> Dallas
		{3, 8, "LIVES_IN", map[string]interface{}{"since": 2015}}, // Carol -> Houston
		{4, 6, "LIVES_IN", map[string]interface{}{"since": 2019}}, // David -> Austin
		{5, 9, "LIVES_IN", map[string]interface{}{"since": 2021}}, // Eve -> San Antonio

		// WATCHED relationships (People -> Movies)
		{1, 10, "WATCHED", map[string]interface{}{"rating": 9, "date": "2023-01-15"}},  // Alice -> Matrix
		{1, 11, "WATCHED", map[string]interface{}{"rating": 8, "date": "2023-02-10"}},  // Alice -> Inception
		{2, 12, "WATCHED", map[string]interface{}{"rating": 10, "date": "2023-01-20"}}, // Bob -> Godfather
		{3, 14, "WATCHED", map[string]interface{}{"rating": 9, "date": "2023-03-05"}},  // Carol -> Shawshank
		{4, 13, "WATCHED", map[string]interface{}{"rating": 8, "date": "2023-02-25"}},  // David -> Pulp Fiction

		// KNOWS relationships (People -> People)
		{1, 4, "KNOWS", map[string]interface{}{"since": 2019, "relationship": "friend"}},    // Alice -> David
		{2, 3, "KNOWS", map[string]interface{}{"since": 2020, "relationship": "colleague"}}, // Bob -> Carol
	}

	for _, rel := range relationships {
		props, _ := json.Marshal(rel.properties)
		_, err := tx.Exec("INSERT INTO relationships (from_id, to_id, type, properties) VALUES (?, ?, ?, ?)",
			rel.fromID, rel.toID, rel.relType, string(props))
		if err != nil {
			return fmt.Errorf("failed to insert relationship: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Neo4j-matching data inserted successfully")
	return nil
}

// Execute a SQL query and return results as columns and rows (for compatibility with main.go)
func executeSQL(sqlQuery string) ([]string, [][]string, error) {
	// Execute the SQL query against the database
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names from the result set
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Initialize slice to store all result rows
	var result [][]string

	// Iterate through each row returned by the query
	for rows.Next() {
		// Create slices to hold the raw column values for this row
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		// Set up pointers - Scan() needs memory addresses to write to
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Read the current row's data into our value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert raw database values to strings for JSON response
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				// Handle NULL database values
				row[i] = ""
			} else {
				// Convert different database types to strings
				switch v := val.(type) {
				case []byte:
					// TEXT/BLOB columns come as byte slices
					row[i] = string(v)
				case string:
					// String columns
					row[i] = v
				case int64:
					// INTEGER columns - use strconv for proper conversion
					row[i] = strconv.FormatInt(v, 10)
				case float64:
					// REAL columns
					row[i] = strconv.FormatFloat(v, 'f', -1, 64)
				case bool:
					// BOOLEAN columns
					row[i] = strconv.FormatBool(v)
				default:
					// Fallback: convert to string representation
					row[i] = fmt.Sprintf("%v", v)
				}
			}
		}
		// Add this processed row to our result set
		result = append(result, row)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("row iteration error: %w", err)
	}

	return columns, result, nil
}

// Alternative function that returns map results (from Rust version)
// This can be useful for more complex JSON responses
func executeSQLAsMap(sqlQuery string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create slice of interface{} for row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert to map
		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}

		results = append(results, row)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
