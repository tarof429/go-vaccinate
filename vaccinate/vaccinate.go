package vaccinate

import (
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

// Person is someone who can get sick.  A sick Person has a 10% chance of infecting up to 4 other people every 5 seconds.
// Infections are manifested by symptoms. If the person has 3 symptoms then he will
// call a help line. Once he is administered a test he will be diagnosed as being "sick" and
// go into quarantine for 1 minute. At the end of the quarantine he will lose all his symptoms.
// However after going out of quarantine he may be infected again.
type Person struct {
	ID           int
	SickDay      int
	age          int
	InfectedFlag bool
	attributes   *PersonListAttributes
}

// PersonNode is a node in a circular list
type PersonNode struct {
	person   Person
	next     *PersonNode
	previous *PersonNode
}

// PersonListAttributes are attributes of the Personlist
type PersonListAttributes struct {
	maxSneeze         int
	infectionRate     int
	maxSickDays       int
	numberOfPeople    int
	infectedCount     int
	sneezeProbability *rand.Rand
}

// PersonList is a list of Persons
type PersonList struct {
	attributes PersonListAttributes
	head       *PersonNode
	tail       *PersonNode
}

func (list *PersonList) newPerson(id int, sickDay int, age int, infectedFlag bool) Person {
	return Person{id, sickDay, age, infectedFlag, &list.attributes}
}

func newPersonNode(person Person) *PersonNode {
	return &PersonNode{person, nil, nil}
}

func newPersonList(maxSneeze, infectionRate, maxSickDays, numberOfPeople, infectedCount int) *PersonList {

	attributes := PersonListAttributes{
		maxSneeze,
		infectionRate,
		maxSickDays,
		numberOfPeople,
		infectedCount,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	list := PersonList{
		attributes,
		nil,
		nil,
	}

	for i := 0; i < numberOfPeople; i++ {
		p := list.newPerson(i, 0, 0, false)
		list.add(p)
	}

	return &list
}

func (list *PersonList) add(p Person) {

	node := newPersonNode(p)

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

func (list *PersonList) visit(times int) {

	cur := list.head

	iteration := 0

	for cur != nil {

		if times != 0 && iteration > times {
			break
		}

		list.attributes.infectedCount = 0

		cur.epoch()

		if cur.person.InfectedFlag == true {
			list.attributes.infectedCount++
		}
		cur = cur.next

		if times != 0 {
			iteration++
		}
	}
}

func (list *PersonList) reverseVisit() {

	cur := list.head

	for cur != nil {
		fmt.Println(cur.person)
		cur = cur.previous
		time.Sleep(100 * time.Millisecond)
	}
}

func (p Person) String() string {
	return fmt.Sprintf("ID: %d, Infected: %v, SickDay: %d, Age: %d", p.ID, p.InfectedFlag, p.SickDay, p.age)
}

func (node *PersonNode) epoch() {

	if node.person.InfectedFlag == true {

		listAttributes := node.person.attributes

		var probability int

		// Previous
		probability = listAttributes.sneezeProbability.Intn(100)

		if probability <= listAttributes.infectionRate {
			node.previous.person.InfectedFlag = true
		}

		// Next
		probability = listAttributes.sneezeProbability.Intn(100)

		if probability <= listAttributes.infectionRate {
			node.next.person.InfectedFlag = true
			node.person.SickDay++
		}

		if node.person.SickDay > node.person.attributes.maxSickDays {
			node.person.InfectedFlag = false
			node.person.SickDay = 0
		}
	}
}

func (list *PersonList) gatherStats() {

	cur := list.head

	headAddr := list.head

	list.attributes.infectedCount = 0

	for cur != nil {
		if cur.person.InfectedFlag == true {
			list.attributes.infectedCount++
		}

		//fmt.Println(cur.person)
		cur = cur.next

		if cur == headAddr {
			break
		}

	}
}

func (list *PersonList) printStats() {

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)

	defer w.Flush()

	show := func(a, b interface{}) {
		fmt.Fprintf(w, "%v\t%v\n", a, b)
	}

	show("COLUMN", "VALUE")
	show("People", list.attributes.numberOfPeople)
	show("Infected", list.attributes.infectedCount)

}

func sleep() {
	time.Sleep(time.Nanosecond * 1000)
}

// Run runs the simulation
func Run() {
	persons := newPersonList(3, 10, 3, 100, 0)

	persons.head.person.InfectedFlag = true

	persons.visit(10000)
	//persons.list()
	persons.gatherStats()
	persons.printStats()
}
