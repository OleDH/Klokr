package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

// Burde Clocked item være laveste element? frekvens kan bli flyttet ut til å denotere antall objekter i kategorien.
type ClockedItem struct {
	Activity  string    `json:"activity"`
	Frequency int       `json:"frequency"`
	CreatedAt time.Time `json:"createdat"`
	ClosedAt  time.Time `json:"closedat"`
	Dueby     time.Time `json:"dueby"`
	Priority  int       `json:"priority"`
	Recurring bool      `json:"recurring"`
	Snoozed   int       `json:"snoozed"`
	Baseline  int       `json:"baseline"`
}

//implementerer duebuty og under etter http
//lag en reminder funkjson, dette blir nok interfacing med crontab, kan se på andre måter å purre på bruker, via notifications feks i cinnamon evt.
//bruker nok time std for datoer, skal være snill nok til at gjeldende datoer må være innen midnatt til frist feks.

type Clockhandler struct {
	Data     map[string]ClockedItem `json:"data"`
	JSONpath string
	OptPath  string
}

func MakeclockHandler() *Clockhandler {

	return &Clockhandler{Data: make(map[string]ClockedItem)}

}

func (c *Clockhandler) ensureInit() {
	if c.Data == nil {
		c.Data = make(map[string]ClockedItem)
	}
}

type Settings struct {
	DeleteEmpty bool `json:"delete_empty"`
}

func (c *Clockhandler) LoadSettings() (Settings, error) {
	var s Settings
	b, err := os.ReadFile(c.OptPath)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(b, &s)
}

func (c *Clockhandler) SaveSettings(s Settings) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(c.OptPath, b, 0o644)
}

func makeClock(s string, z int) ClockedItem {

	NewActivity := ClockedItem{
		Activity:  s,
		Frequency: z,
		CreatedAt: time.Now(),
	}
	return NewActivity

}

func userDataDir(app string) (string, error) {
	if d, err := os.UserConfigDir(); err == nil && d != "" {
		// for data prefer XDG_DATA_HOME if set
		if x := os.Getenv("XDG_DATA_HOME"); x != "" {
			return filepath.Join(x, app), nil
		}
		// fallback: put under config dir sibling
		return filepath.Join(d, app), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		return filepath.Join(x, app), nil
	}
	return filepath.Join(home, ".local", "share", app), nil
}

func (c *Clockhandler) dataEntry() error {
	wrapper := struct {
		Data map[string]ClockedItem `json:"data"`
	}{Data: c.Data}

	b, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	return os.WriteFile(c.JSONpath, b, 0o644)
}

func (c *Clockhandler) readFromFile() error {
	// Ensure map exists
	c.ensureInit()

	// If file doesn’t exist, treat as empty
	if _, err := os.Stat(c.JSONpath); os.IsNotExist(err) {
		return nil
	}

	b, err := os.ReadFile(c.JSONpath)
	if err != nil {
		return err
	}

	// Unmarshal into a temporary wrapper matching on-disk shape
	var disk struct {
		Data map[string]ClockedItem `json:"data"`
	}
	if err := json.Unmarshal(b, &disk); err != nil {
		return err
	}

	// Replace current map with loaded one
	c.Data = disk.Data
	if c.Data == nil {
		c.Data = make(map[string]ClockedItem)
	}
	return nil
}

func (c *Clockhandler) clockIn(input string) error {

	c.ensureInit()
	//decrement func
	key := strings.TrimSpace(input)
	item, ok := c.Data[key]
	if !ok {
		return fmt.Errorf("activity not found: %q", key)
	}
	if item.Frequency <= 0 {
		return fmt.Errorf("frequency already zero for %q", key)
	}

	item.Frequency--

	s, _ := c.LoadSettings()

	if item.Frequency <= 0 && s.DeleteEmpty {
		delete(c.Data, key)
	} else if item.Frequency <= 0 && !s.DeleteEmpty {
		item.ClosedAt = time.Now()
		c.Data[key] = item

	} else {
		c.Data[key] = item
	}

	if err := c.dataEntry(); err != nil {
		return err
	}
	fmt.Printf("Clocked in Activity: %s!, Frequency at: %d ", item.Activity, item.Frequency)

	return nil
}

func (C *Clockhandler) interactiveInit(input string, freq int) error {

	C.ensureInit()
	key := strings.TrimSpace(input)
	C.Data[key] = makeClock(key, freq)

	return C.dataEntry()

}

func listAll(m map[string]ClockedItem) {
	fmt.Println("Activity|Frequency")
	keys := keysByFreq(m)
	for _, k := range keys {
		fmt.Printf("%.20s\t%d\n", k, m[k].Frequency)
	}
}

func keysByFreq(m map[string]ClockedItem) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]].Frequency == m[keys[j]].Frequency {
			return keys[i] < keys[j]
		}
		return m[keys[i]].Frequency < m[keys[j]].Frequency
	})
	return keys
}

// flytt server logikk ut i egen fil
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func (c *Clockhandler) httpHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		//return entire json
		respondWithJSON(w, 200, c.Data)
		return

	case "POST":

		defer r.Body.Close()

		var params ClockedItem

		dat, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, 500, "couldn't read request")
			return
		}

		err = json.Unmarshal(dat, &params)
		if err != nil {
			respondWithError(w, 500, "couldn't unmarshal parameters")
			return
		}
		c.Data[params.Activity] = params
		c.dataEntry()

		respondWithJSON(w, 200, params)

	default:
		respondWithError(w, 405, "Method not allowed")
	}
}

func main() {

	dataDir, _ := userDataDir("klokr")
	configDir, _ := os.UserConfigDir()
	appConfigDir := filepath.Join(configDir, "klokr")

	//d := time.Now()
	//println(d.Date())

	if err := ensureDir(dataDir); err != nil {
		log.Fatal(err)
	}
	if err := ensureDir(appConfigDir); err != nil {
		log.Fatal(err)
	}

	handler := MakeclockHandler()
	handler.JSONpath = filepath.Join(dataDir, "clocked.json")
	handler.OptPath = filepath.Join(appConfigDir, "opts.json")

	if err := handler.readFromFile(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/todos", handler.httpHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))

	activity := ""
	frequency := 0
	clocking := ""
	printelement := ""
	printstruct := ""
	list_all := false
	delete_flag := ""
	delete_empty := false
	var listForCompletion bool
	reset := false

	flag.BoolVar(&listForCompletion, "list", false, "list activities for completion")
	flag.BoolVar(&reset, "reset", false, "Delete all data and start fresh")
	flag.StringVar(&activity, "a", "myActivity", "Add activity or update activity")
	flag.StringVar(&delete_flag, "d", "", "Delete activity")
	flag.StringVar(&clocking, "c", "", "Clock Activity")
	flag.IntVar(&frequency, "f", 0, "Sets frequency for desired activity")
	flag.StringVar(&printelement, "p", "", "Prints activity and frequency")
	flag.StringVar(&printstruct, "P", "", "Prints struct")
	flag.BoolVar(&list_all, "ls", false, "Pretty prints the todolist")
	flag.BoolVar(&delete_empty, "sde", false, "Set delete when activity empty")

	flag.Parse()

	if activity != "" && frequency != 0 {
		if err := handler.interactiveInit(activity, frequency); err != nil {
			log.Fatal(err)
		}

	}
	if delete_flag != "" {
		delete(handler.Data, strings.TrimSpace(delete_flag))
		if err := handler.dataEntry(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %s\n", delete_flag)
	}

	if listForCompletion {
		keys := make([]string, 0, len(handler.Data))
		for k := range handler.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Println(k)
		}
		return
	}

	if clocking != "" {

		if err := handler.clockIn(clocking); err != nil {
			log.Fatal(err)
		}

	}

	if delete_flag != "" {
		delete(handler.Data, strings.TrimSpace(delete_flag))
		if err := handler.dataEntry(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %s\n", delete_flag)
	}

	if printelement != "" {

		a := handler.Data[printelement]

		fmt.Printf("Activity: %s \nFrequency: %d\n", a.Activity, a.Frequency)

	}

	if printstruct != "" {

		a := handler.Data[printstruct]

		//fmt.Printf("Activity: %s \nFrequency: %d\n", a.Activity, a.Frequency)
		fmt.Println(a)

	}

	if list_all {

		listAll(handler.Data)

	}

	// ...
	if reset {
		os.Remove(handler.JSONpath)
		fmt.Println("Data reset!")
		return
	}

	//capitalize eller standardiser før en check.

}

//DEMO: implementer tid
// implementer resten av struct .
//update metode
// http server
// sql logikk

//0.1:
//Installscript, kompilasjon, etter å ha fikset feilmeldinger, 0.1 bør bli ferdig snart.
//QA
//

//0.2
//Temp flagg, slettes når tom per struct
//tab completion? Hver gang man cruder, oppdater complete i linux
//Implementer Receiver funksjoner. _/
//Refactor _/
//Priority?
//Legge inn timestamps, tidsenheter o.l
//Rename activity

//0.3
//REPL

//reoccuring som forhindrer sletting

//0.x
//TODO: Multiinput, sjekk om det fins noe ala unix opts.
//TODO: Fiks litt mer feilmeldinger, implementer lister og appending, senere tid og sånn
//Cutoff på tekst/formatering.
//Grab bag, velg noe med lav prioritet, kan kanskje ta en optional arg med så mye tid har jeg idag, fyll inn
//Cleanup etter lang tid
// default til 1? + hvis man ikke setter flagg og den ikke fins registrer den, eller clock inn hvis den fins. er dette dårlig oppførsel?
//spaces i string i aktivitet? er dette et problem med hvordan map funker?
//waterfall? dependent activities, mulig man trenger en egen datastruktur for det? tris?
//må fikse litt quotes før inputs går inn
//kalender/rekkefølgepåting
