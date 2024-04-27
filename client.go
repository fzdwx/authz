package authz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samber/lo"
)

type IdType interface {
	int64 | int | string
}

type Client[ID IdType] interface {
	// Login creates a new session for the user with the given ID and returns a token
	Login(ctx context.Context, opt *LoginOption[ID]) (string, error)
	// GetSession returns the session associated with the given token
	GetSession(ctx context.Context) (*Session[ID], error)
	// Set sets the value of the given key in the session metadata
	Set(ctx context.Context, key, value string) error
	// SetMetadata sets the metadata of the session
	SetMetadata(ctx context.Context, metadata map[string]string) error
}

type client[ID IdType] struct {
	store              Store
	permissionSupplier PermissionSupplier[ID]
	keyPrefix          string
}

func NewClient[ID IdType](
	store Store,
	permissionSupplier PermissionSupplier[ID],
) Client[ID] {
	return &client[ID]{
		store:              store,
		permissionSupplier: permissionSupplier,
		keyPrefix:          "authz",
	}
}

func makeToken() string {
	return randomString(32)
}

type LoginOption[ID IdType] struct {
	ID ID
	// == Optional start
	// Metadata is a map of key-value pairs that can be used to store
	Metadata map[string]string
	// Platform is the platform from which the user is logging in
	Platform string
	// == Optional end
}

func (c *client[ID]) Login(ctx context.Context, opt *LoginOption[ID]) (string, error) {
	if opt == nil {
		return "", fmt.Errorf("opt is nil")
	}
	opt.prepare()

	_, token, err := c.getOrCreateSession(ctx, opt.ID, opt.Platform, opt.Metadata)
	if err != nil {
		return "", fmt.Errorf("get or create session: %w", err)
	}

	if err := c.setTokenMappingID(ctx, opt.ID, token); err != nil {
		return "", fmt.Errorf("set token mapping: %w", err)
	}

	return token, nil
}

func (c *client[ID]) GetSession(ctx context.Context) (*Session[ID], error) {
	token, ok := GetToken(ctx)
	if !ok {
		return nil, ErrNoToken
	}

	id, err := c.getIDFromToken(ctx, token)
	if err != nil {
		return nil, err
	}
	sessionKey := c.getSessionKey(id)
	sessionString, err := c.store.Get(ctx, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	var session Session[ID]
	if err := json.Unmarshal([]byte(sessionString), &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, err
}

func (c *client[ID]) SetMetadata(ctx context.Context, metadata map[string]string) error {
	session, err := c.GetSession(ctx)
	if err != nil {
		return err
	}

	session.mergeMetadata(metadata)
	return c.saveSession(ctx, session, c.getSessionKey(session.ID))
}

func (c *client[ID]) Set(ctx context.Context, key, value string) error {
	return c.SetMetadata(ctx, map[string]string{key: value})
}

func (c *client[ID]) getOrCreateSession(ctx context.Context, id ID, plat string, metadata map[string]string) (*Session[ID], string, error) {
	sessionKey := c.getSessionKey(id)
	oldSession, err := c.store.Get(ctx, sessionKey)
	if err != nil {
		if !errors.Is(err, ErrStoreValueNotFound) {
			return nil, "", fmt.Errorf("get old session: %w", err)
		}
	}

	var session *Session[ID]
	if oldSession != "" {
		session, err = c.parseSession(oldSession)
		if err != nil {
			return nil, "", fmt.Errorf("parse old session: %w", err)
		}
	} else {
		session = &Session[ID]{
			Model: &Model[ID]{
				ID:       id,
				Metadata: metadata,
			},
		}
	}

	token := makeToken()
	session.addToken(token, plat)
	session.mergeMetadata(metadata)
	if err := c.saveSession(ctx, session, sessionKey); err != nil {
		return nil, "", err
	}
	return session, token, nil
}

func (c *client[ID]) saveSession(ctx context.Context, session *Session[ID], sessionKey string) error {
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal new session: %w", err)
	}

	if err := c.store.Set(ctx, sessionKey, string(sessionBytes)); err != nil {
		return fmt.Errorf("set session: %w", err)
	}
	return nil
}

func (o *LoginOption[ID]) prepare() {
	if o.Platform == "" {
		o.Platform = "default"
	}
}

func (c *client[ID]) getSessionKey(id ID) string {
	return fmt.Sprintf("%s:%v", c.keyPrefix, id)
}

func (c *client[ID]) parseSession(session string) (*Session[ID], error) {
	var result Session[ID]
	if err := json.Unmarshal([]byte(session), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *client[ID]) getTokenMappingKey(token string) string {
	return fmt.Sprintf("%s:token:%s", c.keyPrefix, token)
}

func (c *client[ID]) setTokenMappingID(ctx context.Context, id ID, token string) error {
	var m = idMapping[ID]{ID: id}

	if bytes, err := json.Marshal(&m); err != nil {
		return err
	} else {
		return c.store.Set(ctx, c.getTokenMappingKey(token), string(bytes))
	}
}

type idMapping[ID IdType] struct {
	ID ID `json:"id"`
}

func (c *client[ID]) getIDFromToken(ctx context.Context, token string) (ID, error) {
	key := c.getTokenMappingKey(token)
	value, err := c.store.Get(ctx, key)
	if err != nil {
		if errors.Is(err, ErrStoreValueNotFound) {
			return lo.Empty[ID](), fmt.Errorf("token not found")
		}
	}
	var m idMapping[ID]
	if err := json.Unmarshal([]byte(value), &m); err != nil {
		return lo.Empty[ID](), fmt.Errorf("unmarshal id mapping: %w", err)
	}
	return m.ID, nil
}
