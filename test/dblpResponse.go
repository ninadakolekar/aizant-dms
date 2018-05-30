package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
)

func DBLPResponse() {

	dbconf := []string{"VLDB", "SIGMOD", "PODS", "ICDE", "ICDT", "EDBT"}

	for _, conf := range dbconf {

		resp, err := getJSONResponse(conf)

		if err != nil {
			log.Fatal(err, " Line 21")
		}

		// fmt.Println(string(resp))

		jsonparser.ArrayEach(resp, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

			fmt.Println(jsonparser.Get(value, "url"))

		}, "result", "hits", "hit")

		if err != nil {
			log.Fatal(err)
		}
	}

}

func getJSONResponse(conf string) ([]byte, error) {

	url := getVenueURL(conf)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func getVenueURL(conf string) string {
	return "http://dblp.org/search/publ/api?q=venue%3A" + conf + "%3A&format=json"
}
