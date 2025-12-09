// filepath: project/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Success bool       `json:"success"`
	Cypher  string     `json:"cypher,omitempty"`
	SQL     string     `json:"sql,omitempty"`
	Columns []string   `json:"columns,omitempty"`
	Rows    [][]string `json:"rows,omitempty"`
	Error   string     `json:"error,omitempty"`
}

func main() {
	// Initialize database
	if err := initDatabase("graph.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Ensure transformer binary is built
	if err := buildTransformer(); err != nil {
		log.Fatal("Failed to build transformer:", err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.POST("/query", handleQuery) // Original Cypher query endpoint
	e.GET("/health", healthCheck) // Health check
	e.GET("/data", showData)      // View sample data
	e.GET("/test-sql", testSQL)   // Test SQL queries directly

	// Start server
	log.Println("Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func handleQuery(c echo.Context) error {
	var req QueryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, QueryResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	if req.Query == "" {
		return c.JSON(http.StatusBadRequest, QueryResponse{
			Success: false,
			Error:   "No Cypher query provided",
		})
	}

	// Transform Cypher to SQL using your Rust transformer
	sqlQuery, err := transformCypherToSQL(req.Query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, QueryResponse{
			Success: false,
			Cypher:  req.Query,
			Error:   "Failed to transform Cypher query: " + err.Error(),
		})
	}

	// Execute the transformed SQL
	columns, rows, err := executeSQL(sqlQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, QueryResponse{
			Success: false,
			Cypher:  req.Query,
			SQL:     sqlQuery,
			Error:   "Database query failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, QueryResponse{
		Success: true,
		Cypher:  req.Query,
		SQL:     sqlQuery,
		Columns: columns,
		Rows:    rows,
	})
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// View all nodes
func showData(c echo.Context) error {
	columns, rows, err := executeSQL("SELECT * FROM nodes LIMIT 10")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, QueryResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, QueryResponse{
		Success: true,
		Columns: columns,
		Rows:    rows,
	})
}

// Test SQL queries directly
func testSQL(c echo.Context) error {
	// Test query to get people living in Austin
	sql := `SELECT 
                n1.id,
                json_extract(n1.properties, '$.name') as person_name,
                json_extract(n1.properties, '$.occupation') as occupation,
                json_extract(n2.properties, '$.name') as city_name
            FROM nodes n1
            JOIN relationships r ON n1.id = r.from_id
            JOIN nodes n2 ON r.to_id = n2.id
            WHERE n1.label = 'Person' 
              AND n2.label = 'City' 
              AND r.type = 'LIVES_IN'
              AND json_extract(n2.properties, '$.name') = 'Austin'`

	columns, rows, err := executeSQL(sql)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, QueryResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, QueryResponse{
		Success: true,
		Columns: columns,
		Rows:    rows,
		SQL:     sql,
	})
}

func buildTransformer() error {
	// Check if binary already exists
	binaryPath := "../cypher-transformer/target/release/cypher_transformer.exe"
	if _, err := os.Stat(binaryPath); err == nil {
		log.Println("Transformer binary already exists")
		return nil // Already built
	}

	log.Println("Building Rust transformer...")
	cmd := exec.Command("cargo", "build", "--release")
	cmd.Dir = "../cypher-transformer"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build transformer: %w", err)
	}

	log.Println("Transformer built successfully")
	return nil
}

func transformCypherToSQL(cypherQuery string) (string, error) {
	binaryPath := "../cypher-transformer/target/release/cypher_transformer.exe"

	// Double-check binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return "", fmt.Errorf("transformer binary not found at %s", binaryPath)
	}

	cmd := exec.Command(binaryPath, cypherQuery)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("transformation failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
