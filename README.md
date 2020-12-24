# go-x-app

Some functionalities utils Chinese language Application for cross-platform.

For personal usage, but feel free to view and comment, even reuse it.

## functionalities

### Insert image to excel file

![scene snap](https://github.com/pascallin/go-x-app/blob/main/images/excel_method.png?raw=true)

## development

```shell script
go run ./cmd/pascalxapp/main.go
```

## build

```shell script
# win10
GOOS=windows go build -x -v -o windows/pascalxapp-v1.10.exe ./cmd/pascalxapp/main.go

# MacOS
GOOS=darwin go build -x -v -o macos/pascalxapp ./cmd/pascalxapp/main.go
```

## NOTE

### need to fix package

jump to `github.com/360EntSecGroup-Skylar/excelize/v2` package `picture.go` file, search `drawingResize`, and comment rows `630-633`, target to `autofix` just base on row width

```go
//if float64(cellHeight) < height {
//	asp := float64(cellHeight) / height
//	height, width = float64(cellHeight), width*asp
//}
```

jump to `github.com/360EntSecGroup-Skylar/excelize/v2` package `col.go` file and comment rows `608-611`, target to `autofix` just base on row width

```go
// Subtract the underlying cell heights to find end cell of the object.
// for height >= f.getRowHeight(sheet, rowEnd) {
// 	height -= f.getRowHeight(sheet, rowEnd)
// 	rowEnd++
// }
```