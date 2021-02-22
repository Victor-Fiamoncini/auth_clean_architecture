package loaduserbyemailrepository

import (
	"context"

	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/domain/entities"
	shared_custom_errors "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LoadUserByEmailRepository struct
type LoadUserByEmailRepository struct {
	Email     string
	User      entities.IUser
	UserModel *mongo.Collection
}

// NewLoadUserByEmailRepository func
func NewLoadUserByEmailRepository(userModel *mongo.Collection) ILoadUserByEmailRepository {
	return &LoadUserByEmailRepository{
		UserModel: userModel,
	}
}

// GetEmail LoadUserByEmailRepository method
func (luber *LoadUserByEmailRepository) GetEmail() string {
	return luber.Email
}

// SetEmail LoadUserByEmailRepository method
func (luber *LoadUserByEmailRepository) SetEmail(email string) {
	luber.Email = email
}

// GetUser LoadUserByEmailRepository method
func (luber *LoadUserByEmailRepository) GetUser() entities.IUser {
	return luber.User
}

// SetUser LoadUserByEmailRepository method
func (luber *LoadUserByEmailRepository) SetUser(user entities.IUser) {
	luber.User = user
}

// Load LoadUserByEmailRepository method
func (luber *LoadUserByEmailRepository) Load() (entities.IUser, shared_custom_errors.IDefaultError) {
	if luber.Email == "" {
		return nil, shared_custom_errors.NewMissingParamError("Email")
	}

	user := entities.NewUser()
	ctx := context.Background()

	result := luber.UserModel.FindOne(ctx, bson.D{{
		Key:   "email",
		Value: luber.Email,
	}}, options.FindOne().SetProjection(bson.M{
		"password": 1,
	}))

	err := result.Decode(user)

	if err != nil {
		return nil, shared_custom_errors.NewDefaultError("LoadUserByEmailRepository.Load()")
	}

	luber.User = user

	return luber.User, nil
}
