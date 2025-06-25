package main

import (
	"encoding/json"
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

func clockIn(input string) error {

	//stringen som blir tatt inn er kun fil, fiks dette så det er også fra filen.

	//readFromFile(input)

	editable := getFromFile(input)

	editable.Frequency--
	editable.Frequency--
	editable.Frequency--
	//editable.Frequency--

	if editable.Frequency < 0 {
		//fmt.Println("too smol")
		return fmt.Errorf("Too smol")
		//return fmt.Errorf("Too smol")
	}

	dataEntry(editable)

	fmt.Println("Clocked in !")

	readFromFile(input)
	return nil

	//clocked in spør etter en spesifikk fil, den må også spørre om en aktivitet. Lurer på om interfaces og slikt her.

}

func main() {

	dataEntry(makeClock("Krita", 3))

	clockIn("clocked.json")

	//hva er det du clocker inn?, clockIn bør ta en activity

}

//TODO: Fiks litt mer feilmeldinger, implementer lister og appending, senere tid og sånn

//filen bør være statisk kanskje til og med en global variabel. input bør være dynamisk clockin bør ta et input, sjekke i filen om det fins, og så dekrementere det
