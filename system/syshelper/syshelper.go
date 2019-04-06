package syshelper

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"github.com/pkg/errors"
	"math/rand"
)

func GenerateNewSessionId() (string) {
	var b= string("AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789")
	var r= ""
	for i := 0; i < 20; i++ {
		r += string(b[rand.Int()%len(b)])
	}
	return r
}

type Serialized map[string]interface{}

func SerializeStruct(m *Serialized) (string, error) {
	var b = bytes.Buffer{}
	var encoder = gob.NewEncoder(&b)
	err := encoder.Encode(*m)
	if err != nil {
		return "", errors.Wrap(err, "failed gob encode")
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func DeserializeStruct(str string) (*Serialized, error) {
	m := &Serialized{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.Wrap(err, "failed base64 decode")
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed gob decode")
	}
	return m, nil
}

