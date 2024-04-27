# authz with gin

## Install

```shell
go get github.com/fzdwx/authz/contrib/middleware/gin
```

## Run the demo

```shell
go run main.go
```

1. access `http://localhost:8080/` you will see login button
2. click login button, you will see the user info get login success mage
3. access `http://localhost:8080/` you will get the user info


```go
package main

import (
	"github.com/fzdwx/authz"
	ginx "github.com/fzdwx/authz/contrib/middleware/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	authzClient := authz.NewClient[string](
		authz.NewMemoryStore(),
		authz.DefaultPermissionSupplier[string]{},
	)

	r.Use(
		ginx.Middleware(authzClient).
			WhiteList([]string{"/", "/auth"}).
			Build(),
	)

	r.GET("/auth", func(c *gin.Context) {
		if token, err := authzClient.Login(c, &authz.LoginOption[string]{
			ID: "Hello",
			Metadata: map[string]string{
				"role": "admin",
			},
		}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		} else {
			c.SetCookie(ginx.TokenKey, token, 3600, "/", "localhost", false, true)
		}
		c.JSON(200, gin.H{"message": "login success"})
	})

	r.GET("/", func(context *gin.Context) {
		session := ginx.GetSession[string](context)

		if session == nil {
			context.Header("Content-Type", "text/html")
			context.String(200, `<html>
<a href="/auth">Login</a>
</html>`)
			return
		}

		context.JSON(200, gin.H{"message": session})
	})

	r.Run(":8080")
}
```

