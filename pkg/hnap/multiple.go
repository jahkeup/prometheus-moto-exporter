package hnap

import (
	"encoding/json"
	"errors"
)

type getMultipleRequest struct {
	HNAP map[string]string `json:"GetMultipleHNAPs"`
}

func GetMultipleRequestData(names ...string) interface{} {
	req := getMultipleRequest{map[string]string{}}

	for _, name := range names {
		req.HNAP[name] = ""
	}

	return req
}

type GetMultipleHNAPsResponse struct {
	HNAP map[string]json.RawMessage `json:"GetMultipleHNAPsResponse"`
}

func (g *GetMultipleHNAPsResponse) GetJSON(name string) (json.RawMessage, error) {
	data, ok := g.HNAP[name]
	if ok {
		return data, nil
	}
	// Might be namespaced under the request, check there too.
	data, ok = g.HNAP[name+"Response"]
	if !ok {
		return nil, errors.New("no response found ")
	}
	return data, nil
}
