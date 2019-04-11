webserver
=========

Generates dynamic children URLs from given URL (To help test crawling for example)
## Usage

```
go run webserver.go
```

Server starts @ http://localhost:3030

## Generated Dynamic URLs

`http://localhost:3030/this-is-a-path`

will display this list of children URLs:

```
http://localhost:3030/this-is-a-path/this-is-a-pat
http://localhost:3030/this-is-a-path/this-is-a-pa
http://localhost:3030/this-is-a-path/this-is-a-p
http://localhost:3030/this-is-a-path/this-is-a-
http://localhost:3030/this-is-a-path/this-is-a
http://localhost:3030/this-is-a-path/this-is-
http://localhost:3030/this-is-a-path/this-is
http://localhost:3030/this-is-a-path/this-i
http://localhost:3030/this-is-a-path/this-
http://localhost:3030/this-is-a-path/this
http://localhost:3030/this-is-a-path/thi
http://localhost:3030/this-is-a-path/th
http://localhost:3030/this-is-a-path/t
```

`http://localhost:3030/this-is-a-path/this`

will display children URLs

```
http://localhost:3030/this-is-a-path/this/thi
http://localhost:3030/this-is-a-path/this/th
http://localhost:3030/this-is-a-path/this/t
```

**etc ...**