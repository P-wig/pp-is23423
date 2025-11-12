// filepath: project/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Success bool       `json:"success"`
	Columns []string   `json:"columns,omitempty"`
	Rows    [][]string `json:"rows,omitempty"`
	Error   string     `json:"error,omitempty"`
}

func main() {
	// Initialize database
	if err := initDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.POST("/query", handleQuery)
	e.GET("/health", healthCheck)
	e.GET("/data", showData)
	e.GET("/test-sql", testSQL) // New test endpoint

	// Start server
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

	// TODO: Replace with Cypher -> SQL transformation
	// For now, execute a sample SQL query based on common Cypher patterns
	var sqlQuery string
	if req.Query != "" {
		// Simple example: return all people for any query
		sqlQuery = `SELECT json_extract(properties, '$.name') as name, 
                           json_extract(properties, '$.age') as age,
                           json_extract(properties, '$.occupation') as occupation
                    FROM nodes WHERE label = 'Person'`
	} else {
		sqlQuery = "SELECT 'No query provided' as message"
	}

	columns, rows, err := executeSQL(sqlQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, QueryResponse{
			Success: false,
			Error:   "Database query failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, QueryResponse{
		Success: true,
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
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"columns": columns,
		"rows":    rows,
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
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"description": "People living in Austin",
		"columns":     columns,
		"rows":        rows,
	})
}
