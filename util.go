package authz

import "math/rand"

var tokenChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(i int) string {
	b := make([]rune, i)
	for i := range b {
		b[i] = tokenChars[rand.Intn(len(tokenChars))]
	}
	return string(b)
}
