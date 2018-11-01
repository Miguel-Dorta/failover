package main

import (
	"encoding/json"
	"io/ioutil"
)

type setting struct {
	Email string `json:"email,omitempty"`
	Key   string `json:"key,omitempty"`
	Zone  []zone `json:"zone,omitempty"`
}

type zone struct {
	ID     string   `json:"id,omitempty"`
	Record []record `json:"record,omitempty"`
}

type record struct {
	ID         string `json:"id,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	Name       string `json:"name,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	Proxied    bool   `json:"proxied"`
}

func readSettings(path string) (setting, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return setting{}, err
	}

	var sett setting
	json.Unmarshal(data, &sett)
	return sett, nil
}
