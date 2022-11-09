# cmd/agent

## Build info

Setup the -ldflags option for build info

Build:
```sh
go build -v -ldflags="-X 'main.buildVersion=0.1.0' \
 -X 'main.buildDate=$(date)' \
 -X 'main.buildCommit=test'" -o ./agent \
 cmd/agent/main.go
```
Run:
```sh
go run -v -ldflags="-X 'main.buildVersion=0.2.0' \
 -X 'main.buildDate=$(date)' \
 -X 'main.buildCommit=test'" \
cmd/agent/main.go
```

