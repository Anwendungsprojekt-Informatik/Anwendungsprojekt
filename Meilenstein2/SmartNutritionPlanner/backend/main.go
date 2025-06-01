// main.go
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Anwendungsprojekt-Informatik/Anwendungsprojekt/api"
	"github.com/Anwendungsprojekt-Informatik/Anwendungsprojekt/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Drei Token-Buckets für die jeweiligen Limits:
var (
	// 100 Anfragen/Minute ≈ 1.666 req/s, Burst bis 100
	productLimiter = rate.NewLimiter(rate.Limit(100.0/60.0), 100)
	// 10 Anfragen/Minute ≈ 0.166 req/s, Burst bis 10
	searchLimiter = rate.NewLimiter(rate.Limit(10.0/60.0), 10)
	// 2 Anfragen/Minute ≈ 0.033 req/s, Burst bis 2
	facetLimiter = rate.NewLimiter(rate.Limit(2.0/60.0), 2)
)

// rateLimitMiddleware prüft, ob im jeweiligen Bucket noch ein Token verfügbar ist.
// Ist keins mehr vorhanden, gibt es HTTP 429 zurück.
func rateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			// Optional: Retry-After auf 60 Sekunden setzen
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"description": "Zu viele Anfragen in kurzer Zeit. Bitte später erneut versuchen.",
			})
			return
		}
		c.Next()
	}
}

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
		apiGroup.GET("/search",
			rateLimitMiddleware(searchLimiter), // 10 req/minute
			func(c *gin.Context) {
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
			},
		)

		// b) Barcode-Abfrage: GET /api/product/:barcode
		apiGroup.GET("/product/:barcode",
			rateLimitMiddleware(productLimiter), // 100 req/minute
			func(c *gin.Context) {
				barcode := c.Param("barcode")
				if barcode == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Barcode is required"})
					return
				}

				p, err := client.GetProductByBarcode(barcode)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, p)
			},
		)

		// c) Beispiel-Facet-Route: GET /api/facet/categories
		//    (2 req/minute)
		apiGroup.GET("/facet/categories",
			rateLimitMiddleware(facetLimiter), // 2 req/minute
			func(c *gin.Context) {
				// TODO: Datenquelle für Kategorien
				c.JSON(http.StatusOK, gin.H{
					"categories": []string{"Brot", "Obst", "Getränke"},
				})
			},
		)

		// Beispiel: weitere Facet-Routen ebenfalls mit facetLimiter
		apiGroup.GET("/facet/label/organic",
			rateLimitMiddleware(facetLimiter),
			func(c *gin.Context) {
				// TODO: Datenquelle für "organic"-Label
				c.JSON(http.StatusOK, gin.H{
					"label": "organic",
					"data":  []string{"Bio-Apfel", "Bio-Brot"},
				})
			},
		)

		// d) Produkte hinzufügen: POST /api/daily-entries
		apiGroup.POST("/daily-entries", middleware.HandleDailyEntries)

		// e) DailySummary Home: GET /api/daily-summary
		apiGroup.GET("/daily-summary", middleware.HandleDailySummary)
	}

	// ---------------------------------------------------
	// 4) Statisches Frontend unter "/public/*"
	//    -> Alle Dateien aus "../frontend/public" ausliefern
	// ---------------------------------------------------
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
	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.ListenAndServe()
}
