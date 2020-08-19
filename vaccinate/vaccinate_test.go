package vaccinate

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

const testdataDir = "testdata"

var s Simulator

func TestMain(m *testing.M) {
	err := os.RemoveAll(testdataDir)

	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(testdataDir, 0755)

	if err != nil {
		log.Fatal(err)
	}

	s = ConsoleSimulator{}

	status := m.Run()

	os.Exit(status)
}

func TestWriteReadConfig(t *testing.T) {
	dir, err := filepath.Abs("testdata")

	if err != nil {
		log.Fatal(err)
	}

	attr := defaultPersonListAttributes()

	err = writeConfig(dir, attr)

	if err != nil {
		t.Fatalf("Unable to write config file: %s", err.Error())
	}

	err2 := readConfig(dir, attr)

	if err2 != nil {
		t.Fatalf("Unable to read config file: %s", err.Error())
	}

	defaultNumberOfPeople := attr.NumberOfPeople

	if defaultNumberOfPeople != attr.NumberOfPeople {
		t.Fatalf("Default number of people was not correct in the read config file")
	}

}

func TestRun(t *testing.T) {
	user, err := user.Current()

	if err != nil {
		log.Fatalf(err.Error())
	}

	var attr PersonListAttributes

	err = s.Load(user.HomeDir, &attr)

	if err != nil {
		log.Fatalf(err.Error())
	}

}

func TestNewPersonList(t *testing.T) {
	list := newPersonList(defaultPersonListAttributes())
	if list.head == nil || list.tail == nil {
		t.Fail()
	}
}
