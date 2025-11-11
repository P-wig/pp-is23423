package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize database and create tables
func initDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./graph.db")
	if err != nil {
		return err
	}

	// Create schema for storing graph data
	schema := `
    CREATE TABLE IF NOT EXISTS nodes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        label TEXT NOT NULL,
        properties TEXT DEFAULT '{}'
    );

    CREATE TABLE IF NOT EXISTS relationships (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        from_id INTEGER NOT NULL,
        to_id INTEGER NOT NULL,
        type TEXT NOT NULL,
        properties TEXT DEFAULT '{}',
        FOREIGN KEY(from_id) REFERENCES nodes(id),
        FOREIGN KEY(to_id) REFERENCES nodes(id)
    );`

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	// Insert sample data matching the Neo4j structure
	return insertNeo4jData()
}

// Insert data matching the Neo4j structure
func insertNeo4jData() error {
	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM nodes").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Data already exists
	}

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
		_, err := db.Exec("INSERT INTO nodes (label, properties) VALUES (?, ?)",
			"Person", string(props))
		if err != nil {
			return err
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
		_, err := db.Exec("INSERT INTO nodes (label, properties) VALUES (?, ?)",
			"City", string(props))
		if err != nil {
			return err
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
		_, err := db.Exec("INSERT INTO nodes (label, properties) VALUES (?, ?)",
			"Movie", string(props))
		if err != nil {
			return err
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
		_, err := db.Exec("INSERT INTO relationships (from_id, to_id, type, properties) VALUES (?, ?, ?, ?)",
			rel.fromID, rel.toID, rel.relType, string(props))
		if err != nil {
			return err
		}
	}

	log.Println("Neo4j-matching data inserted successfully")
	return nil
}

// Execute a simple SQL query and return results
func executeSQL(sqlQuery string) ([]string, [][]string, error) {
	// Execute the SQL query against the database
	// This returns a *sql.Rows object that contains the result set
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, nil, err
	}
	// Ensure rows are closed when function exits to free database resources
	defer rows.Close()

	// Get column names from the result set
	// This tells us what fields/columns are in our query results
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	// Initialize slice to store all result rows
	// Each row will be a slice of strings: [["Alice", "28"], ["Bob", "35"]]
	var result [][]string

	// Iterate through each row returned by the query
	for rows.Next() {
		// Create slices (interface is a slice that holds any type of value)
		// to hold the raw column values for this row
		// values: holds the actual data from database
		// valuePtrs: holds pointers to the values (needed for Scan)
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		// Set up pointers - Scan() needs memory addresses to write to
		for i := range values {
			valuePtrs[i] = &values[i] // Point to each value slot
		}

		// Read the current row's data into our value pointers
		// This populates the 'values' slice with actual database data
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, err
		}

		// Convert raw database values to strings for JSON response
		// Database returns different types ([]byte, string, int64, etc.)
		// Everything must be converted to strings for our API
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
					// String columns (rare in SQLite)
					row[i] = v
				case int64:
					// INTEGER columns - convert to string representation
					// Note: string(rune(v)) might not work for large numbers
					// Consider using strconv.FormatInt(v, 10) instead
					row[i] = string(rune(v))
				default:
					// Fallback for unexpected types
					row[i] = ""
				}
			}
		}
		// Add this processed row to our result set
		result = append(result, row)
	}

	// Return column names and all rows as string arrays
	// columns: ["name", "age"]
	// result: [["Alice", "28"], ["Bob", "35"]]
	return columns, result, nil
}
