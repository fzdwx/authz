# authz

A simple authorization library for Go.

```shell
go get github.com/fzdwx/authz@main
```

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

	var session, _ = c.GetSession(ctx, token)
	fmt.Println(session.ID)
}

```
