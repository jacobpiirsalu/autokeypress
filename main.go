package main

import (
	"strconv"
	"strings"
	"sync"
	"time"
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

// TODO: make the clicking togglable by JKL
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("AutoKeyPress")

	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder("Interval in milliseconds")

	var stopChan chan struct{}
	var wg sync.WaitGroup
	var selectedOption string

	startButton := widget.NewButton("Start", func() {
		intervalStr := intervalEntry.Text
		interval, err := strconv.Atoi(strings.TrimSpace(intervalStr))
		if err != nil {
			return
		}

		stopChan = make(chan struct{})
		wg.Add(1)

		go func() {
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if selectedOption == "spacebar" {
						robotgo.KeyTap("space")
					} else if selectedOption == "left click" {
						robotgo.Click()
					}

				case <-stopChan:
					return
				}
			}
		}()
       
		 // Goroutine to capture keyboard input
		 go func() {
            fmt.Println("--- Please hold down the keys JKL to exit ---")
            hook.Register(hook.KeyDown,[]string{"j","k","l"}, func(e hook.Event) {
                fmt.Println("stopping...")
                close(stopChan)
                hook.End()
            })

            s := hook.Start()
            <-hook.Process(s)
        }()
	})

	stopButton := widget.NewButton("Stop", func() {
		if stopChan != nil {
			close(stopChan)
			wg.Wait()
			stopChan = nil
		}
	})

	selection := widget.NewSelect([]string{"spacebar", "left click"}, func(option string) {
		selectedOption = option
	})

	content := container.NewVBox(
		widget.NewLabel("Enter the interval between key presses:"),
		intervalEntry,
		selection,
		startButton,
		stopButton,
		widget.NewLabel("Hold down the keys JKL at any time to stop the program"),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}