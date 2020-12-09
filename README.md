# go-x-app

## development

```shell script
go run ./cmd/pascalxapp/main.go
```

## compile

only support for now because of Chinese Font location using absolute path.

```shell script
GOOS=windows go build -x -v -o windows/pascalxapp.exe ./cmd/pascalxapp/main.go
```