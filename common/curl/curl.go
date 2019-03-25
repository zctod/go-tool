// http连接
package curl

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func httpCurl(method, urlVal, data string, headerMap map[string]string, cookieMap map[string]string) ([]byte, error) {

	client := &http.Client{}
	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal, "?")
		if len(urlArr) == 2 {
			//将GET请求的参数进行转义
			urlVal = urlArr[0] + "?" + url.PathEscape(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	} else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	//添加cookie，key为X-Xsrftoken，value为df41ba54db5011e89861002324e63af81
	//可以添加多个cookie
	//cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	//req.AddCookie(cookie1)
	if len(cookieMap) > 0 {
		for ck, cv := range cookieMap {
			req.AddCookie(&http.Cookie{
				Name: ck,
				Value: cv,
				HttpOnly: true,
			})
		}
	}

	//添加header，key为X-Xsrftoken，value为b6d695bbdcd111e8b681002324e63af81
	//req.Header.Add("X-Xsrftoken", "b6d695bbdcd111e8b681002324e63af81")
	if len(headerMap) > 0 {
		for hk, hv := range headerMap {
			req.Header.Add(hk, hv)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	return b, nil
}

func ToMap(b []byte) (map[string]interface{}, error) {

	var mapResult map[string]interface{}
	if err := json.Unmarshal([]byte(string(b)), &mapResult); err != nil {
		return nil, err
	}
	return mapResult, nil
}

func CurlGet(url string, headerMap map[string]string) (map[string]interface{}, error) {

	b, err := httpCurl("GET", url, "", headerMap, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = json.Unmarshal(b, &res)
	return res, err
}

func CurlPost(url string, dataMap map[string]interface{}, headerMap map[string]string) (map[string]interface{}, error) {

	var dataStr = ""
	for k, v := range dataMap {
		if dataStr != "" {
			dataStr += "&"
		}
		dataStr += k + "=" + v.(string)
	}

	b, err := httpCurl("POST", url, dataStr, headerMap, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = json.Unmarshal(b, &res)
	return res, err
}
