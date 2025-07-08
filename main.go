package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var filepath = "clocked.json"
var delete_empty = true

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

}

func getFromFile(input string, m map[string]ClockedItem) ClockedItem {

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

func getMapfromFile() map[string]ClockedItem {

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
	//decrement func
	editable := getFromFile(input, data)

	if editable.Frequency <= 0 {
		return fmt.Errorf("too small")
	}

	editable.Frequency--
	data[editable.Activity] = editable

	if editable.Frequency == 0 && delete_empty {

		delete(data, input)

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

func listAll(m map[string]ClockedItem) {

	println("Activity|Frequency")
	for activ, freq := range m {

		fmt.Printf("%.20s\t%d\n", activ, freq.Frequency)

		//ser noenlunde grei ut, finn lit mer ut med formatering.
		//bør og se hvordan man kan print ut mer når structen blir større m/timestamps etc

	}

	fmt.Printf("OPS: Delete empty is %t \n", delete_empty)
}

func main() {

	lookupMap := getMapfromFile()

	//generiske verdier for flaggene
	activity := ""
	frequency := 0
	clocking := ""
	printelement := ""
	list_all := false
	delete_flag := ""

	flag.StringVar(&activity, "a", "myActivity", "Add activity or update activity")
	flag.StringVar(&delete_flag, "d", "", "Delete activity")
	flag.StringVar(&clocking, "k", "", "Clock Activity")
	flag.IntVar(&frequency, "f", 0, "Sets frequency for desired activity")
	flag.StringVar(&printelement, "p", "", "Prints activity")
	flag.BoolVar(&list_all, "ls", false, "Pretty prints the todolist")
	flag.BoolVar(&delete_empty, "sde", false, "Set delete when activity empty")

	//bruker default verdien når man klokker inn, clocking må kanskje ta inn andre args

	flag.Parse()

	if activity != "" && frequency != 0 {

		interactiveInit(activity, frequency, lookupMap)
		//bruke dataentry her istedet?

	}

	//clocked og freq samtidig kan bli sketch

	if clocking != "" {

		clockIn(clocking, lookupMap)

	}

	if delete_flag != "" {

		delete(lookupMap, delete_flag)
		dataEntry(lookupMap)
		fmt.Printf("Deleted %s \n", delete_flag)

		//her bør det nok være noe feilhåndtering

	}

	if printelement != "" {

		//kanskje ha noe listing av json filen på forhånd.

		//readFromFile(lookupMap)
		a := lookupMap[printelement]
		fmt.Printf("Activity: %s \nFrequency: %d\n", a.Activity, a.Frequency)

	}

	if list_all {

		listAll(lookupMap)

	}

	//capitalize eller standardiser før en check.

}

//TODO: Fiks litt mer feilmeldinger, implementer lister og appending, senere tid og sånn
//TODO: Implementer Receiver funksjoner.
//TODO: User Input --Done! Basic flaggtyper er fikset, negative verdier ser ut til å bli håndtert.
//TODO:Hvis fil ikke finnes, lag den. --Done! os.writefile ser ut til å håndtere dette, kan ta en titt på custom filer senere.
//TODO: Print flagg/funksjon, spesifikt og alt --Done! Tanker videre: Ikke skalerbar til evt større structs
//TODO: Multi input, sjekk om det er forskjell i stor og lite flagg. Per nå tar den kun siste inputet som gyldig.
//TODO: Rename activity?
//Settings: Delete activity when done? Temp flagg?
//Legge inn timestamps
//Installscript, kompilasjon, etter å ha fikset feilmeldinger, 0.1 bør bli ferdig snart.
//REPL i 0.2 med en kortversjon (rediger et enkeltelement, og en komplett en)
//unngå dobbel/trippelprinting, alt bør bruke samme print funksjon, og kun gjøre det engang.
//Cutoff på tekst/formatering.
//Priority?
//Grab bag
//Cleanup etter lang tid
//tab completion?
// default til 1? + hvis man ikke setter flagg og den ikke fins registrer den, eller clock inn hvis den fins. er dette dårlig oppførsel?
//bug? Hva er my acitivity?
//spaces i string i aktivitet? er dette et problem med hvordan map funker?
//waterfall? dependent activities
