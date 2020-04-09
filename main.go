package main

import (
	"os"
	"log"
	"time"
	"strconv"
	"net/http"
	"encoding/csv"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

/*
Welcome to my main file:)
This is the heart of the program.
You can find here all essential functions.

- main.go: core functions;
- check.go: parsing functions;
- models.go: struct that I'm using;
- requirements.txt: dependencies;
- install.sh: script to install and execute program;
*/

// Adds the date and the number of an article.
func addDate(Articles []Article, newArticle Article, info string) []Article {
	newArticle.Date = fullDate + " " + info
	if len(Articles) != 0 {
		newArticle.Count = Articles[len(Articles)-1].Count + 1
	} else {
		newArticle.Count = 1
	}
	Articles = append(Articles, newArticle)
	return Articles
}

// Adds title to the article.
func addHeader(Articles []Article, info string) []Article {
	Articles[len(Articles)-1].Title = info
	return Articles
}

// Adds the link of the article.
func addBody(Articles []Article, elem *colly.HTMLElement, URL string, c *colly.Collector) []Article {
	link := elem.Attr("href")
	link = URL + link + "/"
	Articles[len(Articles)-1].Link = link
    return Articles
}

// Creates new Articles and appends it to the slices.
func createArticle(Articles []Article, info string, flag int) []Article {
	var newArticle Article

	if flag == 1 {
		return addDate(Articles, newArticle, info)
	} else if flag == 2 {
		var err error
		Articles[len(Articles)-1].CommQuantity, err = strconv.Atoi(info)
		if err != nil {
			return nil
		}
		return Articles
	}

	return Articles
}

// Parses, creates and retrieves articles.
func getArticles(Articles []Article, data string) []Article {
	date, err := checkInput(data)
	if err != nil {
		log.Println("[Error]:", err.Error())
		os.Exit(1)
	}

	if date != 0 {
		Articles = createArticle(Articles, data, date)
	}
	return Articles
}

// Scraps the page and builds the slice of Articles.
func scrapURL(Articles []Article, URL string) []Article {
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
		Parallelism: 4,
	})

	// Setting random user-agent to avoid getting detected by the source.
	extensions.RandomUserAgent(c)

	// Trying to find div block with class="cat_news_item" and inside the span and links(href) in the <a> tag
	c.OnHTML(".cat_news_item", func(e *colly.HTMLElement) {
		e.ForEach("span", func(_ int, elem *colly.HTMLElement) {
			Articles = getArticles(Articles, elem.Text)
	    })
	    e.ForEach("a[href]", func(_ int, elem *colly.HTMLElement) {
	    	Articles = addHeader(Articles, elem.Text)
	    	Articles = addBody(Articles, elem, URL, c)
	    })
	})

	// Logging if we catch the request.
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
		log.Println("UserAgent", r.Headers.Get("User-Agent"))
	})

	c.Visit(URL)
	c.Wait()
	return Articles
}

// Scraps all the internal links.
func scrapInternalURL(Articles Article, URL string) Article {
	c := colly.NewCollector(
		colly.Async(true), // Making collector asynchronic
	)

	// Setting random user-agent to avoid getting detected by the source.
	extensions.RandomUserAgent(c)

	// Trying to find div block with class="fullnews white_block" and inside the paragraph
	c.OnHTML(".fullnews.white_block", func(e *colly.HTMLElement) {
	    e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
	    	Articles.Body = Articles.Body + elem.Text
	    })
	})

	// Trying to find div block with class="full_text" and inside the paragraph
	c.OnHTML(".full_text", func(e *colly.HTMLElement) {
	    e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
	    	Articles.Body = Articles.Body + elem.Text
	    })
	})

	// Trying to find div block with class="WordSection1" and inside the paragraph
	c.OnHTML(".WordSection1", func(e *colly.HTMLElement) {
	    e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
	    	Articles.Body = Articles.Body + elem.Text
	    })
	})

	// Logging if we catch the request.
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting->", r.URL)
	})

	c.Visit(URL) // Making request to the URL.
	c.Wait() // Waiting for all go-routines to be done.
	return Articles
}

// Create a CSV file, called zakon.csv and fills it with all articles.
func main() {
	var Articles []Article

	// Checking if URL is reachable.
	URL := "https://www.zakon.kz/news"
	_, err := http.Get(URL)
	if err != nil {
		log.Println(err.Error())
		return
	}

	Articles = scrapURL(Articles, URL)

	// Creating CSV file.
	fileName := "zakon.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Could not create %s", fileName)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing the output to the CSV file.
	writer.Write([]string{"Title", "Text", "Date", "CommentsQuantity"})
	for i, article := range Articles {
		Articles[i] = scrapInternalURL(article, article.Link)
	}

	for _, art := range Articles {
		writer.Write([]string{
			art.Title,
			art.Body,
			art.Date,
			strconv.Itoa(art.CommQuantity),
		})
		// log.Printf("title:[%s]\nlink:[%s]\nbody:[%s]\nComments:[%d]\n\n", art.Title, art.Link, art.Body, art.CommQuantity)
	}
}
