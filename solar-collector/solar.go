package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type SolarFlare struct {
	Date   string `json:"date"`
	Time   string `json:"time"`
	Class  string `json:"class"`
	Region string `json:"region"`
}

func main() {
	// contexto de Chrome headless
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var rowsText []string
	url := "https://www.spaceweatherlive.com/en/solar-activity/solar-flares.html"

	fmt.Println("Requesting:", url)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#datatable-solarflares tbody`, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll("#datatable-solarflares tbody tr"))
				.map(row => Array.from(row.children).map(cell => cell.textContent.trim()).join(";"))
		`, &rowsText),
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(rowsText) == 0 {
		fmt.Println("No solar flares found.")
		fmt.Println("Error collecting flares.")
		return
	}

	var flares []SolarFlare
	for _, row := range rowsText {
		parts := strings.Split(row, ";")
		if len(parts) < 4 {
			continue
		}
		flares = append(flares, SolarFlare{
			Date:   parts[0],
			Time:   parts[1],
			Class:  parts[2],
			Region: parts[3],
		})
	}

	// carpeta resultados
	os.MkdirAll("resultados", os.ModePerm)

	// archivo con timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join("resultados", fmt.Sprintf("solar_flares_%s.json", timestamp))

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error al crear archivo:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(flares)

	fmt.Println("Datos guardados en:", filename)
}
