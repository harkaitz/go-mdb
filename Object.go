package mdb

import (
	"encoding/json"
)

type Object map[string]string 

func (s Object) Reset() {
	s = Object{}
}

func (s Object) Validate() error {
	return nil
}

func (s Object) Set(key, val string) error {
	s[key] = val
	return nil
}

func (s Object) Marshal() (b []byte, err error) {
	return json.Marshal(s)
}

func (s Object) Unmarshal(b []byte) (err error) {
	return json.Unmarshal(b, &s)
}
