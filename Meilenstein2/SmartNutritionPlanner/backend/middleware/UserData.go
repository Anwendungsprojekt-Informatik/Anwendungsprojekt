package middleware

import (
	"encoding/json"
	"fmt"
	"os"
)

// Struct passend zur JSON-Struktur
type UserData struct {
	Name          string `json:"name"`
	Age           string `json:"age"`
	Height        string `json:"height"`
	Weight        string `json:"weight"`
	CalorieTarget string `json:"calorieTarget"`
	TypeOfDiet    string `json:"typeOfDiet"`
}

// Funktion zum Erzeugen und Speichern der JSON-Daten
func SaveUserData(name, age, height, weight, calorieTarget, typeOfDiet string) error {
	data := UserData{
		Name:          name,
		Age:           age,
		Height:        height,
		Weight:        weight,
		CalorieTarget: calorieTarget,
		TypeOfDiet:    typeOfDiet,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der JSON: %v", err)
	}

	err = os.WriteFile("UserProfile.json", jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("Fehler beim Schreiben der Datei: %v", err)
	}

	return nil
}

// Funktion zum Einlesen und Zurückgeben der Daten
func LoadUserData() (UserData, error) {
	var data UserData

	fileBytes, err := os.ReadFile("UserProfile.json")
	if err != nil {
		return data, fmt.Errorf("Fehler beim Lesen der Datei: %v", err)
	}

	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return data, fmt.Errorf("Fehler beim Parsen der JSON: %v", err)
	}

	return data, nil
}
