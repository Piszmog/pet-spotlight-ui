package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"os"
	"pet-spotlight/io"
	"sort"
	"strings"
)

type flags struct {
	dogs             string
	baseDirectory    string
	determineFosters bool
}

func main() {
	// Create app
	mainApp := app.New()
	// Create quit button
	quitButton := widget.NewButton("Quit", func() {
		mainApp.Quit()
	})
	// Create error window
	errorWindow := mainApp.NewWindow("Errors")
	errorEntry := widget.NewEntry()
	errorCloseButton := widget.NewButton("Close", func() {
		errorEntry.SetText("")
		errorWindow.Hide()
	})
	errorWindow.SetContent(widget.NewVBox(errorEntry, errorCloseButton))
	errorChannel := make(chan error, 10)
	go func() {
		for err := range errorChannel {
			errorWindow.Show()
			errorEntry.SetText(errorEntry.Text + "\n" + fmt.Sprintf("%+v", err))
		}
	}()
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		errorEntry.SetText(fmt.Sprintf("%+v", err))
		errorWindow.Show()
		return
	}
	// Create main window
	mainWindow := mainApp.NewWindow("Pet Spotlight")
	// Create directory entry
	baseDirectoryEntry := widget.NewEntry()
	baseDirectoryEntry.SetText(dir)
	// Create the dog entry
	dogEntry := widget.NewEntry()
	// Create progress bar
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Stop()
	progressBar.Hide()
	// Create download button
	downloadWindow := mainApp.NewWindow("Download Progress")
	downloadEntry := widget.NewMultiLineEntry()
	downloadCloseButton := widget.NewButton("Close", func() {
		downloadEntry.SetText("")
		downloadWindow.Hide()
	})
	downloadWindow.SetContent(widget.NewVBox(downloadEntry, downloadCloseButton))
	downloadButton := widget.NewButton("Download", func() {
		progressChannel := make(chan string, 10)
		// Create directory where the dog info will go
		if err := io.MakeDir(baseDirectoryEntry.Text); err != nil {
			errorEntry := widget.NewEntry()
			errorEntry.SetText(fmt.Sprintf("%+v", err))
			errorWindow.SetContent(widget.NewVBox(errorEntry, errorCloseButton))
			errorWindow.Show()
		}
		progressBar.Start()
		progressBar.Show()
		downloadWindow.Show()
		go func() {
			if err := RunDogDownloads(dogEntry.Text, baseDirectoryEntry.Text, progressChannel, errorChannel); err != nil {
				errorEntry := widget.NewEntry()
				errorEntry.SetText(fmt.Sprintf("%+v", err))
				errorWindow.SetContent(widget.NewVBox(errorEntry, errorCloseButton))
				errorWindow.Show()
			}
		}()
		downloadEntry.SetText("Downloading...\n")
		for progress := range progressChannel {
			downloadEntry.SetText(downloadEntry.Text + progress + "\n")
		}
		progressBar.Stop()
		progressBar.Hide()
	})
	// Create foster window
	boardingDogsWindow := mainApp.NewWindow("Boarding Dogs")
	boardingCloseButton := widget.NewButton("Close", func() {
		boardingDogsWindow.Hide()
	})
	// Set the window content
	mainWindow.SetContent(widget.NewVBox(
		// Lookup fosters group
		widget.NewGroup("Lookup", widget.NewButton("Get Boarding List", func() {
			progressBar.Start()
			progressBar.Show()
			downloadButton.Disable()
			baseDirectoryEntry.Disable()
			dogEntry.Disable()
			fosters, err := RunGetFosters(errorChannel)
			if err != nil {
				errorEntry := widget.NewEntry()
				errorEntry.SetText(fmt.Sprintf("%+v", err))
				errorWindow.SetContent(widget.NewVBox(errorEntry, errorCloseButton))
				errorWindow.Show()
			}
			dogs := widget.NewMultiLineEntry()
			dogs.SetText(joinDogs(fosters))
			boardingDogsWindow.SetContent(widget.NewVBox(dogs, boardingCloseButton))
			progressBar.Stop()
			progressBar.Hide()
			downloadButton.Enable()
			baseDirectoryEntry.Enable()
			dogEntry.Enable()
			boardingDogsWindow.Show()
		})),
		// Download dogs group
		widget.NewGroup("Dog Download", widget.NewForm(&widget.FormItem{
			Text:   "Output Directory:",
			Widget: baseDirectoryEntry,
		}, &widget.FormItem{
			Text:   "Dogs (comma separated):",
			Widget: dogEntry,
		}), downloadButton),
		progressBar,
		// Quit
		quitButton,
	))
	// Run it
	mainWindow.ShowAndRun()
	close(errorChannel)
}

func joinDogs(fosters []string) string {
	sort.Sort(sort.StringSlice(fosters))
	return insertNewLine(strings.Join(fosters, ","), 100)
}

func insertNewLine(input string, index int) string {
	var buffer bytes.Buffer
	var currentPos = index - 1
	var lastIndex = len(input) - 1
	for i, r := range input {
		buffer.WriteRune(r)
		if i%index == currentPos && i != lastIndex {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}
