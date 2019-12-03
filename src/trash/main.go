package main

import (
	"encoding/json"

	"github.com/mmcloughlin/geohash"
)

type inlocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func main() {
	var input inlocation

	anan := `{"lat" : 38.41822780 , "lng" : 27.14777520}`
	fuck := []byte(anan)

	err := json.Unmarshal(fuck, &input)

	if err != nil {
		println(err.Error())
	}

	hash := geohash.Encode(input.Lat, input.Lng)

	println(hash)
}
