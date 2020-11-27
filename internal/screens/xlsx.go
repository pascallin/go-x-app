package screens

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

	infProgress *widget.ProgressBarInfinite
	resultLabel *widget.Label
	paramsCard *widget.Card
)

func xlsxScreen(win fyne.Window) fyne.CanvasObject {
	infProgress = widget.NewProgressBarInfinite()
	resultLabel = widget.NewLabel("准备中")
	// NOTE: to soon for execution
	infProgress.Hide()
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

				//children, err := list.List()
				//if err != nil {
				//	dialog.ShowError(err, win)
				//	return
				//}
				//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
				//dialog.ShowInformation("Folder Open", out, win)
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
		infProgress, // NOTE: not showing
		paramsCard,
		resultLabel,
	))
}

func makeForm() fyne.CanvasObject {
	sheetNameFormItem := widget.NewEntry()
	sheetNameFormItem.SetPlaceHolder("请输入表格名称")

	kayColumnName := widget.NewEntry()
	kayColumnName.SetPlaceHolder("请输入excel列名")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Sheet名称", Widget: sheetNameFormItem},
			{Text: "需要匹配的列名", Widget: kayColumnName},
		},
		SubmitText: "执行",
	}

	form.OnSubmit = func() {
		fmt.Println("Form submitted")

		resultLabel.Text = "进行中...."
		sheetNameFormItem.Disable()
		kayColumnName.Disable()
		//infProgress.Start()

		err := internal.Insert(&internal.InsertOptions{
			ExcelFile: xlsxPath,
			ImageDir: imgDirPath,
			KeyColumn:kayColumnName.Text,
			SheetName:sheetNameFormItem.Text},
		)
		if err != nil {
			fyne.LogError("Execution Faild", err)
		}

		if err != nil {
			fyne.LogError("Failed close reader", err)
		}

		resultLabel.Text = "已完成"
		sheetNameFormItem.Enable()
		kayColumnName.Enable()
		//infProgress.Stop()
	}

	return form
}