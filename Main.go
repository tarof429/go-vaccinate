package main

import (
	"fmt"
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
	InfectedFlag bool
}

type PersonNode struct {
	person Person
	next   *PersonNode
}

type PersonList struct {
	maxSneeze     int
	infectionRate int
	maxSickDays   int
	head          *PersonNode
	tail          *PersonNode
}

func NewPerson(id int, sickDay int, infectedFlag bool) Person {
	return Person{id, sickDay, infectedFlag}
}

func NewPersonNode(person Person) *PersonNode {
	return &PersonNode{person, nil}
}

func NewPersonList(maxSneeze, infectionRate, maxSickDays, count int) *PersonList {
	list := PersonList{maxSneeze, infectionRate, maxSickDays, nil, nil}

	for i := 0; i < count; i++ {
		p := NewPerson(i, 0, false)
		list.Add(p)
	}

	return &list
}

func (list *PersonList) Add(p Person) {

	node := NewPersonNode(p)

	if list.head == nil {
		list.head = node
		node.next = node
		list.tail = node
	} else {
		node.next = list.head
		list.tail.next = node
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

func (list *PersonList) Visit() {

	cur := list.head

	for cur != nil {
		fmt.Println(cur.person)
		cur = cur.next
		time.Sleep(100 * time.Millisecond)
	}
}

func (p Person) String() string {
	return fmt.Sprintf("ID: %d, Infected: %v, SickDay: %d, ", p.ID, p.InfectedFlag, p.SickDay)
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

// func Epoch(persons []Person) {

// 	for _, person := range persons {
// 		person.Sneeze()

// 		fmt.Println(person)
// 	}

// }

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "load", Description: "Load simulation"},
		{Text: "run", Description: "Run simulation"},
		{Text: "quit", Description: "Quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// func PersonList() []Person {

// 	ids := r.Perm(MaxPeople)

// 	persons := make([]Person, MaxPeople)

// 	for i, id := range ids {
// 		persons[i] = Person{ID: id}
// 	}

// 	return persons
// }

func Sleep() {
	time.Sleep(time.Nanosecond * 1000)
}

func main() {

	// Init()

	persons := NewPersonList(3, 3, 3, 10)

	persons.List()

	// Sleep()

	// for _, person := range persons {
	// 	person.Sneeze()
	// 	fmt.Println(person)
	// }

	// for _, person := range persons {
	// 	fmt.Println(person)
	// }

	// f := Facility{
	// 	id:           "1234",
	// 	maxOccupancy: 4,
	// 	patients:     []Patient{p},
	// }

	// for {
	// 	fmt.Println("Please select command")
	// 	t := prompt.Input("> ", completer)
	// 	fmt.Println("You selected " + t)

	// 	if strings.ToLower(t) == "quit" {
	// 		break
	// 	}
	// }
}
