package role

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	role *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		role: db.Collection("roles"),
	}
}

func (r *Repository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	cur, err := r.role.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	roles := make([]models.Role, 0)

	for cur.Next(ctx) {
		var role RoleDTO

		err = cur.Decode(&role)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		roles = append(roles, DTOToRole(role))
	}

	return roles, nil
}

func (r *Repository) GetRoleByID(ctx context.Context, roleID string) (models.Role, error) {
	objID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectIDFromHex err:", err)
		return models.Role{}, err
	}

	var roleDTO RoleDTO
	err = r.role.FindOne(ctx, bson.M{"_id": objID}).Decode(&roleDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Role{}, fmt.Errorf("role not found")
		}
		log.Println("Error finding role:", err)
		return models.Role{}, fmt.Errorf("error finding role: %v", err)
	}

	role := DTOToRole(roleDTO)

	return role, nil
}

func (r *Repository) CreateRole(ctx context.Context, dto CreateRoleDTO) (string, error) {
	role := bson.M{
		"role_name":   dto.RoleName,
		"org_id":      dto.OrgId,
		"permissions": []models.Permission{},
	}

	result, err := r.role.InsertOne(ctx, role)
	if err != nil {
		log.Println("Error inserting role:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting role: %v", err)
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *Repository) UpdateRole(ctx context.Context, dto UpdateRoleDTO, roleID string) (string, error) {
	objID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectIDFromHex err")
		return "", err
	}

	// Convert incoming permissions to the correct format
	newPermissions := make([]bson.M, 0, len(dto.Permissions))
	for _, perm := range dto.Permissions {
		resourceObjID := perm.ResourceId
		if err != nil {
			log.Printf("Invalid resource_id: %v\n", perm.ResourceId)
			return "", fmt.Errorf("invalid resource_id: %v", perm.ResourceId)
		}

		if !enum.IsValidPermissionActions(perm.Action) {
			return "", fmt.Errorf("invalid action: %v", perm.Action)
		}

		newPermissions = append(newPermissions, bson.M{
			"_id":             bson.NewObjectID(),
			"permission_name": perm.PermissionName,
			"action":          perm.Action,
			"resource_id":     resourceObjID,
		})
	}

	update := bson.M{
		"$set": bson.M{
			"role_name":  dto.RoleName,
			"permissions": newPermissions,
			"updated_at": time.Now(),
		},
	}

	// Update the role in MongoDB
	result, err := r.role.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error updating role:", err)
		return "", fmt.Errorf("error updating role: %v", err)
	}

	if result.MatchedCount == 0 {
		return "", fmt.Errorf("role not found")
	}

	return objID.Hex(), nil
}


func (r *Repository) DeleteRole(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.role.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error deleting role:", err)
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *Repository) AddPermission(ctx context.Context, dto ModifyPermissionDTO, roleID string) (string, error) {
	if !enum.IsValidPermissionActions(dto.Action) {
		return "", fmt.Errorf("invalid resource type: %v", dto.Action)
	}
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return "", err
	}
	// check role exists
	count, err := r.role.CountDocuments(ctx, bson.M{"_id": roleObjID})
	if err != nil {
		log.Println("Error checking role existence:", err)
		return "", fmt.Errorf("error checking role existence: %v", err)
	}
	if count == 0 {
		log.Println("This role doesn't exist")
		return "", fmt.Errorf("this role doesn't exist")
	}
	// check if action and resource is unique
	uniquePermissionFilter := bson.M{
		"_id": roleObjID,
		"permissions": bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					"action":      dto.Action,
					"resource_id": dto.ResourceId,
				},
			},
		},
	}
	count, err = r.role.CountDocuments(ctx, uniquePermissionFilter)
	if err != nil {
		log.Println("Error checking role/permission existence:", err)
		return "", fmt.Errorf("error checking role/permission existence: %v", err)
	}
	if count == 0 {
		log.Println("Permission with this action on this resource already exists")
		return "", fmt.Errorf("Permission with this action on this resource already exists")
	}
	// add new perm to role
	permission := bson.M{
		"_id":             bson.NewObjectID(),
		"permission_name": dto.PermissionName,
		"action":          dto.Action,
		"resource_id":     dto.ResourceId,
	}
	update := bson.M{
		"$push": bson.M{
			"permissions": permission,
		},
	}
	_, err = r.role.UpdateOne(ctx, bson.M{"_id": roleObjID}, update)
	if err != nil {
		log.Println("Error adding permission to role:", err)
		return "", fmt.Errorf("error adding permission to role: %v", err)
	}

	return roleObjID.Hex(), nil
}

func (r *Repository) UpdatePermission(ctx context.Context, dto ModifyPermissionDTO, roleID string, permID string) (string, error) {
	if !enum.IsValidPermissionActions(dto.Action) {
		return "", fmt.Errorf("invalid resource type: %v", dto.Action)
	}
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return "", err
	}

	permObjID, err := bson.ObjectIDFromHex(permID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return "", err
	}
	// check role exists
	count, err := r.role.CountDocuments(ctx, bson.M{"_id": roleObjID})
	if err != nil {
		log.Println("Error checking role existence:", err)
		return "", fmt.Errorf("error checking role existence: %v", err)
	}
	if count == 0 {
		log.Println("This role doesn't exist")
		return "", fmt.Errorf("this role doesn't exist")
	}
	// check permission exists
	count, err = r.role.CountDocuments(ctx, bson.M{"_id": roleObjID, "permissions._id": permObjID})
	if err != nil {
		log.Println("Error checking permission existence:", err)
		return "", fmt.Errorf("error checking permission existence: %v", err)
	}
	if count == 0 {
		log.Println("This permission doesn't exist in this role")
		return "", fmt.Errorf("this permission doesn't exist in this role")
	}
	// check if action and resource is unique
	filter := bson.M{
		"_id": roleObjID,
		"permissions": bson.M{
			"$elemMatch": bson.M{
				"_id":         permObjID,
				"action":      bson.M{"$ne": dto.Action},     // Ensure new action is unique
				"resource_id": bson.M{"$ne": dto.ResourceId}, // Ensure new resource_id is unique
			},
		},
	}
	count, err = r.role.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("Error checking role/permission existence:", err)
		return "", fmt.Errorf("error checking role/permission existence: %v", err)
	}
	if count == 0 {
		log.Println("Permission with this action on this resource already exists")
		return "", fmt.Errorf("Permission with this action on this resource already exists")
	}
	// update permission
	update := bson.M{
		"$set": bson.M{
			"permissions.$.permission_name": dto.PermissionName,
			"permissions.$.action":          dto.Action,
			"permissions.$.resource_id":     dto.ResourceId,
		}}

	_, err = r.role.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating permission in role:", err)
		return "", fmt.Errorf("error updating permission in role: %v", err)
	}

	return roleObjID.Hex(), nil
}

func (r *Repository) DeletePermission(ctx context.Context, roleID string, permID string) (int64, error) {
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	permObjID, err := bson.ObjectIDFromHex(permID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	update := bson.M{
		"$pull": bson.M{
			"permissions": bson.M{"_id": permObjID},
		},
	}

	result, err := r.role.UpdateOne(ctx, bson.M{"_id": roleObjID}, update)
	if err != nil {
		log.Println("Error deleting permission from role:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}