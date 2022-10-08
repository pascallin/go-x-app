# go-x-app

Some functionalities utils Chinese language Application for cross-platform.

For personal usage, but feel free to view and comment, even reuse it.

## functionalities

### Insert image to excel file

![scene snap](https://github.com/pascallin/go-x-app/blob/main/images/excel_method.png?raw=true)

## development

Prerequisites: [fyne Getting Started](https://developer.fyne.io/started/#prerequisites)

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
