package main

import "fmt"

type Role struct {
	ID          int32
	Name        string
	Permissions []*string
}

var mockRoles = []string{
	"DEFAULT",
	"SITE_ADMINISTRATOR",
	"OPERATOR",
}

//func (r *Role) AddPermission(permission *Permission) error {
//	return nil
//}
//
//func (r *Role) AddPermissions(permissions []*Permission) error {
//	return nil
//}
//
//func (r *Role) RemovePermission(permission *Permission) error {
//	return nil
//}
//
//func (r *Role) getPermissions() ([]*Permission, error) {
//	return nil, nil
//}
//
//func (r *Role) CheckNamespaceAccess(namespace string) bool {
//	permissions, err := r.getPermissions()
//	if err != nil {
//		return false
//	}
//
//	if permissions == nil || len(permissions) == 0 {
//		return false
//	}
//
//	for _, p := range permissions {
//		if p.HasNamespaceAccess(namespace) {
//			return true
//		}
//	}
//	return false
//}

func seedRoles(c *Controller) ([]*Role, error) {
	roles := make([]*Role, len(mockRoles), len(mockRoles))

	for idx, name := range mockRoles {
		fmt.Printf("Seeding role %s", name)
		role, err := c.CreateRole(name)
		if err != nil {
			return roles, err
		}

		roles[idx] = role
	}

	return roles, nil
}
