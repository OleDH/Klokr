package main

import "fmt"

type ClockedItem struct {
	activity  string
	frequency int
}

func main() {
	var (
		a string
		b int
	)
	fmt.Println("Skriv inn ønsket aktivitet")
	fmt.Scanln(&a)
	fmt.Println("Hvor mange ganger i uken vil du gjøre", a, "?")
	fmt.Scanln(&b)
	Testactivity := ClockedItem{activity: a, frequency: b}
	fmt.Println(Testactivity)
	fmt.Println("helloworld")
}

//TODO: Legg til fil init, lagre og lese operasjoner, dekrementeringsoperasjoner
