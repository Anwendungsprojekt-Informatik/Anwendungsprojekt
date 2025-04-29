package main

import (
	"net/http"
	"strconv"

	"github.com/Anwendungsprojekt-Informatik/Anwendungsprojekt/Meilenstein2/SmartNutritionPlanner/backend/api"
	"github.com/gin-gonic/gin"
)

func main() {
	// Erstelle den OpenFoodFacts Client
	client := api.NewOpenFoodFactClient()

	// Initialisiere den Gin Router
	router := gin.Default()

	// Stelle statische Dateien zur Verfügung (für das HTML-Frontend)
	router.StaticFile("/", "./index.html")

	// API-Endpunkte
	api := router.Group("/api")
	{
		// Endpunkt für die Suche nach Produkten
		api.GET("/search", func(c *gin.Context) {
			query := c.Query("q")
			pageStr := c.DefaultQuery("page", "1")

			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				page = 1
			}

			if query == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
				return
			}

			results, err := client.SearchProducts(query, page)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, results)
		})

		// Endpunkt für die Suche nach einem Produkt anhand des Barcodes
		api.GET("/product/:barcode", func(c *gin.Context) {
			barcode := c.Param("barcode")

			if barcode == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Barcode is required"})
				return
			}

			product, err := client.GetProductByBarcode(barcode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, product)
		})
	}

	// Server starten
	router.Run(":8080")
}
