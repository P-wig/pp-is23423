// filepath: project/server/main.go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type QueryRequest struct {
    Query string `json:"query"`
}

type QueryResponse struct {
    Success bool        `json:"success"`
    Columns []string    `json:"columns,omitempty"`
    Rows    [][]string  `json:"rows,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func main() {
    e := echo.New()
    
    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())
    
    // Routes
    e.POST("/query", handleQuery)
    e.GET("/health", healthCheck)
    
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
    
    // For now, return a placeholder response
    return c.JSON(http.StatusOK, QueryResponse{
        Success: true,
        Columns: []string{"name", "age"},
        Rows:    [][]string{{"Alice", "25"}, {"Bob", "30"}},
    })
}

func healthCheck(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{
        "status": "healthy",
    })
}