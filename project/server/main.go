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
	e.GET("/data", showData) // New endpoint to see sample data

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

	// For now, execute a sample SQL query
	// TODO: Replace with Cypher -> SQL transformation
	sampleSQL := "SELECT json_extract(properties, '$.name') as name, json_extract(properties, '$.age') as age FROM nodes WHERE label = 'Person'"

	columns, rows, err := executeSQL(sampleSQL)
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

// New endpoint to view sample data
func showData(c echo.Context) error {
	columns, rows, err := executeSQL("SELECT * FROM nodes")
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
