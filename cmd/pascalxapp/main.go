package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/pascallin/go-x-app/internal/screens"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func main() {
	fmt.Println(runtime.GOOS)
	if runtime.GOOS == "windows" {
		os.Setenv("FYNE_FONT", path.Join("C://Windows/Fonts", "msyh.ttc"))
	}
	if runtime.GOOS == "darwin" {
		os.Setenv("FYNE_FONT", path.Join("/System/Library/Fonts", "PingFang.ttc"))
	}

	a := app.NewWithID("pascal-x-app-v1")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("pascal-x-app")
	topWindow = w

	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t screens.Screen) {
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
		w.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(720, 480))
	w.ShowAndRun()
}

func makeNav(setTutorial func(tutorial screens.Screen), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return screens.ScreensIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := screens.ScreensIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := screens.Screens[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := screens.Screens[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		widget.NewButton("暗黑主题", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("白色主题", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
