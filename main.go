package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

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

func dataEntry(c ClockedItem) error {

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile("clocked.json", data, 0644)
	if err != nil {
		return err
	}

	return nil

}

func readFromFile(input string) {

	file_contents, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	var clocked ClockedItem
	if err := json.Unmarshal(file_contents, &clocked); err != nil {

		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	fmt.Println(clocked)

	//kan bare lese et element enn så lenge

}

func getFromFile(input string) ClockedItem {

	file_contents, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	var clocked ClockedItem
	if err := json.Unmarshal(file_contents, &clocked); err != nil {

		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return clocked

	//kan bare lese et element enn så lenge

}

func clockIn(input string) {

	editable := getFromFile(input)

	if editable.Frequency <= 0 {
		fmt.Println("too smol")
		return
	}

	editable.Frequency--

	dataEntry(editable)

	fmt.Println("Clocked in !")

	readFromFile(input)

}

func main() {

	dataEntry(makeClock("Krita", 3))

	//readFromFile("clocked.json")
	clockIn("clocked.json")

}

//TODO: Fiks litt mer feilmeldinger, implementer lister og appending, senere tid og sånn
