package main

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"

	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

const WIKI_ENDPOINT = "https://en.wikipedia.org/w/api.php"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getPeople() {
	reader := bufio.NewReader(os.Stdin)
	println("Enter the text:\n")

	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSuffix(keyword, "\n")

	url := WIKI_ENDPOINT + "?action=query&list=search&srsearch=" + url.QueryEscape(keyword) + "&format=json"
	res, err := http.Get(url)

	check(err)

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		check(err)

		bodyString := string(bodyBytes)

		value := gjson.Get(bodyString, "query.search")

		options := make([]string, 0)

		for _, val := range value.Array() {
			title := gjson.Get(val.String(), "title")
			pageid := gjson.Get(val.String(), "pageid")
			options = append(options, pageid.String()+"|"+title.String())
		}

		choice := ""

		prompt := &survey.Select{
			Message: "Choose what do you want to lookup!",
			Options: options,
		}

		survey.AskOne(prompt, &choice)
		picked := strings.Split(choice, "|")
		pageid := picked[0]

		getPageUrl := WIKI_ENDPOINT + "?action=parse&prop=text&pageid=" + pageid + "&format=json"

		res, err := http.Get(getPageUrl)

		check(err)

		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {

			pageBytes, err := ioutil.ReadAll(res.Body)
			check(err)

			pageString := string(pageBytes)

			htmlText := gjson.Get(pageString, "parse.text")

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText.String()))

			check(err)

			paragraph := doc.Find("p").Text()

			unquotedParagraph, err := strconv.Unquote(`"` + paragraph + `"`)
			check(err)

			result := html.UnescapeString(unquotedParagraph)
			w := 10

			fmt.Printf(fmt.Sprintf("%%-%ds", w/2), fmt.Sprintf(fmt.Sprintf("%%%ds", w/2), result))
		}

	}

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

	getPeople()
}
