package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// OpenFoodFactClient repräsentiert einen Client für die OpenFoodFacts API
type OpenFoodFactClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Product repräsentiert ein Lebensmittelprodukt aus der OpenFoodFacts API
type Product struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	ProductName     string `json:"product_name"`
	Brands          string `json:"brands"`
	ImageURL        string `json:"image_url"`
	Ingredients     string `json:"ingredients_text"`
	NutritionGrades string `json:"nutrition_grades"`
	Quantity        string `json:"quantity"`
}

// Product Nährwerte
type Nutriment struct {
	Calories string `json:"nutriments.energy-kcal_100g"`
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
		BaseURL:    "https://world.openfoodfacts.org/api/v0",
		HTTPClient: &http.Client{},
	}
}

// GetProductByBarcode sucht ein Produkt anhand seines Barcodes
func (c *OpenFoodFactClient) GetProductByBarcode(barcode string) (*Product, error) {
	url := fmt.Sprintf("%s/product/%s.json", c.BaseURL, barcode)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Product   Product   `json:"product"`
		Nutriment Nutriment `json:"nutriment"`
		Status    int       `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != 1 {
		return nil, fmt.Errorf("product not found or API error")
	}

	return &result.Product, nil
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
	params.Add("page", fmt.Sprintf("%d", page))
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
