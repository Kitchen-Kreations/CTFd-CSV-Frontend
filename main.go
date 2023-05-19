package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"math/rand"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	TYPE         string = "standard" // challenge type
	STATE        string = "visible"  // challenge state
	MAX_ATTEMPTS string = "0"        // challenge max attemptes (0 = unlimited)
)

var (
	csv_path         string   = ""
	category_options []string = []string{"Assessment", "Practice Range", "Main Range"}
	headers          []string = []string{"name", "description", "category", "value", "type", "state", "max_attempts", "flags"}

	charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// n is the length of random string we want to generate
func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		// randomly select 1 character from given charset
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	// initialize fyne app
	App := app.New()

	// starting window
	Start_Window := App.NewWindow("Select CSV")
	Start_Window.Resize(fyne.NewSize(600, 600))

	// config app and main window
	Main_Window := App.NewWindow("CTFd CSV Creator")
	Main_Window.Resize(fyne.NewSize(600, 600))

	// starting window widgets
	csv_new_file_entry := widget.NewEntry()
	csv_open_new := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		csv_new_file_entry.SetText(uc.URI().Path())
	}, Start_Window)
	csv_open_new_button := widget.NewButton("New", func() {
		csv_open_new.Show()
		csv_open_new.Resize(fyne.NewSize(600, 600))
	})

	csv_open_file_entry := widget.NewEntry()
	csv_open_old := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		csv_open_file_entry.SetText(uc.URI().Path())
	}, Start_Window)
	csv_open_old_button := widget.NewButton("Open", func() {
		csv_open_old.Show()
		csv_open_old.Resize(fyne.NewSize(600, 600))
	})

	// starting window form
	start_window_form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "CSV New File: ", Widget: csv_new_file_entry},
			{Text: "", Widget: csv_open_new_button},
			{Text: "CSV Open File: ", Widget: csv_open_file_entry},
			{Text: "", Widget: csv_open_old_button},
		},
		OnSubmit: func() {
			// error handling
			if csv_new_file_entry.Text != "" && csv_open_file_entry.Text != "" {
				log.Fatal("only one entry can have data")
			}

			if csv_new_file_entry.Text != "" {
				// change csv path
				csv_path = csv_new_file_entry.Text

				// truncate or create file & write csv headers
				csv_file, err := os.Create(csv_path)
				if err != nil {
					log.Fatal("failed to truncate and/or create csv file\ncsv_path: " + csv_path + "\nError: " + err.Error())
				}

				csv_writer := csv.NewWriter(csv_file)
				defer csv_writer.Flush()

				err = csv_writer.Write(headers)
				if err != nil {
					log.Fatal("failed to write headers to csv file\nError: " + err.Error())
				}
			} else if csv_open_file_entry.Text != "" {
				csv_path = csv_open_file_entry.Text

				// check if file contains the correct headers
				read_file, err := os.Open(csv_path)
				if err != nil {
					log.Fatal(err)
				}
				file_scanner := bufio.NewScanner(read_file)

				file_scanner.Split(bufio.ScanLines)

				file_scanner.Scan()
				if file_scanner.Text() != "name,description,category,value,type,state,max_attempts,flags" {
					log.Fatal("invalid csv file")
				}
			} else {
				log.Fatal("one option must be selected")
			}

			// cycle to new window
			Start_Window.Hide()
			Main_Window.Show()
		},
		OnCancel: func() {
			App.Quit()
		},
	}

	// start window content
	Start_Window.SetContent(start_window_form)

	// main window widgets
	challenge_name_entry := widget.NewEntry()
	challenge_description_entry := widget.NewMultiLineEntry()
	challenge_category_dropdown := widget.NewSelect(category_options, func(s string) {})
	challenge_value_entry := widget.NewEntry()
	challenge_flag_entry := widget.NewEntry()

	generate_flag_button := widget.NewButton("Generate Flag", func() {
		challenge_flag_entry.SetText(randStr(32))
	})

	save_check := widget.NewCheck("Save Entries", func(b bool) {})

	success_message_label := widget.NewLabel("")

	// main window form
	main_window_form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Challenge Name: ", Widget: challenge_name_entry},
			{Text: "Challenge Description: ", Widget: challenge_description_entry},
			{Text: "Challenge Category: ", Widget: challenge_category_dropdown},
			{Text: "Challenge Value: ", Widget: challenge_value_entry},
			{Text: "Challenge Flag: ", Widget: challenge_flag_entry},
			{Text: "", Widget: generate_flag_button},
			{Text: "", Widget: save_check},
			{Text: "", Widget: success_message_label},
		},
		OnSubmit: func() {
			name := challenge_name_entry.Text
			description := challenge_description_entry.Text
			category := challenge_category_dropdown.Selected
			value := challenge_value_entry.Text
			flag := challenge_flag_entry.Text

			data := []string{name, description, category, value, TYPE, STATE, MAX_ATTEMPTS, flag}

			// open csv file
			csv_file, err := os.OpenFile(csv_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal("failed to truncate and/or create csv file\ncsv_path: " + csv_path + "\nError: " + err.Error())
			}

			// write to csv file
			csv_writer := csv.NewWriter(csv_file)
			defer csv_writer.Flush()

			err = csv_writer.Write(data)
			if err != nil {
				log.Fatal("failed to write headers to csv file\nError: " + err.Error())
			}

			// write success message
			success_message_label.SetText("Successfully Wrote Challenge")

			// clear fields
			if !save_check.Checked {
				challenge_name_entry.SetText("")
				challenge_description_entry.SetText("")
				challenge_value_entry.SetText("")
				challenge_flag_entry.SetText("")
			}
		},
		OnCancel: func() {
			App.Quit()
		},
	}

	// set main window content
	Main_Window.SetContent(main_window_form)

	// start the app
	Start_Window.ShowAndRun()
}
