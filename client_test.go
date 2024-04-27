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

	ctx = SetToken(ctx, token)
	var s, err2 = c.GetSession(ctx)
	if err2 != nil {
		t.Fatal(err2)
	}
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, 1, len(s.Tokens))
	assert.Equal(t, token, s.Tokens[0].Value)
}
