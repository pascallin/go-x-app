package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"

	"github.com/pascallin/go-x-app/internal"
)

var (
	imgDirPath string
	xlsxPath string

	progressBar *widget.ProgressBar
	resultLabel *widget.Label
	paramsCard *widget.Card
)

func xlsxScreen(win fyne.Window) fyne.CanvasObject {
	progressBar = widget.NewProgressBar()
	progressBar.SetValue(0)
	resultLabel = widget.NewLabel("准备中")
	paramsCard := widget.NewCard("获取到的参数", "", widget.NewLabel(fmt.Sprintf("图片文件夹地址: %s\nExcel文件地址: %s\n", imgDirPath, xlsxPath)))

	return container.NewVScroll(container.NewVBox(
		widget.NewButton("选择图片文件夹", func() {
			dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if list == nil {
					return
				}
				imgDirPath = strings.TrimPrefix(list.String(), "file://")
				paramsCard.SetContent(widget.NewLabel(fmt.Sprintf("pictures folder: %s\nxlsx file path: %s\n", imgDirPath, xlsxPath)))
			}, win)
		}),
		widget.NewButton("选择Excel文件", func() {
			fd := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
				if err == nil && f == nil {
					return
				}
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				xlsxPath = strings.TrimPrefix(f.URI().String(), "file://")
				paramsCard.SetContent(widget.NewLabel(fmt.Sprintf("pictures folder: %s\nxlsx file path: %s\n", imgDirPath, xlsxPath)))
			}, win)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
			fd.Show()
		}),
		makeForm(),
		paramsCard,
		resultLabel,
		progressBar,
	))
}

func makeForm() fyne.CanvasObject {
	sheetNameFormItem := widget.NewEntry()
	sheetNameFormItem.SetPlaceHolder("请输入表格名称")
	sheetNameFormItem.SetText("Sheet1")

	keyColumnName := widget.NewEntry()
	keyColumnName.SetPlaceHolder("请输入匹配列头名称")
	keyColumnName.SetText("款号")

	selectColumnName := widget.NewEntry()
	selectColumnName.SetPlaceHolder("请输入需要插入图片的列头名称")
	selectColumnName.SetText("图片")

	var (
		isCoverFile = false
		coverFileText = "覆盖原文件"
		saveAsText = "另存文件（文件名_XXXX）"
	)
	radio := widget.NewRadioGroup([]string{coverFileText, saveAsText}, func(s string) {
		if s == coverFileText {
			isCoverFile = true
		} else {
			isCoverFile = false
		}
	})
	radio.SetSelected(saveAsText)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Sheet名称", Widget: sheetNameFormItem},
			{Text: "匹配关键字列的列头名称", Widget: keyColumnName},
			{Text: "匹配插入列的列头名称", Widget: selectColumnName},
			{Text: "文件保存方式", Widget: radio},
		},
		SubmitText: "执行",
	}

	form.OnSubmit = func() {
		fmt.Println("Form submitted")

		resultLabel.Text = "进行中...."
		sheetNameFormItem.Disable()
		keyColumnName.Disable()
		selectColumnName.Disable()
		radio.Disable()
		progressBar.Show()
		progressBar.SetValue(0)

		ep := internal.NewExcelProgress()
		err := ep.InsertImage(&internal.InsertOptions{
			ExcelFile: xlsxPath,
			ImageDir: imgDirPath,
			KeyColumn: keyColumnName.Text,
			SelectColumn: selectColumnName.Text,
			SheetName: sheetNameFormItem.Text,
			IsCoverFile: isCoverFile,
		})

		go listenProgress(ep, progressBar)

		if err != nil {
			fyne.LogError("Execution Faild", err)
		}

		if err != nil {
			fyne.LogError("Failed close reader", err)
		}

		resultLabel.Text = "已完成，保存文件路径：" + ep.NewFilePath
		sheetNameFormItem.Enable()
		keyColumnName.Enable()
		selectColumnName.Enable()
		radio.Enable()
	}

	return form
}

func listenProgress(ep *internal.ExcelProgress, bar *widget.ProgressBar) {
	for {
		select {
			case progress := <-ep.ProgressChannel:
			bar.SetValue(progress)
		}
	}
}