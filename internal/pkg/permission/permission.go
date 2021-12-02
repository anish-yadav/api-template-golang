package permission

import (
	"errors"

	"github.com/anish-yadav/api-template-golang/internal/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Permission struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

const collection = "permissions"

func NewPermission(name string, p []string) *Permission {
	return &Permission{name, p}
}

func ClearDB() error {
	return db.DelAll(collection)
}

// TODO: should return error
func GetPermissionByName(name string) *Permission {
	permissionDB, err := db.GetByPKey(collection, "name", name)
	if err != nil {
		return nil
	}
	bsonBytes, err := bson.Marshal(permissionDB)
	if err != nil {
		return nil
	}
	var permission Permission
	if err = bson.Unmarshal(bsonBytes, &permission); err != nil {
		return nil
	}
	return &permission
}

func (p *Permission) AddToDB() (string, error) {
	bin, err := bson.Marshal(p)
	if err != nil {
		return "", errors.New("failed to marshal permission data")
	}
	var bsonData bson.D
	_ = bson.Unmarshal(bin, &bsonData)
	return db.InsertOne(collection, bsonData)
}

func (p *Permission) HasPermission(permission string) bool {
	if len(permission) == 0 {
		return true
	}
	for _, p := range p.Permissions {
		if permission == p {
			return true
		}
	}
	return false
}
