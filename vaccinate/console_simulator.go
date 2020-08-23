package vaccinate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"
)

// ConsoleSimulator is a kind of Simulator for the console
type ConsoleSimulator struct {
	list *PersonList
}

// printInfo is used to print the statistics in a columnar format
func printInfo(info InfectionInfo) {

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)

	defer w.Flush()

	show := func(a, b interface{}) {
		fmt.Fprintf(w, "%v\t%v\n", a, b)
	}

	show("COLUMN", "VALUE")
	show("Common name", info.CommonName)
	show("People", info.Total)
	show("Visits", info.Visits)
	show("Infection rate", info.InfectionRate)
	show("Infected count", info.InfectedCount)
	show("Number of  times infected", info.NumberInfections)
	show("Number of times cured", info.NumberCured)
}

// defaultPersonListAttributes returns a *PersonListAttributes with default values
func defaultPersonListAttributes() *PersonListAttributes {
	return &PersonListAttributes{InfectionRate: 10, MaxSickDays: 3, Visits: 10000, NumberOfPeople: 100}
}

// writeConfig writes PersonListAttributes to the config file under dir
func writeConfig(dir string, attr *PersonListAttributes) error {
	sep := string(os.PathSeparator)

	path := dir + sep + configFile

	data, err := json.MarshalIndent(attr, "", "\t")

	if err != nil {
		log.Fatalf(err.Error())
	}

	mode := int(0644)

	err = ioutil.WriteFile(path, data, os.FileMode(mode))

	return err
}

// readConfig reads a config file under dir and populates the list attributes
func readConfig(dir string, attr *PersonListAttributes) error {
	sep := string(os.PathSeparator)

	f, err := ioutil.ReadFile(dir + sep + configFile)

	if err != nil {
		log.Fatalf(err.Error())
	}

	//fmt.Println("Unmarshalling")
	err = json.Unmarshal(f, attr)

	//fmt.Println("Returning from ReadConfig")
	return err
}

// Load loads the configuration file under dir and populates the list attributes
func (s *ConsoleSimulator) Load(dir string) error {

	sep := string(os.PathSeparator)

	path := dir + sep + configFile

	_, err := os.Stat(path)

	if err != nil {
		err = writeConfig(dir, defaultPersonListAttributes())

		if err != nil {
			return err
		}
	}
	attr := PersonListAttributes{}

	err = readConfig(dir, &attr)

	if err != nil {
		return err
	}

	s.list = newPersonList(&attr)

	return nil

}

// Run runs the simulation based on the provided attributes.
// The first person will be infected by default.
// This function is useful for running the simulation in console mode when
// only the results are desired.
func (s *ConsoleSimulator) Run() {

	if s.list.head == nil {
		err := errors.New("List is empty")
		fmt.Println(err.Error())
		return
	}

	s.list.resetStats()
	s.list.infectTheHead()
	s.list.visit()
	printInfo(s.list.InfectionInfo())

	// for {
	// 	s.list.visit()
	// 	s.list.gatherStats()
	// 	s.list.printStats()
	// 	s.list.resetStats()
	// 	time.After(time.Second)
	// 	//return nil
	// }
}
