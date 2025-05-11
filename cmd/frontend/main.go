package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Gokeep")

	// Suchfeld
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search...")

	// Neue Notiz Button
	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		// Hier später Dialog zum Erstellen
	})

	// Kopfzeile mit Suche und + Button
	header := container.NewBorder(nil, nil, widget.NewButtonWithIcon("", theme.MenuIcon(), func() {}), addButton, searchEntry)

	// Beispielkarten (dynamisch später per HTTP aus Backend laden)
	notes := []string{"Title 1", "Title 2", "Title 3"}
	var cards []fyne.CanvasObject

	for _, title := range notes {
		card := widget.NewCard(title, "Important notes...\nImportant notes...\nImportant notes...",
			container.NewHBox(
				widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}),
				widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {}),
			),
		)
		cards = append(cards, card)
	}

	// Grid mit Karten
	grid := container.NewGridWithColumns(3, cards...)

	// Hauptlayout
	content := container.NewBorder(header, nil, nil, nil, grid)

	w.SetContent(content)
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
