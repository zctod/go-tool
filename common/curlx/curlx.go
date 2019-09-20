package curlx

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
)

type FileInfo struct {
	Name   string
	Stream io.Reader
}

type FormData struct {
	File   map[string]FileInfo
	Params map[string]string
}

type HttpReq struct {
	Url        string
	Method     string
	Header     map[string]string
	Query      map[string]string
	Params     map[string]string
	FormData   FormData
	Body       []byte
	BodyReader io.Reader
	CertFile   string
	KeyFile    string
	Timeout    time.Duration
}

func (h *HttpReq) buildUrl() {
	if h.Query == nil || len(h.Query) == 0 {
		return
	}

	query := url.Values{}
	for k, v := range h.Query {
		query.Set(k, v)
	}

	urlSet := strings.Split(h.Url, "?")
	switch len(urlSet) {
	case 1:
		h.Url += "?" + query.Encode()
	case 2:
		if urlSet[1] != "" {
			urlSet[1] += "&"
		}
		h.Url = urlSet[0] + "?" + url.PathEscape(urlSet[1]+query.Encode())
	}
}

func (h *HttpReq) buildBody() {
	if h.Body != nil || h.BodyReader != nil {
		return
	}

	params := url.Values{}
	for k, v := range h.Params {
		params.Set(k, v)
	}
	h.Body = []byte(params.Encode())
}

func (h *HttpReq) Do() ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if h.CertFile != "" {
		cert, err := tls.LoadX509KeyPair(h.CertFile, h.KeyFile)
		if err != nil {
			return nil, err
		}
		tr.DisableCompression = true
		tr.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
	var client = &http.Client{
		Transport: tr,
		Timeout:   h.Timeout,
	}

	h.buildUrl()
	h.buildBody()

	var bReader io.Reader
	if h.BodyReader != nil {
		bReader = h.BodyReader
	} else {
		bReader = bytes.NewReader(h.Body)
	}

	req, err := http.NewRequest(h.Method, h.Url, bReader)
	if err != nil {
		return nil, err
	}

	if h.Header != nil {
		for k, v := range h.Header {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (h *HttpReq) Get() ([]byte, error) {

	h.Method = METHOD_GET
	return h.Do()
}

func (h *HttpReq) Post() ([]byte, error) {

	h.Method = METHOD_POST
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	if _, ok := h.Header["Content-Type"]; !ok {
		h.Header["Content-Type"] = "application/x-www-form-urlencoded"
	}

	return h.Do()
}

func (h *HttpReq) PostForm() ([]byte, error) {

	h.Method = METHOD_POST

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if h.FormData.File != nil {
		for k, file := range h.FormData.File {
			part, err := w.CreateFormFile(k, file.Name)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(part, file.Stream); err != nil {
				return nil, err
			}
		}
	}

	if h.FormData.Params != nil {
		for k, v := range h.FormData.Params {
			if err := w.WriteField(k, v); err != nil {
				return nil, err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	h.Header["Content-Type"] = w.FormDataContentType()
	h.BodyReader = &buf

	return h.Do()
}
