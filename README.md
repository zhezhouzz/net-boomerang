# net-boomerang

#### get this repo

```sh
go get github.com/zhezhouzz/net-boomerang
```

#### build

+ server

```sh
cd src
go run server.go
```

+ client

```sh
cd src
go run client.go
```

+ client would sent a file(this path in the header of client.go) to server, then server would save this to a file(path in the header of server.go).
+ `section transfer`. This feature need to test and review.