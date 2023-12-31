package goe

import (
	"encoding/json"
	"log"
)

type Model interface {
	String() string
}

type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		log.Print(err)
	}
	return Models(b)
}
