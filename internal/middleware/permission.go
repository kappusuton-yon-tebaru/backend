package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
)
// Check if the user has allowed action on that resourceId
func (m *Middleware) HavePermission(allowedAction enum.PermissionActions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resourceId := ctx.Param("id")
		if resourceId == "" {
			ctx.Next() // assume using get all or post endpoint where the permission will be checked in the handler
			return
		}

		userId, err := utils.GetUserID(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		permissions, err := m.roleService.GetUserPermissions(ctx, userId)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, permission := range permissions {
			// Check if user has permission that relates to the resourceId
			if permission.ResourceId == resourceId {
				if allowedAction == enum.PermissionActionsRead { // if allowedAction is read, check if permission is read or write
					if permission.Action == enum.PermissionActionsRead || permission.Action == enum.PermissionActionsWrite {
						ctx.Next()
						return
					}
				} else if allowedAction == enum.PermissionActionsWrite { // if allowedAction is write, check if permission is write
					if permission.Action == enum.PermissionActionsWrite {
						ctx.Next()
						return
					}
				}
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]any{
			"error": "user does not have permission to access this resource",
		})
		return

	}
}
