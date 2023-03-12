package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GET[T any](host, endpoint string) ([]*T, error) {
	resp, err := http.Get(host + endpoint)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []*T

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
