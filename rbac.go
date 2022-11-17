package main

import "context"

type RBAC struct {
	UserID int32
}

func (r *RBAC) CheckAccess(ctx context.Context, namespace Namespace) bool {
	return false
}

func (r *RBAC) CheckAccessKind(ctx context.Context, namespace Namespace) string {
	return "VIEW"
}

func (r *RBAC) CheckResourceAccess(ctx context.Context, namespace Namespace, resourceId int32) bool {
	return false
}
