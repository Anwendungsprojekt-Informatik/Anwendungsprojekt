package main

import (
	"net/http"
	"strconv"

	"github.com/Anwendungsprojekt-Informatik/Anwendungsprojekt/api"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Erstelle den OpenFoodFacts Client
	client := api.NewOpenFoodFactClient()

	// 2) Initialisiere den Gin Router
	router := gin.Default()

	// ---------------------------------------------------
	// 3) API-Endpoints unter "/api"
	// ---------------------------------------------------
	apiGroup := router.Group("/api")
	{
		// a) Produktsuche: GET /api/search?q=<Suchbegriff>&page=<Seite>
		apiGroup.GET("/search", func(c *gin.Context) {
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

		// b) Barcode-Abfrage: GET /api/product/:barcode
		apiGroup.GET("/product/:barcode", func(c *gin.Context) {
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

	// ---------------------------------------------------
	// 4) Statisches Frontend unter "/public/*"
	//    -> Alle Dateien aus "../frontend/public" ausliefern
	// ---------------------------------------------------
	//
	//    Wenn im Browser "/public/index.html" angefragt wird,
	//    sucht Gin im Dateisystem nach "../frontend/public/index.html".
	//    Das gleiche gilt für "/public/add-product.html", "/public/info.html" etc.
	//
	//    WICHTIG: Das relative "../frontend/public" funktioniert nur,
	//    wenn du das Binary / "go run main.go" aus dem Ordner
	//    "SmartNutritionPlanner/backend" startest.
	//
	router.Static("/public", "../frontend/public")

	// ---------------------------------------------------
	// 5) Root-Route "/" → Redirect auf "/public/index.html"
	// ---------------------------------------------------
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/public/index.html")
	})

	// ---------------------------------------------------
	// 6) Starte den Server auf Port 8080
	// ---------------------------------------------------
	router.Run(":8080")
}
