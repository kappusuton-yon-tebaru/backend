package role

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (r *Repository) CreateRole(ctx context.Context, dto CreateRoleDTO) (string, error) {
	role := bson.M{
		"role_name": dto.RoleName,
		"org_id": dto.OrgId,
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
	update := map[string]any{
		"$set": map[string]any{
			"role_name": dto.RoleName,
			"updated_at":    time.Now(), 
		},
	}
	// Update the role in MongoDB
	result, err := r.role.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error updating role:", err)
		return "", fmt.Errorf("error updating role: %v", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		return "", fmt.Errorf("role not found")
	}

	return objID.Hex(), nil
}

func (r *Repository) DeleteRole(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.role.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error deleting role:", err)
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *Repository) AddPermissionToRole(ctx context.Context, dto ModifyPermissionDTO, roleID string) (string, error) {
	objID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectIDFromHex err")
		return "", err
	}

	permission := map[string]any{
		"_id":            bson.NewObjectID(), 
		"permission_name": dto.PermissionName,
		"action": dto.Action,
		"resource_id": dto.ResourceId,
	}
	update :=  map[string]any{
		"$push":  map[string]any{
			"permissions": permission,
		}}

	// Update the role in MongoDB
	result, err := r.role.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error adding permission to role:", err)
		return "", fmt.Errorf("error adding permission to role: %v", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		return "", fmt.Errorf("role not found")
	}

	return objID.Hex(), nil
}

func (r *Repository) UpdatePermission(ctx context.Context, dto ModifyPermissionDTO, roleID string, permID string) (string, error) {
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

	update :=  map[string]any{
		"$set":  map[string]any{
			"permissions.$.permission_name": dto.PermissionName,
            "permissions.$.action":          dto.Action,
            "permissions.$.resource_id":     dto.ResourceId,
		}}

	// Update the role in MongoDB
	result, err := r.role.UpdateOne(ctx, bson.M{"_id": roleObjID, "permissions._id": permObjID}, update)
	if err != nil {
		log.Println("Error updating permission in role:", err)
		return "", fmt.Errorf("error updating permission in role: %v", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		return "", fmt.Errorf("role not found")
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

	update := map[string]any{
		"$pull": map[string]any{
			"permissions": map[string]any{"_id": permObjID},
		},
	}

	result, err := r.role.UpdateOne(ctx, bson.M{"_id": roleObjID}, update)
	if err != nil {
		log.Println("Error deleting permission from role:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}
