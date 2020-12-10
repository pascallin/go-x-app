package internal

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type InsertOptions struct {
	ExcelFile string
	SheetName string
	ImageDir string
	KeyColumn string
	SelectColumn string
	IsCoverFile bool
}

type ExcelProgress struct {
	ProgressChannel chan float64
	NewFilePath string
}

func NewExcelProgress() *ExcelProgress {
	return &ExcelProgress{}
}

func (ep *ExcelProgress) InsertImage(ops *InsertOptions) error {
	if !ops.IsCoverFile {
		rand.Seed(time.Now().Unix())
	}

	xlsx, err := excelize.OpenFile(ops.ExcelFile)
	if err != nil {
		return err
	}

	rows, err := xlsx.GetRows(ops.SheetName)
	if err != nil {
		return err
	}

	headers := rows[0]
	// get key axis
	keyIndex := SliceIndex(len(headers), func(i int) bool { return headers[i] == ops.KeyColumn })
	if keyIndex == -1 {
		return errors.New("could not found key column")
	}
	keyAxis := toChar(keyIndex)
	// get image axis
	imageIndex := SliceIndex(len(headers), func(i int) bool { return headers[i] == ops.SelectColumn })
	if imageIndex == -1 {
		return errors.New("could not found image column")
	}
	imageAxis := toChar(imageIndex)

	//imageFormat := `{
	//	"autofit": true,
	//	"locked": true,
	//	"print_obj": true,
	//	"lock_aspect_ratio": true
	//}`
	imageFormat := `{
		"autofit": true
	}`

	ep.ProgressChannel = make(chan float64)
	//defer close(ep.ProgressChannel)
	for row, _ := range rows {
		cell, err := xlsx.GetCellValue(ops.SheetName, keyAxis + strconv.Itoa(row+1))
		if err != nil {
			return err
		}
		imagePath := findImageByKey(ops.ImageDir, cell)
		// Insert a picture.
		if imagePath != "" {
			// Set cell width and height
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

			//err = xlsx.AddPicture(ops.SheetName, imageAxis + strconv.Itoa(row+1), imagePath, imageFormat)
			//if err != nil {
			//	return err
			//}

			//imageColWith, _ := xlsx.GetColWidth(ops.SheetName, imageAxis)
			//imageWidth, _, err := getImageWidthAndHeight(imagePath)
			//ratio := imageColWith / float64(imageWidth)
			//fmt.Println(imageAxis, imageColWith, imageWidth, ratio)
			//imageFormat := fmt.Sprintf(`{"x_scale":%f,"y_scale":%f}`, ratio, ratio)

			err = xlsx.AddPicture(ops.SheetName, imageAxis + strconv.Itoa(row+1), imagePath, imageFormat)
			if err != nil {
				return err
			}
		}
		go func() { ep.ProgressChannel <- float64(row + 1) / float64(len(rows)) }()
	}

	// Save the xlsx file with the origin path.
	if ops.IsCoverFile {
		err = xlsx.Save()
	} else {
		newFile := fmt.Sprintf("%s_%s.xlsx", strings.TrimSuffix(ops.ExcelFile, ".xlsx"), strconv.Itoa(rand.Intn(10000)))
		ep.NewFilePath = newFile
		err = xlsx.SaveAs(newFile)
	}
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

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
func toChar(i int) string {
	return arr[i]
}
func findImageByKey(dir string, key string) string {
	suffixs := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, suffix := range(suffixs) {
		isExist, _ := PathExists(path.Join(dir, key) + suffix)
		if isExist {
			fmt.Println("found match file: ", path.Join(dir, key), suffix)
			return path.Join(dir, key) + suffix
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

func getImageWidthAndHeight(picture string) (int, int, error) {
	file, _ := ioutil.ReadFile(picture)
	img, _, err := image.DecodeConfig(bytes.NewReader(file))
	if err != nil {
		return 0, 0, err
	}
	return img.Width,img.Height, nil
}