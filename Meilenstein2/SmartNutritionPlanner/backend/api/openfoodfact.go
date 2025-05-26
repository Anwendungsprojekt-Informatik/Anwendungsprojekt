// openfoodfact.go
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// OpenFoodFactClient repräsentiert einen Client für die OpenFoodFacts API
type OpenFoodFactClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Product repräsentiert ein Lebensmittelprodukt aus der OpenFoodFacts API
type Product struct {
	ID              string                     `json:"id"`
	Code            string                     `json:"code"`
	ProductName     string                     `json:"product_name"`
	Brands          string                     `json:"brands"`
	ImageURL        string                     `json:"image_url"`
	Ingredients     string                     `json:"ingredients_text"`
	NutritionGrades string                     `json:"nutrition_grades"`
	Quantity        string                     `json:"quantity"`
	Nutriments      map[string]json.RawMessage `json:"nutriments"`

	// Neue Felder für Nährwerte pro 100 g
	EnergyKcal100g    float64 `json:"energy-kcal_100g"`
	Carbohydrates100g float64 `json:"carbohydrates_100g"`
	Sugars100g        float64 `json:"sugars_100g"`
	Proteins100g      float64 `json:"proteins_100g"`
	Fat100g           float64 `json:"fat_100g"`
}

// SearchResponse repräsentiert die Antwort der OpenFoodFacts API bei einer Suche
type SearchResponse struct {
	Count    int       `json:"count"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Products []Product `json:"products"`
}

// NewOpenFoodFactClient erstellt einen neuen OpenFoodFacts API Client
func NewOpenFoodFactClient() *OpenFoodFactClient {
	return &OpenFoodFactClient{
		BaseURL:    "https://world.openfoodfacts.org/api/v2",
		HTTPClient: &http.Client{},
	}
}

// Hilfsfunktion: Wert als float64 extrahieren (ebenfalls Strings parsen)
func getFloatValue(data map[string]json.RawMessage, key string) float64 {
	if value, ok := data[key]; ok {
		// Versuch Float direkt
		var f float64
		if err := json.Unmarshal(value, &f); err == nil {
			return f
		}
		// Versuch String → Float
		var s string
		if err := json.Unmarshal(value, &s); err == nil {
			if p, err := strconv.ParseFloat(s, 64); err == nil {
				return p
			}
		}
	}
	return 0
}

// GetProductByBarcode sucht ein Produkt anhand seines Barcodes und liest Nährwerte
func (c *OpenFoodFactClient) GetProductByBarcode(barcode string) (*Product, error) {
	url := fmt.Sprintf("%s/product/%s.json", c.BaseURL, barcode)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Product Product `json:"product"`
		Status  int     `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Status != 1 {
		return nil, fmt.Errorf("product not found or API error")
	}

	p := &result.Product
	p.EnergyKcal100g = getFloatValue(p.Nutriments, "energy-kcal_100g")
	p.Carbohydrates100g = getFloatValue(p.Nutriments, "carbohydrates_100g")
	p.Sugars100g = getFloatValue(p.Nutriments, "sugars_100g")
	p.Proteins100g = getFloatValue(p.Nutriments, "proteins_100g")
	p.Fat100g = getFloatValue(p.Nutriments, "fat_100g")

	return p, nil
}

// SearchProducts sucht Produkte anhand eines Suchbegriffs
func (c *OpenFoodFactClient) SearchProducts(query string, page int) (*SearchResponse, error) {
	if page < 1 {
		page = 1
	}
	baseURL := "https://world.openfoodfacts.org/cgi/search.pl"
	params := url.Values{}
	params.Add("search_terms", query)
	params.Add("json", "1")
	params.Add("page", strconv.Itoa(page))
	params.Add("page_size", "10")
	params.Add("search_simple", "1")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
