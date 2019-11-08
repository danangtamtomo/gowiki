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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type WikiOptions struct {
	Pageid string
	Title string
}

type WikiSuggestions struct {
	Status int
	Data []WikiOptions
}

func startWiki() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

		//reader := bufio.NewReader(os.Stdin)
		//println("Enter the text:\n")

		//keyword, _ := reader.ReadString('\n')
		//keyword = strings.TrimSuffix(keyword, "\n")
		keyword := r.URL.Query().Get("keyword")

		url := WIKI_ENDPOINT + "?action=query&list=search&srsearch=" + url.QueryEscape(keyword) + "&format=json"
		res, err := http.Get(url)

		check(err)

		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(res.Body)
			check(err)

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

			//wikiOptionsJSON, err := json.Marshal(options)
			//check(err)

			wikiSuggestions := WikiSuggestions{
				Status: 200,
				Data: options,
			}


			wikiSuggestionsJSON, err := json.Marshal(wikiSuggestions)

			w.Header().Set("Content-Type:", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write(wikiSuggestionsJSON)


			//choice := ""
			//
			//prompt := &survey.Select{
			//	Message: "Choose what do you want to lookup!",
			//	Options: options,
			//}
			//
			//survey.AskOne(prompt, &choice)
			//picked := strings.Split(choice, "|")
			//pageid := picked[0]

			//getPageUrl := WIKI_ENDPOINT + "?action=parse&prop=text&pageid=" + pageid + "&format=json"
			//
			//res, err := http.Get(getPageUrl)
			//
			//check(err)
			//
			//defer res.Body.Close()
			//
			//if res.StatusCode == http.StatusOK {
			//
			//	pageBytes, err := ioutil.ReadAll(res.Body)
			//	check(err)
			//
			//	pageString := string(pageBytes)
			//
			//	htmlText := gjson.Get(pageString, "parse.text")
			//
			//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText.String()))
			//
			//	check(err)
			//
			//	paragraph := doc.Find("p").Text()
			//
			//	unquotedParagraph, err := strconv.Unquote(`"` + paragraph + `"`)
			//	check(err)
			//
			//	result := html.UnescapeString(unquotedParagraph)
			//	w := 10
			//
			//	fmt.Printf(fmt.Sprintf("%%-%ds", w/2), fmt.Sprintf(fmt.Sprintf("%%%ds", w/2), result))
			//}

		}

	})


}

func reduce(arr []interface{}, cb func(acc interface{}, a interface{}, idx int, orArr []interface{}) interface{}, initialData *interface{}) {
	result := *initialData
	for i := range arr {
		result = cb(result, arr[i], i, arr)
	}

	*initialData = result

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

	startWiki()
	http.ListenAndServe(os.Getenv("PORT"), nil)
}
