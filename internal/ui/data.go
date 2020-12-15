package ui

import (
	"fyne.io/fyne"
)

// Tutorial defines the data structure for a tutorial
type Screen struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Screens = map[string]Screen{
		"welcome": {"欢迎页", "", welcomeScreen},
		"xlsx": {"Excel(xlsx)文件工具",
			"自定义列名，往行末尾插入在图片文件夹中找到的图片文件\n支持（png, jpg, jpeg, gif格式）",
			xlsxScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	ScreensIndex = map[string][]string{
		"":            {"welcome", "xlsx"},
	}
)
