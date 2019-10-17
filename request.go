package yz_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

const (
	BaseURL   = "https://open.youzanyun.com"
	TokenPath = "/auth/token"
	OrderPath = "/api/youzan.trades.sold.get"
)

const (
	typeJSON      = "application/json"
	typeForm      = "application/x-www-form-urlencoded"
	typeMultipart = "multipart/form-data"
)

type YZBaseResponse struct {
	Success bool
	Code    int
}

type uploadFile struct {
	FileName  string
	FieldName string
	Reader    io.Reader
}

func (c *Client) httpJSON(path interface{}, params url.Values, req interface{}, res interface{}) error {
	data, err := c.httpRequest(path, params, req)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, res)
	return err
}

func (c *Client) httpRequest(path interface{}, params url.Values, requestData interface{}) ([]byte, error) {
	var request *http.Request
	var requestUrl string
	client := c.HTTPClient

	requestUrl = BaseURL + path.(string) + "?" + params.Encode()
	fmt.Printf("requestUrl=%s\n", requestUrl)
	if requestData != nil {
		switch v := requestData.(type) {
		case *uploadFile:
			var b bytes.Buffer
			if v.Reader == nil {
				return nil, errors.New("upload file is empty")
			}
			w := multipart.NewWriter(&b)
			fw, err := w.CreateFormFile(v.FieldName, v.FileName)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(fw, v.Reader); err != nil {
				return nil, err
			}
			if err = w.Close(); err != nil {
				return nil, err
			}
			request, _ = http.NewRequest("POST", requestUrl, &b)
			request.Header.Set("Content-Type", w.FormDataContentType())
		default:
			d, _ := json.Marshal(requestData)
			request, _ = http.NewRequest("POST", requestUrl, bytes.NewReader(d))
			request.Header.Set("Content-Type", typeJSON+"; charset=UTF-8")
		}
	} else {
		request, _ = http.NewRequest("GET", requestUrl, nil)
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Server Error: " + resp.Status)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
