package vaccinate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

const configFile = ".vaccinate"

// Person is someone who can get sick.  A sick Person has a 10% chance of infecting up to 4 other people every 5 seconds.
// Infections are manifested by symptoms. If the person has 3 symptoms then he will
// call a help line. Once he is administered a test he will be diagnosed as being "sick" and
// go into quarantine for 1 minute. At the end of the quarantine he will lose all his symptoms.
// However after going out of quarantine he may be infected again.
type Person struct {
	id                    int
	sickDay               int
	infectedFlag          bool
	numberOfTimesInfected int
	numberOfTimesCured    int
}

func (p Person) infected() bool {
	return p.infectedFlag == true
}

func (p *Person) infect() {
	p.infectedFlag = true
	p.numberOfTimesInfected++
}

func (p *Person) disinfect() {
	p.infectedFlag = false
	p.numberOfTimesCured++
}

func (p *Person) stayInBed() {
	p.sickDay++
}

func (p *Person) getWell() {
	p.sickDay = 0
}

// PersonNode is a node in a circular list
type PersonNode struct {
	person   Person
	next     *PersonNode
	previous *PersonNode
}

// PersonListAttributes are attributes of the Personlist
type PersonListAttributes struct {
	InfectionRate     int
	MaxSickDays       int
	NumberOfPeople    int
	Visits            int
	stats             PersonListStats
	sneezeProbability *rand.Rand
}

// PersonListStats are statistics about the simulation.
// It needs to be initialized before use and populated by traversing the list.
type PersonListStats struct {
	infectedCount         int
	numberOfTimesInfected int
	numberOfTimesCured    int
}

// PersonList is a list of Persons. It is meant to be used as a circular list.
type PersonList struct {
	attr *PersonListAttributes
	head *PersonNode
	tail *PersonNode
}

// newPersonList is a factory function that generates a PersonList based on attributes.
func newPersonList(attr *PersonListAttributes) *PersonList {

	attr.sneezeProbability = rand.New(rand.NewSource(time.Now().UnixNano()))

	list := PersonList{attr, nil, nil}

	for i := 0; i < attr.NumberOfPeople; i++ {
		list.add(Person{i, 0, false, 0, 0})
	}

	return &list
}

// add adds a new person to the list
func (list *PersonList) add(p Person) {

	node := &PersonNode{p, nil, nil}

	if list.head == nil {
		list.head = node
		node.next = list.head
		node.previous = list.tail
		list.head.previous = node
		list.tail = node

	} else {
		node.next = list.head
		node.previous = list.tail
		list.tail.next = node
		list.head.previous = node
		list.tail = node
	}
}

// list lists each person in the list
func (list *PersonList) list() {

	cur := list.head

	headAddr := list.head

	for cur != nil {
		fmt.Println(cur.person)
		cur = cur.next

		if cur == headAddr {
			break
		}

	}
}

//reverseList is the same as list() but traverses the list in reverse
func (list *PersonList) reverseList() {

	cur := list.tail

	headAddr := list.tail

	for cur != nil {
		fmt.Println(cur.person)
		cur = cur.previous

		if cur == headAddr {
			break
		}
	}
}

// visit traverses the list and applies an epoch() to each node
func (list *PersonList) visit() {

	if list.attr.Visits == 0 {
		return
	}

	cur := list.head
	iteration := 0

	for cur != nil && iteration < list.attr.Visits {
		list.epoch(cur)
		cur = cur.next
		iteration++

	}
}

// epoch gives the possibility of each person sneezing on its neighbor
func (list *PersonList) epoch(node *PersonNode) {
	if node.person.infected() {
		list.sneeze(node.previous)
		list.sneeze(node.next)
	}
}

// reverseList prints the list in reverse
func (list *PersonList) reverseVisit() {

	cur := list.head

	for cur != nil {
		fmt.Println(cur.person)
		cur = cur.previous
		time.Sleep(100 * time.Millisecond)
	}
}

// String is a stringer function used to print a person
func (p Person) String() string {
	return fmt.Sprintf("ID: %d, Infected: %v, Number of times infected: %d", p.id, p.infectedFlag, p.numberOfTimesInfected)
}

// sneeze is a method that a person invokes if he's infected
func (list *PersonList) sneeze(on *PersonNode) {

	sneezeProbabilty := list.attr.sneezeProbability
	maxSickDays := list.attr.MaxSickDays
	infectionRate := list.attr.InfectionRate
	probability := sneezeProbabilty.Intn(100)

	if on.person.infected() == false {
		if probability <= infectionRate {
			on.person.infect()
			on.person.getWell()
		}
	} else {
		on.person.stayInBed()

		if on.person.sickDay > maxSickDays {
			on.person.disinfect()
			on.person.getWell()
		}
	}
}

// ResetStats resets the list statistics. This method must be caled prior to calling GatherStats(),
func (list *PersonList) ResetStats() {
	list.attr.stats = PersonListStats{0, 0, 0}
}

// GatherStats iterates through the list and gathers statistics
func (list *PersonList) GatherStats() {

	cur := list.head

	for cur != nil {
		if cur.person.infected() {
			list.attr.stats.infectedCount++
		}
		list.attr.stats.numberOfTimesInfected += cur.person.numberOfTimesInfected
		list.attr.stats.numberOfTimesCured += cur.person.numberOfTimesCured

		//fmt.Printf("Number of times infected vs cured (per person): %d (total) %d infected: %v\n", cur.person.numberOfTimesInfected, cur.person.numberOfTimesCured, cur.person.infected())
		//fmt.Println(cur.person)

		cur = cur.next

		if cur == list.head {
			break
		}
	}
}

// PrintStats is used to print the statistics in a columnar format
func (list *PersonList) PrintStats() {

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

func sleep() {
	time.Sleep(time.Nanosecond * 1000)
}

// DefaultPersonListAttributes returns a *PersonListAttributes with default values
func DefaultPersonListAttributes() *PersonListAttributes {
	return &PersonListAttributes{InfectionRate: 10, MaxSickDays: 3, Visits: 10000, NumberOfPeople: 100}
}

// WriteConfig writes PersonListAttributes to the config file under dir
func WriteConfig(dir string, attr *PersonListAttributes) error {
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

// ReadConfig reads a config file under dir and populates the list attributes
func ReadConfig(dir string, attr *PersonListAttributes) error {
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
func Load(dir string, attr *PersonListAttributes) error {

	sep := string(os.PathSeparator)

	path := dir + sep + configFile

	_, err := os.Stat(path)

	if err != nil {
		err = WriteConfig(dir, DefaultPersonListAttributes())

		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	err = ReadConfig(dir, attr)

	return err

}

// Run runs the simulation based on the provided attributes.
// The first person will be infected by default.
// This function is useful for running the simulation in console mode when
// only the results are desired.
func Run(attr *PersonListAttributes) error {

	persons := newPersonList(attr)

	if persons.head == nil {
		return errors.New("List is empty")
	}

	persons.head.person.infect()
	persons.visit()
	persons.ResetStats()
	persons.GatherStats()
	persons.PrintStats()

	return nil
}
