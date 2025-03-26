package enum

type PermissionActions string

const (
	PermissionActionsRead PermissionActions = "read"
	PermissionActionsWrite PermissionActions = "write"
	PermissionActionsExecute PermissionActions = "execute"
	PermissionActionsBuild PermissionActions = "build"
)

func IsValidPermissionActions(t PermissionActions) bool {
	switch t {
	case PermissionActionsRead, PermissionActionsWrite, PermissionActionsExecute, PermissionActionsBuild:
		return true
	default:
		return false
	}
}
