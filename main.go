package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var filepath = "clocked.json"

type ClockedItem struct {
	Activity  string `json:"activity"`
	Frequency int    `json:"frequency"`
}

func makeClock(s string, z int) ClockedItem {

	NewActivity := ClockedItem{
		Activity:  s,
		Frequency: z,
	}

	return NewActivity

}

func dataEntry(filedata map[string]ClockedItem) error {

	data, err := json.Marshal(filedata)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil

}

func readFromFile(m map[string]ClockedItem) {

	file_contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	var clocked ClockedItem
	if err := json.Unmarshal(file_contents, &clocked); err != nil {

		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	fmt.Println(m)

	//kan bare lese et element enn så lenge

}

func getFromFile(filepath, input string, m map[string]ClockedItem) ClockedItem {

	file_contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	var clocked ClockedItem
	if err := json.Unmarshal(file_contents, &clocked); err != nil {

		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	//bør ha en sjekk om den fins?
	return m[input]

	//kan bare lese et element enn så lenge

}

func getMapfromFile(filepath string) map[string]ClockedItem {

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// File doesn't exist, return empty map
		return make(map[string]ClockedItem)
	}

	file_contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var struturedData map[string]ClockedItem
	if err := json.Unmarshal(file_contents, &struturedData); err != nil {

		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return struturedData

}

func clockIn(input string, data map[string]ClockedItem) error {

	editable := getFromFile(filepath, input, data)

	editable.Frequency--

	data[editable.Activity] = editable

	if editable.Frequency < 0 {
		return fmt.Errorf("too small")
	}

	dataEntry(data)

	fmt.Println("Clocked in !")

	readFromFile(data)

	return nil

}

func interactiveInit(key string, freq int, m map[string]ClockedItem) {

	m[key] = makeClock(key, freq)

	dataEntry(m)

}

func main() {

	lookupMap := getMapfromFile(filepath)

	activity := ""
	frequency := 0
	clocking := ""

	flag.StringVar(&activity, "a", "myActivity", "your activity here")
	flag.StringVar(&clocking, "c", "", "your clocked activity here")
	flag.IntVar(&frequency, "f", 0, "your frequency here")

	//bruker default verdien når man klokker inn, clocking må kanskje ta inn andre args

	flag.Parse()

	if activity != "" && frequency != 0 {

		interactiveInit(activity, frequency, lookupMap)

	}

	//clocked og freq samtidig kan bli sketch

	if clocking != "" {

		//kanskje ha noe listing av json filen på forhånd.

		clockIn(clocking, lookupMap)

	}

	//capitalize eller standardiser før en check.

}

//TODO: Fiks litt mer feilmeldinger, implementer lister og appending, senere tid og sånn
//TODO: Implementer Receiver funksjoner.
//TODO: User Input --kinda, sanitize med stor forbokstav
//TODO:Hvis fil ikke finnes, lag den.
