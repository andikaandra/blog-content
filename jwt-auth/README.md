Golang auth use JWT and redis 

requirements:
- github.com/gin-gonic/gin
- github.com/go-redis/redis
- github.com/dgrijalva/jwt-go
- github.com/google/uuid

goal:
- create public api
- create api for authenticated user
- create api for its own user (authenticate + authorize)

---

#### Main Application
```
func main() {
	r := gin.Default()

	r.POST("/login", loginRoute)
	r.GET("/users", allUserRoute)

	r.Use(tokenAuthMiddleware())
	{
		r.GET("/books", allBookRoute)
		r.GET("/user/:userID", findUserRoute)
	}

	r.Run()
}
```

public api routes
`/login`
`/users `

api for authenticated user
`/books`

api for its own user
`/user/:userID`

TODO..
