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
	// Создание контейнера для куки
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Unable to create cookie container: %s", err)
		return nil, err
	}
	cookieClient := &http.Client{
		Jar: jar,
	}
	// Создание логгера
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		client: cookieClient,
		logger: logger,
	}
	return app, nil
}

func getPage(app *application, uRL string, method string) (*goquery.Document, error) {
	// Обработка стартовой страницы
	req, err := http.NewRequest(method, uRL, nil)
	if err != nil {
		app.logger.Fatalf("Unable to request using %s link: %s", uRL, err.Error())
		return nil, err
	}

	resp, err := app.client.Do(req)
	if err != nil {
		app.logger.Fatalf("Unable to request using %s link: %s", uRL, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		app.logger.Fatalf("Unable to read document: %s ", err.Error())
		return nil, err
	}
	return doc, err
}

func processWebPage(app *application, uRL string) error {
	_, err := getPage(app, uRL, "GET")
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
		questionLink := "/question/" + number
		// Получаем страницу с очередным вопросом
		doc1, err := getPage(app, uRL+questionLink, "GET")
		if err != nil {
			app.logger.Fatalf("Unable to read current document: %s ", err.Error())
			return err
		}
		// Парсим страницу с вопросами
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
		urlA, err := url.Parse("http://185.204.3.165" + questionLink)
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

		// Отправляем наши ответы с помощью метода POST
		doc, err := getPage(app, urlA.String(), "POST")
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
