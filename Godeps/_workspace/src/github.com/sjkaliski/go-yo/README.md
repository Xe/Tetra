go-yo
=====

Golang client for the Yo API

***Documentation***: http://godoc.org/github.com/sjkaliski/go-yo

***Build Status***: [![Build Status](https://travis-ci.org/sjkaliski/go-yo.png)](https://travis-ci.org/sjkaliski/go-yo)

***API Info***: https://medium.com/@YoAppStatus/yo-developers-api-e7f2f0ec5c3c

## Usage

With token in hand, to create a new Yo client simply

```go
package main

import (
  "github.com/sjkaliski/go-yo"
)

func main() {
  client := yo.NewClient("my_token")
}
```

To send a message to all users who subscribe to you

```go
err := client.YoAll()
```

To send a message to a specific user

```
err := client.YoUser("some_user")
```

## Tests

To run tests

`$ go test ./...`
