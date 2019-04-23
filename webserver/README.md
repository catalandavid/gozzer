webserver
=========

Generates dynamic children URLs from given URL (To help test crawling for example)
## Usage

```
go run webserver.go
```

Server starts @ http://localhost:3030

### Generated Dynamic URLs

`http://localhost:3030/this-is-a-path`

will display this list of children URLs:

```
/this-is-a-path/this-is-a-pat
/this-is-a-path/this-is-a-pa
/this-is-a-path/this-is-a-p
/this-is-a-path/this-is-a-
/this-is-a-path/this-is-a
/this-is-a-path/this-is-
/this-is-a-path/this-is
/this-is-a-path/this-i
/this-is-a-path/this-
/this-is-a-path/this
/this-is-a-path/thi
/this-is-a-path/th
/this-is-a-path/t
```

`http://localhost:3030/this-is-a-path/this`

will display children URLs

```
/this-is-a-path/this/thi
/this-is-a-path/this/th
/this-is-a-path/this/t
```

**etc ...**

## Build

### For Mac

```
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a .
```
### For Linux

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a .
```

### For Windows

```
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a .
```
