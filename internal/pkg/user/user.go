package user

import (
	"errors"

	"github.com/anish-yadav/api-template-golang/internal/constants"
	"github.com/anish-yadav/api-template-golang/internal/pkg/db"
	"github.com/anish-yadav/api-template-golang/internal/pkg/permission"
	"github.com/anish-yadav/api-template-golang/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserDb struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Type     string             `json:"type"`
	Detail   interface{}        `json:"details"`
}

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Type        string             `json:"type"`
	Permissions []string           `json:"permissions"`
}

// TODO: should return error

func (user *UserDb) ToResponse() *User {
	var permissions []string
	per := permission.GetPermissionByName(user.Type)
	if per != nil {
		permissions = per.Permissions
	}
	// sending permissions array too
	return &User{user.ID, user.Name, user.Email, user.Type, permissions}
}

const collection = "users"

func NewUser(name string, email string, password string, typ string) *UserDb {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	id := primitive.NewObjectID()
	return &UserDb{id, name, email, string(hashedPwd), typ, nil}
}

func GetUserById(id string) *UserDb {
	userDb, err := db.GetByID(collection, id)
	if err != nil {
		return nil
	}
	bsonBytes, err := bson.Marshal(userDb)
	if err != nil {
		return nil
	}
	var user UserDb
	if err = bson.Unmarshal(bsonBytes, &user); err != nil {
		return nil
	}
	return &user
}

func GetUserByEmail(email string) *UserDb {
	userDb, err := db.GetByPKey(collection, "email", email)
	if err != nil {
		return nil
	}
	bsonBytes, err := bson.Marshal(userDb)
	if err != nil {
		return nil
	}
	var user UserDb
	if err = bson.Unmarshal(bsonBytes, &user); err != nil {
		return nil
	}
	return &user
}

func DeleteUserByID(id string) error {
	return db.DelByID(collection, id)
}

func (user *UserDb) AddToDB() (string, error) {
	bin, err := bson.Marshal(user)
	if err != nil {
		return "", errors.New("failed to marshal user data")
	}
	var bsonData bson.D
	err = bson.Unmarshal(bin, &bsonData)
	return db.InsertOne(collection, bsonData)
}

func (user *UserDb) ChangePassword(old string, new string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(old))
	if err != nil {
		return errors.New(constants.PasswordMismatch)
	}
	return user.ResetPassword(new)
}

func (user *UserDb) ResetPassword(new string) error {
	newHashedPwd, _ := bcrypt.GenerateFromPassword([]byte(new), bcrypt.DefaultCost)
	user.Password = string(newHashedPwd)
	resetQuery := bson.D{{"$set", bson.D{{"password", user.Password}}}}
	return db.UpdateItem(collection, user.ID.Hex(), resetQuery)
}

func (user *UserDb) Login(pass string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return "", err
	}
	data := map[string]string{
		"user_id": user.ID.Hex(),
	}
	token, err := util.CreateToken(data)
	if err != nil {
		return "", err
	}
	return token, nil
}
