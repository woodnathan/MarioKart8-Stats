package main

import (
	"os"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"fmt"
)

var fields = []string { "type", "name", "speed", "speed_water", "speed_air", "speed_ground", "acceleration", "weight", "handling", "handling_water", "handling_air", "handling_ground", "traction", "mini_turbo"}


type MultipartElement struct {
	Value float64 `json:"value"`
	Water float64 `json:"water"`
	Air float64 `json:"air"`
	Ground float64 `json:"ground"`
}

type Speed struct {
	MultipartElement
}

type Handling struct {
	MultipartElement
}

type Record struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Speed Speed `json:"speed"`
	Acceleration float64 `json:"acceleration"`
	Weight float64 `json:"weight"`
	Handling Handling `json:"handling"`
	Traction float64 `json:"traction"`
	MiniTurbo float64 `json:"mini_turbo"`
}

func fieldIndex(field string) (int) {
	for i, v := range(fields) {
		if (v == field) {
			return i
		}
	}
	return 0
}

func parseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (element *MultipartElement) parseMultipartElement(fields []string) {
	element.Value = parseFloat64(fields[0])
	element.Water = parseFloat64(fields[1])
	element.Air = parseFloat64(fields[2])
	element.Ground = parseFloat64(fields[3])
}

func parseSpeed(fields []string) Speed {
	var s Speed
	s.parseMultipartElement(fields)
	return s
}
func parseHandling(fields []string) Handling {
	var h Handling
	h.parseMultipartElement(fields)
	return h
}

func generateFile(input string, output string) error {
	infile, err := os.Open(input)
	if (err != nil) {
		return err
	}
	
	// Read in the CSV file
	reader := csv.NewReader(infile)
	records, err := reader.ReadAll()
	infile.Close()
	
	// Iterate through the records
	records = records[1:]
	outrecords := make([]Record, len(records))
	for idx, record := range(records) {
		var outrecord Record
		
		outrecord.Type = record[fieldIndex("type")]
		outrecord.Name = record[fieldIndex("name")]
		outrecord.Speed = parseSpeed(record[fieldIndex("speed"):])
		outrecord.Acceleration = parseFloat64(record[fieldIndex("acceleration")])
		outrecord.Weight = parseFloat64(record[fieldIndex("weight")])
		outrecord.Handling = parseHandling(record[fieldIndex("handling"):])
		outrecord.Traction = parseFloat64(record[fieldIndex("traction")])
		outrecord.MiniTurbo = parseFloat64(record[fieldIndex("mini_turbo")])
		
		outrecords[idx] = outrecord
	}
	
	outfile, err := os.Create(output)
	if (err != nil) {
		return err
	}
	b, err := json.MarshalIndent(outrecords, "", "  ")
	outfile.Write(b)
	outfile.Close()
	
	return err
}

func readFile(name string) {
	input := fmt.Sprintf("%s.csv", name)
	output := fmt.Sprintf("json/%s.json", name)
	err := generateFile(input, output)
	if (err != nil) {
		fmt.Println(err)
	}
}

func main() {
	readFile("characters")
	readFile("bodies")
	readFile("tires")
	readFile("gliders")
}