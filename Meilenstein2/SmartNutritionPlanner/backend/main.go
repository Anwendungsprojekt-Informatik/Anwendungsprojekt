// main.go
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Anwendungsprojekt-Informatik/Anwendungsprojekt/api"
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
	apiGroup := router.Group("/api")
	{
		// Endpunkt für die Suche nach Produkten
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

		// Endpunkt für die Suche nach einem Produkt anhand des Barcodes
		apiGroup.GET("/product/:barcode", func(c *gin.Context) {
			barcode := c.Param("barcode")
			if barcode == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Barcode is required"})
				return
			}

			// Produkt aus API holen
			p, err := client.GetProductByBarcode(barcode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Nährwert-Übersicht auf der Konsole ausgeben
			printNutritionSummary(p)
			// Zusätzlich berechnete Energie aus Makros
			total := calculateEnergyFromMacros(p.Proteins100g, p.Carbohydrates100g, p.Fat100g)
			fmt.Printf("Berechnete Energie aus Makros: %.2f kcal\n", total)

			// JSON-Antwort für das Frontend
			c.JSON(http.StatusOK, p)
		})
	}

	// Server starten
	router.Run(":8080")
}

// printNutritionSummary gibt die Nährwerte (pro 100 g) formatiert auf der Konsole aus
func printNutritionSummary(p *api.Product) {
	fmt.Println("–– Nährwert-Übersicht (pro 100 g) ––")
	fmt.Printf("Energie:        %.2f kcal\n", p.EnergyKcal100g)
	fmt.Printf("Kohlenhydrate:  %.2f g\n", p.Carbohydrates100g)
	fmt.Printf("– davon Zucker: %.2f g\n", p.Sugars100g)
	fmt.Printf("Proteine:       %.2f g\n", p.Proteins100g)
	fmt.Printf("Fett:           %.2f g\n", p.Fat100g)
}

// calculateEnergyFromMacros berechnet die Energie aus Protein, Kohlenhydraten und Fett
// (Atwater-Faktoren: 4 kcal/g für Protein & Kohlenhydrate, 9 kcal/g für Fett)
func calculateEnergyFromMacros(protein, carbs, fat float64) float64 {
	return protein*4 + carbs*4 + fat*9
}
