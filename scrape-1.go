package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fetchURL := "https://www.imdb.com/list/ls033609554/"
	fileName := "disney-movies.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("ERROR: Could not create file %q: %s\n", fileName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write column headers of the text file
	writer.Write([]string{"Sl. No.", "Movie Name", "Release Year", "Certificate", "Genre",
		"Running time", "Rating", "Number of Votes", "Gross"})

	// Instantiate the default Collector
	c := colly.NewCollector()

	// Before making a request, print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	// Callback when colly finds the entry point to the DOM segment having a movie info
	c.OnHTML(`.lister-item-content`, func(e *colly.HTMLElement) {
		//Locate and extract different pieces information about each movie
		number := e.ChildText(".lister-item-index")
		name := e.ChildText(".lister-item-index ~ a")
		year := e.ChildText(".lister-item-year")
		runtime := e.ChildText(".runtime")
		certificate := e.ChildText(".certificate")
		genre := e.ChildText(".genre")
		rating := e.ChildText("[class='ipl-rating-star small'] .ipl-rating-star__rating")
		vote := e.ChildAttr("span[name=nv]", "data-value")
		gross := e.ChildText(".text-muted:contains('Gross') ~ span[name=nv]")

		// Write all scraped pieces of information to output text file
		writer.Write([]string{
			number,
			name,
			year,
			certificate,
			runtime,
			genre,
			rating,
			vote,
			gross,
		})
	})

	// start scraping the page under the given URL
	c.Visit(fetchURL)
	fmt.Println("End of scraping: ", fetchURL)
}
