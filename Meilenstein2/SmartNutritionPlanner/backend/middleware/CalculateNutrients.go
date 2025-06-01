package middleware

import (
	"net/http"

	"github.com/Anwendungsprojekt-Informatik/Anwendungsprojekt/api"
)

type NutrientsTotal struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Sugars   float64 `json:"sugar"`
}

func getNutritionValues(barcode string, amountInGrams float64) (NutrientsTotal, error) {
	client := api.OpenFoodFactClient{
		BaseURL:    "https://world.openfoodfacts.org/api/v0",
		HTTPClient: &http.Client{},
	}

	product, err := client.GetProductByBarcode(barcode)
	if err != nil {
		return NutrientsTotal{}, err
	}
	factor := amountInGrams / 100.0
	scaled := NutrientsTotal{
		Calories: product.EnergyKcal100g * factor,
		Protein:  product.Proteins100g * factor,
		Fat:      product.Fat100g * factor,
		Carbs:    product.Carbohydrates100g * factor,
		Sugars:   product.Sugars100g * factor,
	}

	return scaled, nil

}
