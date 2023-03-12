package jmespath

import (
	"encoding/json"
	"fmt"
	"github.com/jmespath/go-jmespath"
)

func TestCondition(condition string, event any) (bool, error) {
	marshal, err := json.Marshal(event)
	if err != nil {
		return false, err
	}
	var data interface{}
	err = json.Unmarshal(marshal, &data)
	if err != nil {
		return false, err
	}
	res, err := jmespath.Search(condition, data)
	if err != nil {
		return false, err
	}

	v, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("result of filter was not of type boolean")
	}

	return v, nil
}
