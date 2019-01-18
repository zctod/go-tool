// http连接
package curl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func CurlGet(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var mapResult map[string]interface{}

	if err := json.Unmarshal([]byte(string(body)), &mapResult); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return mapResult, nil
}

func CurlPost(dataMap map[string]interface{}) map[string]interface{} {
	dataJson, err := json.Marshal(dataMap)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	res, err := http.Post("http://127.0.0.1:8080/json", "application/json; encoding=utf-8", strings.NewReader(string(dataJson)))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var mapResult map[string]interface{}

	if err := json.Unmarshal([]byte(string(body)), &mapResult); err != nil {
		log.Fatal(err)
		return nil
	}
	return mapResult
}