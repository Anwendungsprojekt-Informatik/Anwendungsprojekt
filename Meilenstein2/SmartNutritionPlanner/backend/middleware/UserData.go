package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
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

type Nutrients struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Sugar    float64 `json:"sugar"`
}

type ProductEntry struct {
	Name          string    `json:"name"`
	Barcode       string    `json:"barcode"`
	AmountInGrams float64   `json:"amountInGrams"`
	Nutrients     Nutrients `json:"nutrients"`
}

type DayEntry struct {
	Date    string         `json:"date"`    // z. B. "2025-05-30"
	Entries []ProductEntry `json:"entries"` // Liste von Produkten
}

type DailyFood struct {
	Days []DayEntry `json:"dailyFoodLog.json"`
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

// Funktion zum Einfügen/Ändern der Produkte für den Tag
func AddDailyEntries(newEntries []ProductEntry) error {
	var data DailyFood
	var filename = "DailyFoodLog.json"

	// Datei einlesen (falls vorhanden)
	content, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("file read error: %w", err)
	}

	// Heutiges Datum im Format YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	// Prüfen, ob der Tag bereits existiert
	found := false
	for i, day := range data.Days {
		if day.Date == today {
			// Wenn ja: existierende Einträge ersetzen (oder erweitern)
			data.Days[i].Entries = newEntries
			found = true
			break
		}
	}

	// Falls nicht gefunden: neuen Tag hinzufügen
	if !found {
		data.Days = append(data.Days, DayEntry{
			Date:    today,
			Entries: newEntries,
		})
	}

	// Speichern
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// Funktion für Gesamtnährwerte eines bestimmten Tags
func CalculateTotalNutrientsForToday() (Nutrients, error) {
	var data DailyFood
	total := Nutrients{}
	var filename = "DailyFoodLog.json"

	// Datei einlesen
	content, err := os.ReadFile(filename)
	if err != nil {
		return total, fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return total, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Heutiges Datum im Format YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	// Tageseintrag finden
	for _, day := range data.Days {
		if day.Date == today {
			for _, entry := range day.Entries {
				total.Calories += entry.Nutrients.Calories
				total.Protein += entry.Nutrients.Protein
				total.Fat += entry.Nutrients.Fat
				total.Carbs += entry.Nutrients.Carbs
				total.Sugar += entry.Nutrients.Sugar
			}
			return total, nil
		}
	}

	return total, fmt.Errorf("no entries found for today (%s)", today)
}
