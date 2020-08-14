package main

import (
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"

	"github.com/c-bata/go-prompt"
)

//var r *rand.Rand
//var sneezeProbability *rand.Rand

// A sick Person has a 10% chance of infecting up to 4 other people every 5 seconds.
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

type PersonNode struct {
	person   Person
	next     *PersonNode
	previous *PersonNode
}

type PersonListAttributes struct {
	maxSneeze         int
	infectionRate     int
	maxSickDays       int
	numberOfPeople    int
	infectedCount     int
	sneezeProbability *rand.Rand
}

type PersonList struct {
	attributes PersonListAttributes
	head       *PersonNode
	tail       *PersonNode
}

func (list *PersonList) NewPerson(id int, sickDay int, age int, infectedFlag bool) Person {
	return Person{id, sickDay, age, infectedFlag, &list.attributes}
}

func NewPersonNode(person Person) *PersonNode {
	return &PersonNode{person, nil, nil}
}

func NewPersonList(maxSneeze, infectionRate, maxSickDays, numberOfPeople, infectedCount int) *PersonList {

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
		p := list.NewPerson(i, 0, 0, false)
		list.Add(p)
	}

	return &list
}

func (list *PersonList) Add(p Person) {

	node := NewPersonNode(p)

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

func (list *PersonList) List() {

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

func (list *PersonList) GatherStats() {

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

func (list *PersonList) Print() {

	// cur := list.head

	// headAddr := list.head

	// for cur != nil {
	// 	fmt.Println(cur.person)
	// 	cur = cur.next

	// 	if cur == headAddr {
	// 		break
	// 	}

	// }

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)

	defer w.Flush()

	show := func(a, b interface{}) {
		fmt.Fprintf(w, "%v\t%v\n", a, b)
	}

	show("COLUMN", "VALUE")
	show("People", list.attributes.numberOfPeople)
	show("Infected", list.attributes.infectedCount)

}

func (list *PersonList) ReverseList() {

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

func (list *PersonList) Visit(times int) {

	cur := list.head

	iteration := 0

	for cur != nil {

		if times != 0 && iteration > times {
			break
		}

		list.attributes.infectedCount = 0

		cur.Epoch()

		if cur.person.InfectedFlag == true {
			list.attributes.infectedCount++
		}
		cur = cur.next

		if times != 0 {
			iteration++
		}
	}
}

func (list *PersonList) ReverseVisit() {

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

// func Init() {
// 	//r = rand.New(rand.NewSource(time.Now().UnixNano()))
// 	sneezeProbability = rand.New(rand.NewSource(time.Now().UnixNano()))
// }

// func (p *Person) Sneeze() bool {
// 	threshold := MaxPeople - InfectionRate

// 	sickFlag := false

// 	for i := 0; i < MaxCoughs; i++ {
// 		n := coughProbability.Intn(MaxPeople)

// 		if n > threshold {
// 			p.SymptomCount++
// 			sickFlag = true

// 		}

// 	}
// 	return sickFlag
// }

func (node *PersonNode) Epoch() {

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

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "load", Description: "Load simulation"},
		{Text: "run", Description: "Run simulation"},
		{Text: "quit", Description: "Quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Sleep() {
	time.Sleep(time.Nanosecond * 1000)
}

func main() {

	persons := NewPersonList(3, 10, 3, 100, 0)

	persons.head.person.InfectedFlag = true

	persons.Visit(10000)
	persons.List()
	persons.GatherStats()
	persons.Print()

	// for {
	// 	fmt.Println("Please select command")
	// 	t := prompt.Input("> ", completer)
	// 	fmt.Println("You selected " + t)

	// 	if strings.ToLower(t) == "quit" {
	// 		break
	// 	}
	// }
}
