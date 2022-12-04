package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
)

func checkTextRadio(s *goquery.Selection) (string, string, bool) {
	attrType, exist := s.Find("input").Attr("type")
	if exist && attrType == "radio" {
		value := ""
		name, _ := s.Find("input").Attr("name")
		s.Find("input").Each(func(a int, x *goquery.Selection) {
			val, _ := x.Attr("value")
			if len(val) > len(value) {
				value = val
			}
		})

		return name, value, true
	}
	return "", "", false
}

func checkText(s *goquery.Selection) (string, string, bool) {
	sel := s.Find("input")
	inputType, exist := sel.Attr("type")
	if exist {
		if inputType == "radio" {

		} else if inputType == "text" {
			name, _ := sel.Attr("name")
			return name, "test", true
		}
	}
	return "", "", false
}

func checkSelect(s *goquery.Selection) (string, string, bool) {
	sel := s.Find("select")
	name, exist := sel.Attr("name")
	if exist {
		value := ""
		sel.Find("option").Each(func(a int, x *goquery.Selection) {
			val, _ := x.Attr("value")
			if len(val) > len(value) {
				value = val
			}
		})
		return name, value, true
	}
	return "", "", false
}

type application struct {
	client *http.Client
	logger *log.Logger
}

func app_init(url string) (*application, error) {
	// Cookie container creation
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Unable to create cookie container: %s", err)
		return nil, err
	}
	cookieClient := &http.Client{
		Jar: jar,
	}
	// Create logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		client: cookieClient,
		logger: logger,
	}
	return app, nil
}

func submitStartPage(app *application, uRL string) (*goquery.Document, error) {
	// Обработка стартовой страницы
	req, err := http.NewRequest("GET", uRL, nil)
	if err != nil {
		log.Fatalf("Got an error %s", err.Error())
		return nil, err
	}

	resp, err := app.client.Do(req)
	if err != nil {
		app.logger.Fatalf("Error in the Do(req): %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return doc, err
}

func processWebPage(app *application, uRL string) error {
	_, err := submitStartPage(app, uRL)
	if err != nil {
		app.logger.Fatalf("Error in submitting start page: %s", err.Error())
		return err
	}
	header := ""
	for i := 1; ; i++ {
		if header == "Test successfully passed" {
			break
		}
		number := strconv.Itoa(i)
		startLink := "/question/" + number
		req, err := http.NewRequest("GET", "http://185.204.3.165"+startLink, nil)
		if err != nil {
			app.logger.Fatalf("Unable to request using %s link: %s", startLink, err.Error())
			return err
		}
		resp, err := app.client.Do(req)
		if err != nil {
			app.logger.Fatalf("Unable to request using %s link: %s", startLink, err.Error())
			return err
		}
		defer resp.Body.Close()
		doc1, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			app.logger.Fatalf("Unable to read document: %s ", err.Error())
			return err
		}
		form := doc1.Find("form")
		results := make(map[string]string)
		form.Find("p").Each(func(i int, s *goquery.Selection) {
			name, value, is := checkSelect(s)
			if is == true {
				results[name] = value
			}
			name, value, is = checkTextRadio(s)
			if is == true {
				results[name] = value
			}
			name, value, is = checkText(s)
			if is == true {
				results[name] = value
			}
		})
		urlA, err := url.Parse("http://185.204.3.165" + startLink)
		if err != nil {
			app.logger.Fatalf("Unable to parse: %s ", err.Error())
			return err
		}
		app.logger.Println("Processing: ", urlA.String())
		values := urlA.Query()

		for key, val := range results {
			values.Add(key, val)
		}
		urlA.RawQuery = values.Encode()
		req, err = http.NewRequest("POST", urlA.String(), nil)
		if err != nil {
			app.logger.Fatalf("Unable to POST: %s", err.Error())
			return err
		}

		resp, err = app.client.Do(req)
		if err != nil {
			app.logger.Fatalf("Error in the Do(req): %s", err.Error())
			return err
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			app.logger.Fatalf("Unable to get document header: %s", err.Error())
			return err
		}
		// Обновляем header для проверки
		header = doc.Find("h1").Text()
	}
	return nil
}

func main() {
	// Инициализация приложения
	uRL := "http://185.204.3.165"
	app, err := app_init(uRL)
	if err != nil {
		app.logger.Fatalf("Problem with application launch: %s", err)
		return
	}
	// Старт приложения
	app.logger.Println("Application started...")
	err = processWebPage(app, uRL)
	if err != nil {
		app.logger.Fatal(err)
		return
	}
	app.logger.Println("Test successfully passed.")
}
