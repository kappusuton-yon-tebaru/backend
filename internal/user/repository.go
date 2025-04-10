package user

import (
	"context"
	"fmt"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	user *mongo.Collection
}

func NewRepository(db *mongo.Database) (*Repository, error) {
	repo := &Repository{
		user: db.Collection("users"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) ensureIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys: map[string]any{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.user.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	cur, err := r.user.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	users := make([]models.User, 0)

	for cur.Next(ctx) {
		var user UserDTO

		err = cur.Decode(&user)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		users = append(users, DTOToUser(user))
	}

	return users, nil
}

func (r *Repository) CreateUser(ctx context.Context, dto UserCredentialDTO) (string, error) {
	result, err := r.user.InsertOne(ctx, dto)
	if err != nil {
		return bson.NilObjectID.Hex(), err
	}

	userId := result.InsertedID.(bson.ObjectID)

	return userId.Hex(), nil
}

func (r *Repository) DeleteUser(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.user.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	filter := map[string]any{
		"email": email,
	}

	result := r.user.FindOne(ctx, filter)
	if result.Err() != nil {
		return models.User{}, result.Err()
	}

	var dto UserDTO
	if err := result.Decode(&dto); err != nil {
		return models.User{}, err
	}

	return DTOToUser(dto), nil
}

func (r *Repository) AddRole(ctx context.Context, userID string, roleID string) (string, error) {
	userObjID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("userID FromHex err")
		return "", err
	}

	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("roleID FromHex err")
		return "", err
	}

	update := bson.M{
		"$addToSet": bson.M{
			"role_ids": roleObjID,
		},
	}

	result, err := r.user.UpdateOne(ctx, bson.M{"_id": userObjID}, update)
	if err != nil {
		log.Println("Error adding role to user:", err)
		return "", fmt.Errorf("error adding user to role: %v", err)
	}

	if result.MatchedCount == 0 {
		return "", fmt.Errorf("user not found")
	}

	return userObjID.Hex(), nil
}
func (r *Repository) RemoveRole(ctx context.Context, userID string, roleID string) (int64, error) {
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	userObjID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	update := bson.M{
		"$pull": bson.M{
			"role_ids": roleObjID,
		},
	}

	result, err := r.user.UpdateOne(ctx, bson.M{"_id": userObjID}, update)
	if err != nil {
		log.Println("Error removing role from user:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (r *Repository) RemoveRoleInAllUser(ctx context.Context, roleID string) (int64, error) {
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	update := bson.M{
		"$pull": bson.M{
			"role_ids": roleObjID,
		},
	}

	result, err := r.user.UpdateMany(ctx, bson.M{"role_ids": roleObjID}, update)
	if err != nil {
		log.Println("Error removing role from users:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}
