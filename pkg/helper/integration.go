package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func MakeRequest(method, url string, reqBody []byte) (int, interface{}, error) {
	var response interface{}
	client := &http.Client{}

	reqReader := bytes.NewReader(reqBody)
	req, err := http.NewRequest(method, url, reqReader)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	if resp == nil {
		return http.StatusInternalServerError, response, fmt.Errorf("error while connecting to %s", url)
	}

	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return resp.StatusCode, string(body), err
	}

	return resp.StatusCode, response, nil
}
