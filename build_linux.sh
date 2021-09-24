GOOS=linux GOARCH=amd64 go build -o bin/viewcounter-amd64-linux server.go
GOOS=linux GOARCH=386 go build -o bin/viewcounter-386-linux server.go
