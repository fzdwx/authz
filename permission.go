package authz

import "context"

type Permission string

type PermissionSupplier[ID IdType] interface {
	GetPermissions(ctx context.Context, m *Model[ID]) []Permission
}

type DefaultPermissionSupplier[ID IdType] struct {
}

func (d DefaultPermissionSupplier[ID]) GetPermissions(ctx context.Context, m *Model[ID]) []Permission {
	return []Permission{}
}
