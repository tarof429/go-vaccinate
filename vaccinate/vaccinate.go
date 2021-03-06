package vaccinate

import (
	"fmt"
	"math/rand"
	"time"
)

const configFile = ".vaccinate"

// Person is someone who can get sick.
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
	CommonName        string
	InfectionRate     int
	MaxSickDays       int
	NumberOfPeople    int
	Visits            int
	sneezeProbability *rand.Rand
	stats             PersonListStats
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

// InfectionInfo is a struct for normalizing the simulation result data
type InfectionInfo struct {
	CommonName       string
	Total            int
	Visits           int
	InfectionRate    int
	InfectedCount    int
	NumberInfections int
	NumberCured      int
}

// infectTheHead infects the first person in the list
func (list *PersonList) infectTheHead() {
	list.head.person.infect()
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

// visit traverses the list one time and applies an epoch() to each node.
// Each epoch will have a chance to infect adjacent nodes but only if the current node
// is sick.
func (list *PersonList) visit() {

	cur := list.head

	for cur != nil && cur != list.tail {
		list.epoch(cur)
		cur = cur.next
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

	maxSickDays := list.attr.MaxSickDays
	infectionRate := list.attr.InfectionRate
	probability := list.attr.sneezeProbability.Intn(100)

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

// resetStats resets the list statistics. This method must be caled prior to calling GatherStats(),
func (list *PersonList) resetStats() {
	list.attr.stats = PersonListStats{0, 0, 0}
}

// InfectionInfo iterates through the list and returns InfectionInfo
func (list *PersonList) InfectionInfo() InfectionInfo {

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

	return InfectionInfo{
		CommonName:       list.attr.CommonName,
		Total:            list.attr.NumberOfPeople,
		Visits:           list.attr.Visits,
		InfectionRate:    list.attr.InfectionRate,
		InfectedCount:    list.attr.stats.infectedCount,
		NumberInfections: list.attr.stats.numberOfTimesInfected,
		NumberCured:      list.attr.stats.numberOfTimesCured,
	}
}

// Simulator is something that loads a configuration and runs a simulation
type Simulator interface {
	Run()
	Load(dir string) error
}
