package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// categories like "Protein bar" going into specific types

/*
MVC... kind of
described as:
    Init, a function that returns an initial command for the application to run.
    Update, a function that handles incoming events and updates the model accordingly.
    View, a function that renders the UI based on the data in the model.
but it also includes the MODEL of MVC, which holds data for the entire application

prototyping
    date: struct ->
        year
        yearDay
    stock: struct ->
        item name (string)
        amount (int)
        consumption rate (float) <-- calculated by computer, per day
        estimated consumption rate (float) <-- input by user, per day
        last modified (date struct) <-- offer y/N for whether or not to update this value when modified
    model: struct ->
        stocks (array of stock struct)
        cursor (int)
        altCursor (int)
*/

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type date struct {
	year    int
	yearDay int // time package: YearDay returns the day of the year specified by t in the range [1,365]
}

type stock struct {
	name            string
	amount          int
	consumptionRate float32
	lastModified    date
	creationDate    date
}

type model struct {
	stocks       []stock
	cursor       int
	altCursor    int
	cursorActive int // 0 = main; 1 = alt
	table        table.Model
}

func initialModel() model {
	m := model{
		// choices: []string{"buy 2", "buy 2", "buy 3"},
		stocks: []stock{
			stock{
				name:            "Chocolate Chip Muffins",
				amount:          8,
				consumptionRate: 1.5,
				lastModified:    date{year: 2024, yearDay: 200},
				creationDate:    date{year: 2024, yearDay: 150},
			},
			stock{
				name:            "stock2",
				amount:          15,
				consumptionRate: 1.5,
				lastModified:    date{year: 1, yearDay: 300},
				creationDate:    date{year: 1, yearDay: 200},
			},
		},
		cursor:       0,
		altCursor:    0,
		cursorActive: 0,
	}

	columns := []table.Column{
		{Title: "index", Width: 5},
		{Title: "item", Width: 30},
		{Title: "amount", Width: 6},
	}

	rows := []table.Row{}

	for i := 0; i < len(m.stocks); i++ {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			m.stocks[i].name,
			fmt.Sprintf("%d", m.stocks[i].amount),
		})
		// rows = append({"i + 1", m.stocks[i].name, "m.stocks[i].amount"})
	}
	// rows := []table.Row{
	//     {"1", "Chocolate Chip Muffins", "8"},
	//     {"2", "stock2", "15"},
	// }

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	m.table = t
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// msg being any input

	var cmd tea.Cmd

	// what was the type of the input
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursorActive == 0 && m.cursor > 0 {
				m.cursor--
				m.table, cmd = m.table.Update(msg)
			} else if m.cursorActive == 1 && m.altCursor > 0 {
				m.altCursor--
			}

		case "down", "j":
			if m.cursorActive == 0 && m.cursor < len(m.stocks)-1 {
				m.cursor++
				m.table, cmd = m.table.Update(msg)
			} else if m.cursorActive == 1 && m.altCursor < 5 {
				m.altCursor++
			}

		case "left", "h", "right", "l":
			// m.cursorActive = 0
			if m.cursorActive == 0 {
				m.cursorActive = 1
			} else if m.cursorActive == 1 {
				m.cursorActive = 0
			}

			// case "right", "l":
			//     m.cursorActive = 1

		}

	}

	return m, cmd
}

func (m model) View() string {

	line := "inventory-stock manager\n"

	// for i, currentStock := range m.stocks {
	//
	//     cursor := " "
	//     if m.cursor == i {
	//         cursor = ">"
	//     }
	//
	//     line += fmt.Sprintf("%s [%s]\n", cursor, currentStock.name)
	//
	// }

	line += baseStyle.Render(m.table.View())

	creationDiff := (time.Now().Year()*365 + time.Now().YearDay()) - (m.stocks[m.cursor].creationDate.year*365 + m.stocks[m.cursor].creationDate.yearDay)
	modifyDiff := (time.Now().Year()*365 + time.Now().YearDay()) - (m.stocks[m.cursor].lastModified.year*365 + m.stocks[m.cursor].lastModified.yearDay)

	line += fmt.Sprintf("\nCreated on: %d days ago\nLast modified: %d days ago", creationDiff, modifyDiff)

	line += "\n\n<Q> Quit | <H/J/K/L> Move cursor | <A> Add New Item | <D> Delete Item | <Z/X> Decrement/Increment\n"

	return line
	// return baseStyle.Render(m.table.View()) + "\n " + m.table.HelpView() + "\n"

}

func main() {

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

}
