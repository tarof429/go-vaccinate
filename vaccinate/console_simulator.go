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
type ConsoleSimulator struct{}

// printStats is used to print the statistics in a columnar format
func (list *PersonList) printStats() {

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)

	defer w.Flush()

	show := func(a, b interface{}) {
		fmt.Fprintf(w, "%v\t%v\n", a, b)
	}

	show("COLUMN", "VALUE")
	show("People", list.attr.NumberOfPeople)
	show("Visits", list.attr.Visits)
	show("Infection rate", list.attr.InfectionRate)
	show("Infected count", list.attr.stats.infectedCount)
	show("Number of  times infected", list.attr.stats.numberOfTimesInfected)
	show("Number of times cured", list.attr.stats.numberOfTimesCured)
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
func (s ConsoleSimulator) Load(dir string, attr *PersonListAttributes) error {

	sep := string(os.PathSeparator)

	path := dir + sep + configFile

	_, err := os.Stat(path)

	if err != nil {
		err = writeConfig(dir, defaultPersonListAttributes())

		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	err = readConfig(dir, attr)

	return err

}

// Run runs the simulation based on the provided attributes.
// The first person will be infected by default.
// This function is useful for running the simulation in console mode when
// only the results are desired.
func (s ConsoleSimulator) Run(attr *PersonListAttributes) error {

	persons := newPersonList(attr)

	if persons.head == nil {
		return errors.New("List is empty")
	}

	persons.head.person.infect()
	persons.visit()
	persons.resetStats()
	persons.gatherStats()
	persons.printStats()

	return nil
}
