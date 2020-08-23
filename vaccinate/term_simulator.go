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
	simtablecolumnAWidth = 81
	simtablecolumnBWidth = 27

	simtablex1 = 0
	simtabley1 = 0
	simtablex2 = 110
	simtabley2 = 14 // used 9 previously

	plotx1 = 0
	ploty1 = 51
	plotx2 = 110
	ploty2 = 15
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
		log.Fatalf("failed to initialize the simulator: %v", err)
	}
	defer ui.Close()

	counter := 0

	s.list.visit()
	data := s.list.InfectionInfo()

	simTable := widgets.NewTable()

	simTable.Rows = [][]string{
		{"Column", "Value"},
		{"Common Name", data.CommonName},
		{"People", strconv.Itoa(data.Total)},
		{"Visit", strconv.Itoa(data.Visits)},
		{"Infection Rate", strconv.Itoa(data.InfectionRate) + "%"},
		{"Infection Count", strconv.Itoa(data.InfectedCount)},
		{"Number of times infected", strconv.Itoa(data.NumberInfections)},
		{"Number of times cured", strconv.Itoa(data.NumberCured)},
	}
	simTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	simTable.ColumnWidths = []int{simtablecolumnAWidth, simtablecolumnBWidth}
	simTable.RowSeparator = false
	simTable.Border = true
	simTable.BorderStyle.Fg = ui.ColorBlue
	simTable.FillRow = true
	simTable.Block.Title = " Statistics "
	simTable.SetRect(simtablex1, simtabley1, simtablex2, simtabley2)

	simTable.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)

	simPlot := widgets.NewPlot()
	simPlot.PlotType = widgets.ScatterPlot
	simPlot.Title = " Infection Count "
	simPlot.BorderStyle.Fg = ui.ColorBlue
	simPlot.Data = make([][]float64, 1)
	simPlot.MaxVal = float64(data.Total)
	simPlot.Data[0] = make([]float64, 1)
	simPlot.Data[0][counter] = float64(data.InfectedCount)
	simPlot.SetRect(plotx1, ploty1, plotx2, ploty2)
	simPlot.AxesColor = ui.ColorWhite
	simPlot.LineColors[0] = ui.ColorYellow

	draw := func() {

		s.list.resetStats()

		for i := 0; i < data.Visits; i++ {
			s.list.visit()
		}
		data := s.list.InfectionInfo()
		simTable.Rows = [][]string{
			{"Column", "Value"},
			{"Common Name", data.CommonName},
			{"People", strconv.Itoa(data.Total)},
			{"Visit", strconv.Itoa(data.Visits)},
			{"Infection Rate", strconv.Itoa(data.InfectionRate) + "%"},
			{"Infection Count", strconv.Itoa(data.InfectedCount)},
			{"Number of times infected", strconv.Itoa(data.NumberInfections)},
			{"Number of times cured", strconv.Itoa(data.NumberCured)},
		}
		//lc2.SetRect(plotx1, ploty1, plotx2, ploty2)

		simTable.SetRect(simtablex1, simtabley1, simtablex2, simtabley2)

		// Append to the slice so that we have more data points to display
		simPlot.Data[0] = append(simPlot.Data[0], float64(data.InfectedCount))
		ui.Render(simTable, simPlot)
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
				// case "a":
				// 	tabley1--
				// case "s":
				// 	tabley1++
				// case "d":
				// 	tabley2--
				// case "f":
				// 	tabley2++
			}
		case <-ticker:

			counter++
			if counter >= data.Total {
				time.Sleep(time.Second * 1)
				s.reset()
				counter = 0
				simPlot.Data[0] = nil
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
