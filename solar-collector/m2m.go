package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	// fechas de inicio y fin
	end := time.Now()
	start := end.AddDate(0, -1, 0)

	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	url := fmt.Sprintf(
		"https://kauai.ccmc.gsfc.nasa.gov/DONKI/search/results?startDate=%s&endDate=%s&catalog=M2M_CATALOG&eventType=FLR",
		startDate, endDate,
	)

	fmt.Printf("Consultando: %s\n", url)

	// guardar eventos
	var eventos []map[string]string

	c := colly.NewCollector()

	c.OnHTML("table tbody", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, row *goquery.Selection) {
			evento := make(map[string]string)

			row.Find("td").Each(func(j int, col *goquery.Selection) {
				text := col.Text()
				header := fmt.Sprintf("col%d", j)
				evento[header] = text
			})

			if len(evento) > 0 {
				eventos = append(eventos, evento)
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error en la colección: %v", err)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error al visitar la URL: %v", err)
	}

	// carpeta (si no existe)
	dir := "resultadosM2M"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("Error al crear carpeta: %v", err)
		}
	}

	// archivo con timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("rslt_%s.json", timestamp)
	fullPath := filepath.Join(dir, fileName)

	// resultados
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("Error al crear archivo: %v", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(eventos)
	if err != nil {
		log.Fatalf("Error al guardar JSON: %v", err)
	}

	fmt.Printf("Se extrajeron %d eventos. Guardado en %s\n", len(eventos), fullPath)
}
