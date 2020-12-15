package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/pascallin/go-x-app/internal/ui"
)

var (
	topWindow fyne.Window
	APP_NAME = "pascal-x-app"
	VERSION = "v1"
	)

func main() {
	if runtime.GOOS == "windows" {
		err := os.Setenv("FYNE_FONT", path.Join("C://Windows/Fonts", "msyh.ttc"))
		if err != nil {
			panic(err)
		}
	} else if runtime.GOOS == "darwin" {
		err := os.Setenv("FYNE_FONT", path.Join("/System/Library/Fonts", "PingFang.ttf"))
		if err != nil {
			panic(err)
		}
	}

	a := app.NewWithID(fmt.Sprintf("%s (%s)", APP_NAME, VERSION))
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow(APP_NAME)
	topWindow = w

	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t ui.Screen) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(ui.MakeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(ui.MakeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}