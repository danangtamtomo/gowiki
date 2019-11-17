package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const WIKI_ENDPOINT = "https://en.wikipedia.org/w/api.php"

type WikiOptions struct {
	Pageid string
	Title  string
}

type WikiSuggestions struct {
	Status int
	Data   []WikiOptions
}

type Page struct {
	Title string
	PageID string
	Text string
}

type WikiPage struct {
	Status int
	Content Page
}

func StartWiki(w http.ResponseWriter, r *http.Request) {

	keyword := r.URL.Query().Get("keyword")

	URL := WIKI_ENDPOINT + "?action=query&list=search&srsearch=" + url.QueryEscape(keyword) + "&format=json"
	res, err := http.Get(URL)

	Check(err)

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		Check(err)

		bodyString := string(bodyBytes)

		value := gjson.Get(bodyString, "query.search")

		options := make([]WikiOptions, 0)

		for _, val := range value.Array() {
			title := gjson.Get(val.String(), "title")
			pageid := gjson.Get(val.String(), "pageid")
			options = append(options, WikiOptions{
				Pageid: pageid.String(),
				Title:  title.String(),
			})
		}

		wikiSuggestions := WikiSuggestions{
			Status: 200,
			Data:   options,
		}

		wikiSuggestionsJSON, err := json.Marshal(wikiSuggestions)

		Check(err)

		w.Header().Set("Content-Type:", "application/json")
		w.WriteHeader(http.StatusOK)

		res, _ := w.Write(wikiSuggestionsJSON)

		print(res)
	}

}

func ReadWikiPage(w http.ResponseWriter, r *http.Request) {

	pageID := GetParamByKey(r, "pageID")

	URL := WIKI_ENDPOINT + "?action=parse&prop=text&format=json&pageid=" + pageID[2]
	res, err := http.Get(URL)
	defer res.Body.Close()

	Check(err)
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		Check(err)

		bodyString := string(bodyBytes)
		page := Page{
			Title:  gjson.Get(bodyString, "parse.title").String(),
			PageID: gjson.Get(bodyString, "parse.pageid").String(),
			Text:   gjson.Get(bodyString, "parse.text.*").String(),
		}

		wikiPage := WikiPage{
			Status:  200,
			Content: page,
		}

		pageResult, err := json.Marshal(wikiPage)

		Check(err)

		w.Header().Set("Content-Type:", "application/json")
		w.WriteHeader(http.StatusOK)

		res, _ := w.Write(pageResult)

		print(res)
	}


}

func reduce(arr []interface{}, cb func(acc interface{}, a interface{}, idx int, orArr []interface{}) interface{}, initialData *interface{}) {
	result := *initialData
	for i := range arr {
		result = cb(result, arr[i], i, arr)
	}

	*initialData = result
}

func startRoutes() {
	http.HandleFunc("/search", StartWiki)
	http.HandleFunc("/page/", ReadWikiPage)
}

func main() {
	names := []interface{}{"danang", "aji", "tamtomo",}
	var result interface{} = ""
	reduce(
		names,
		func(a interface{}, b interface{}, i int, arr []interface{}) interface{} {
			res := a.(string) + " dan " + b.(string)
			return res
		},
		&result,
	)
	fmt.Println(result.(string))

	startRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
	}

	http.ListenAndServe(":" + port, nil)
}
