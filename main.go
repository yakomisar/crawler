package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var client http.Client

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client = http.Client{
		Jar: jar,
	}
}

func main() {
	uRL := "http://185.204.3.165"
	req, err := http.NewRequest("GET", uRL, nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}
	cookie := &http.Cookie{
		//Name:   "token",
		//Value:  "my_token",
		//MaxAge: 300,
	}
	urlObj, _ := url.Parse(uRL)
	client.Jar.SetCookies(urlObj, []*http.Cookie{cookie})
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Get start link to initiate Test
	startLink, is := doc.Find("a").Attr("href")
	if is != true {
		log.Println("There is no link to start Test")
		return
	}

	req, err = http.NewRequest("GET", "http://185.204.3.165"+startLink, nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	////Преобразовываем массив байт в строку и выводим на печать
	//sb := string(body)
	//log.Printf(sb)

	// create from a file
	f, err := os.Open("main.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc1, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	} // use the goquery document... _ = doc.Find("h1")

	// Load the HTML document
	//doc1, err := goquery.NewDocumentFromReader(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}

	form := doc1.Find("form")
	results := make(map[string]string)
	form.Find("p").Each(func(i int, s *goquery.Selection) {
		name, value, is := checkSelect(s)
		if is == true {
			results[name] = value
		}
	})
	for key, val := range results {
		fmt.Println("Key " + key + " Value " + val)
	}
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
			//fmt.Printf("option: %s\n", val)
		})
		return name, value, true
	} else {
		fmt.Println("Другое значение")
	}
	return "", "", false
}
