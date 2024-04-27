package gin

import (
	"github.com/fzdwx/authz"
	"github.com/gin-gonic/gin"
)

const (
	TokenKey   string = "authz_token"
	SessionKey string = "authz_session"
)

func GetToken(c *gin.Context) string {
	return c.GetString(TokenKey)
}

func GetSession[ID authz.IdType](c *gin.Context) *authz.Session[ID] {
	value, exists := c.Get(SessionKey)
	if !exists {
		return nil
	}
	return value.(*authz.Session[ID])
}

func Middleware[ID authz.IdType](client authz.Client[ID]) *middleware[ID] {
	return &middleware[ID]{
		client:    client,
		tokenKey:  defaultKey,
		useCookie: true,
		whiteMap:  make(map[string]bool)}
}

type middleware[ID authz.IdType] struct {
	client    authz.Client[ID]
	tokenKey  string
	whiteMap  map[string]bool
	useCookie bool
}

var defaultKey = "Authorization"

func (m *middleware[ID]) TokenKey(tokenKey string) *middleware[ID] {
	m.tokenKey = tokenKey
	return m
}

func (m *middleware[ID]) WhiteList(whiteList []string) *middleware[ID] {
	for _, path := range whiteList {
		m.whiteMap[path] = true
	}
	return m
}

func (m *middleware[ID]) UseCookie(useCookie bool) *middleware[ID] {
	m.useCookie = useCookie
	return m
}

func (m *middleware[ID]) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if m.useCookie {
			if cookie, err := c.Cookie(TokenKey); err == nil {
				token = cookie
			}
		}
		if token == "" {
			tokenVal := c.GetHeader(m.tokenKey)
			if tokenVal != "" {
				if m.tokenKey == defaultKey {
					token = tokenVal[len("Bearer "):]
				} else {
					token = tokenVal
				}
			}
		}

		if token == "" {
			if _, ok := m.whiteMap[c.Request.URL.Path]; ok {
				c.Next()
				return
			}
		}

		ctx := authz.SetToken(c, token)
		session, err := m.client.GetSession(ctx)
		if err != nil {
			if m.useCookie {
				c.SetCookie(TokenKey, "", -1, "/", "", false, true)
			}
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set(TokenKey, token)
		c.Set(SessionKey, session)
		if m.useCookie {
			c.SetCookie(TokenKey, token, 3600, "/", "", false, true)
		}
		c.Next()
	}
}
