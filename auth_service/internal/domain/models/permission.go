package models

type Permission string

const (
	RoleUser  Permission = "User"
	RoleAdmin Permission = "Admin"
)
