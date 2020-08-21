package vaccinate

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// const (
// 	n int = 1000
// )

// TerminalSimulator is a kind of Simulator for a terminal
type TerminalSimulator struct {
	list     *PersonList
	messages chan InfectionInfo
}

// Load loads the configuration file under dir and populates the list attributes
func (s *TerminalSimulator) Load(dir string) error {

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

// Run runs the simulation. This implementation will loop forever if Visits is set to 0.
func (s *TerminalSimulator) Run() {

	if s.list.head == nil {
		err := errors.New("List is empty")
		fmt.Println(err.Error())
		return
	}

	s.list.infectTheHead()
	s.list.resetStats()

	s.messages = make(chan InfectionInfo)

	for i := 0; s.list.attr.Visits == 0 || i < s.list.attr.Visits; {

		go func() {
			s.list.visit()
			time.Sleep(time.Millisecond * 1000)
			data := s.list.InfectionInfo()
			s.messages <- data
		}()

		data := <-s.messages

		printInfo(data)

		if s.list.attr.Visits != 0 {
			i++
		}
	}

}
