package authz

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	c := NewClient[string](NewMemoryStore(), DefaultPermissionSupplier[string]{})
	var token, err = c.Login(ctx, &LoginOption[string]{
		ID: "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetSession(ctx)
	assert.Equal(t, ErrNoToken, err)

	ctx = SetToken(ctx, token)
	var s, err2 = c.GetSession(ctx)
	if err2 != nil {
		t.Fatal(err2)
	}
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, 1, len(s.Tokens))
	assert.Equal(t, token, s.Tokens[0].Value)

	if err := c.Set(ctx, "hello", "world"); err != nil {
		t.Fatal(err)
	}
	session, err := c.GetSession(ctx)
	if err != nil {
		t.Fatal()
	}
	assert.Equal(t, "world", session.Metadata["hello"])
}
