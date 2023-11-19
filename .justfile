build: buildl buildw buildm
    echo "build all"

buildl:
    echo "build linux"
    GOOS=linux GOARCH=amd64 go build -o bin/dumper_l main.go

buildw:
    echo "build windows"
    GOOS=windows GOARCH=amd64 go build -o bin/dumper.exe main.go

buildm:
    echo "build mac"
    GOOS=darwin GOARCH=arm64 go build -o bin/dumper_m main.go
