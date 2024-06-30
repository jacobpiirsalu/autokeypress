package main

import (
    "strconv"
    "sync"
    "time"
	"strings"

    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/go-vgo/robotgo"
)


func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("AutoKeyPress")

    intervalEntry := widget.NewEntry()
    intervalEntry.SetPlaceHolder("Interval in milliseconds")

    var stopChan chan struct{}
    var wg sync.WaitGroup

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
                    robotgo.KeyTap("space")
                case <-stopChan:
                    return
                }
            }
        }()
    })

    stopButton := widget.NewButton("Stop", func() {
        if stopChan != nil {
            close(stopChan)
            wg.Wait()
            stopChan = nil
        }
    })

    content := container.NewVBox(
        widget.NewLabel("Enter the interval between spacebar presses:"),
        intervalEntry,
        startButton,
        stopButton,
    )

    myWindow.SetContent(content)
    myWindow.ShowAndRun()
}
 