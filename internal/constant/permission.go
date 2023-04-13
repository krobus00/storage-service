package constant

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionRead   ActionType = "READ"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"

	PermissionFullAccess  = "FULL_ACCESS"
	PermissionGuestAccess = "GUEST_ACCESS"

	PermissionObjectAll         = "OBJECT_ALL"
	PermissionObjectCreate      = "OBJECT_CREATE"
	PermissionObjectRead        = "OBJECT_READ"
	PermissionObjectReadPrivate = "OBJECT_READ_PRIVATE"
	PermissionObjectDelete      = "OBJECT_DELETE"
	PermissionObjectModifyOther = "OBJECT_MODIFY_OTHER"
)
