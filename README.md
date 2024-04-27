# authz

A simple authorization library for Go.

```shell
go get github.com/fzdwx/authz@main
```

## Supported frameworks

- [x] [Gin](./examples/gin)

## Usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/fzdwx/authz"
)

func foo() {
	ctx := context.Background()
	c := atuhz.NewClient[string](authz.NewMemoryStore(), DefaultPermissionSupplier[string]{})

	var token, _ = c.Login(ctx, &authz.LoginOption[string]{
		ID: "1",
	})
	ctx = SetToken(ctx, token)
	var session, _ = c.GetSession(ctx)
	fmt.Println(session.ID)
}

```
