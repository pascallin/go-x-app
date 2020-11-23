package internal

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type InsertOptions struct {
	ExcelFilePath string
	SheetName string
	ImageDir string
	KeyColumn string
}

func Insert(ops *InsertOptions) error {
	xlsx, err := excelize.OpenFile(ops.ExcelFilePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//imagesRecords := scanImages(ops.ImageDir)

	rows, err := xlsx.GetRows(ops.SheetName)
	for i, _ := range rows {
		cell, err := xlsx.GetCellValue(ops.SheetName, ops.KeyColumn + strconv.Itoa(i+1))
		if err != nil {
			return err
		}
		imagePath := findImageByKey(ops.ImageDir, cell)
		// Insert a picture.
		// NOTE: autofit not working
		err = xlsx.AddPicture(ops.SheetName, "C2", imagePath , `{"autofit": true}`)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	//err = xlsx.SetRowHeight(ops.SheetName, 2, float64(imagesRecords[0].Height))
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//err = xlsx.SetColWidth(ops.SheetName, "B", "C", float64(imagesRecords[0].Width))
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Save the xlsx file with the origin path.
	err = xlsx.Save()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type AppImage struct {
	Key string
	Width int
	Height int
	Path string
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

func findImageByKey(dir string, key string) string {
	suffixs := []string{".jpg", "jpeg", ".png", "gif"}
	for _, item := range(suffixs) {
		fmt.Println(path.Join(dir, key, item))
		isExist, _ := PathExists(path.Join(dir, key, item))
		if isExist {
			return path.Join(dir, key, item)
		}
	}
	return ""
}

func scanImages(dirToScan string) []*AppImage {
	var result []*AppImage
	files, _ := ioutil.ReadDir(dirToScan)
	for _, imgFile := range files {
		fileExt := path.Ext(path.Join(dirToScan, imgFile.Name()))
		if fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" && fileExt != ".jpg" {
			continue
		}
		if reader, err := os.Open(filepath.Join(dirToScan, imgFile.Name())); err == nil {
			defer reader.Close()
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}
			fmt.Printf("%s %d %d\n", imgFile.Name(), im.Width, im.Height)
			result = append(result, &AppImage{
				Width: im.Width,
				Height: im.Height,
				Path: path.Join(dirToScan, imgFile.Name()),
				Key: strings.TrimSuffix(imgFile.Name(), fileExt),
			})
		} else {
			fmt.Println("Impossible to open the file:", err)
		}
	}
	return result
}