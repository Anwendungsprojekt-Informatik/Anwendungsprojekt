package middleware

type NutrientsTotal struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
}

func getNutritionValues(productCode string, amountInGrams float64) (NutrientsTotal, error) {
	var nutrientsTotal NutrientsTotal //Hier müsste ein Funktionsaufruf für die entsprechende API Rückgabe hin
	factor := amountInGrams / 100.0
	scaled := NutrientsTotal{
		Calories: nutrientsTotal.Calories * factor,
		Protein:  nutrientsTotal.Protein * factor,
		Fat:      nutrientsTotal.Fat * factor,
		Carbs:    nutrientsTotal.Carbs * factor,
	}

	return scaled, nil

}
