package main

type RoleName string

var (
	// DefaultRole is the default role anyone can have. It's automatically
	// assigned to a user at the point of creation. It's not meant to be deleted
	// as permissions for everyone can easily be controlled via this role.
	DefaultRole RoleName = "DEFAULT"

	// AdminRole is for users (or supervisors) that are able to manage and configure
	// this instance.
	AdminRole RoleName = "ADMIN"

	// OperatorRole is for users with minimal access.
	OperatorRole RoleName = "OPERATOR"
)

type Role struct {
	ID          int32
	Name        RoleName
	Permissions []*string
}

func (r *Role) AddPermission(permission *Permission) error {
	return nil
}

func (r *Role) AddPermissions(permissions []*Permission) error {
	return nil
}

func (r *Role) RemovePermission(permission *Permission) error {
	return nil
}

func (r *Role) getPermissions() ([]*Permission, error) {
	return nil, nil
}

func (r *Role) CheckNamespaceAccess(namespace string) bool {
	permissions, err := r.getPermissions()
	if err != nil {
		return false
	}

	if permissions == nil || len(permissions) == 0 {
		return false
	}

	for _, p := range permissions {
		if p.HasNamespaceAccess(namespace) {
			return true
		}
	}
	return false
}
