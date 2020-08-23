package vaccinate

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// TerminalSimulator is a kind of Simulator for a terminal
type TerminalSimulator struct {
	list     *PersonList
	messages chan InfectionInfo
}

var (
	columnAWidth = 81
	columnBWidth = 27
	rectx        = 110
)

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

func setupSimulator(s *TerminalSimulator) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	counter := 0

	s.list.visit()
	data := s.list.InfectionInfo()

	table1 := widgets.NewTable()

	table1.Rows = [][]string{
		{"Column", "Value"},
		{"People", strconv.Itoa(data.Total)},
		{"Visit", strconv.Itoa(data.Visits)},
		{"Infection Rate", strconv.Itoa(data.InfectionRate)},
		{"Infection Count", strconv.Itoa(data.InfectedCount)},
		{"Number of times infected", strconv.Itoa(data.NumberInfections)},
		{"Number of times cured", strconv.Itoa(data.NumberCured)},
		{"ColumnA width", strconv.Itoa(columnAWidth)},
		{"ColumnB width", strconv.Itoa(columnBWidth)},
		{"Rectx", strconv.Itoa(rectx)},
	}
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.ColumnWidths = []int{columnAWidth, columnBWidth}
	table1.RowSeparator = false
	table1.Border = true
	table1.BorderStyle.Fg = ui.ColorBlue
	table1.FillRow = true
	table1.Block.Title = " Statistics "
	table1.SetRect(0, 0, rectx, 12)

	table1.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)

	lc2 := widgets.NewPlot()
	lc2.Title = " Infection Count "
	lc2.BorderStyle.Fg = ui.ColorBlue
	lc2.Data = make([][]float64, 1)
	lc2.MaxVal = float64(data.Total)
	lc2.Data[0] = make([]float64, 1)

	lc2.Data[0][counter] = float64(data.InfectedCount)
	lc2.SetRect(0, 50, 110, 15)
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColors[0] = ui.ColorYellow
	lc2.PlotType = widgets.ScatterPlot

	draw := func() {

		s.list.resetStats()

		for i := 0; i < data.Visits; i++ {
			s.list.visit()
		}
		data := s.list.InfectionInfo()
		table1.Rows = [][]string{
			{"Column", "Value"},
			{"People", strconv.Itoa(data.Total)},
			{"Visit", strconv.Itoa(data.Visits)},
			{"Infection Rate", strconv.Itoa(data.InfectionRate)},
			{"Infection Count", strconv.Itoa(data.InfectedCount)},
			{"Number of times infected", strconv.Itoa(data.NumberInfections)},
			{"Number of times cured", strconv.Itoa(data.NumberCured)},
			{"ColumnA width", strconv.Itoa(columnAWidth)},
			{"ColumnB width", strconv.Itoa(columnBWidth)},
			{"Rectx", strconv.Itoa(rectx)},
		}
		// Append to the slice so that we have more data points to display
		lc2.Data[0] = append(lc2.Data[0], float64(data.InfectedCount))
		ui.Render(table1, lc2)
	}

	draw()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Millisecond * 10).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "a":
				columnAWidth++
			case "s":
				columnAWidth--
			case "d":
				columnBWidth++
			case "f":
				columnBWidth--
			case "g":
				rectx--
			case "h":
				rectx++
			}
		case <-ticker:

			counter++
			if counter >= data.Total {
				time.Sleep(time.Second * 1)
				s.reset()
				counter = 0
				lc2.Data[0] = nil
			}
			draw()
		}
	}
}

// reset the simulation
func (s *TerminalSimulator) reset() {

	s.list = newPersonList(s.list.attr)
	s.list.infectTheHead()

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

	setupSimulator(s)

}
