package model

import (
	"encoding/json"
	"fmt"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User use to save user detail
type User struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;"`
	Name         string          `json:"name"`
	NIK          string          `json:"nik"`
	Role         proto.User_Role `json:"role"`
	Username     string          `json:"username"`
	PasswordHash string          `json:"password_hash"`
}

// UserFromJSON use to create user from json string
func UserFromJSON(str string) (user User, err error) {
	err = json.Unmarshal([]byte(str), &user)

	return
}

func ValidateToken(tokenString string) (User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := claims["user"].(map[string]interface{})
		id, _ := uuid.FromString(data["id"].(string))
		role := data["role"].(float64)
		logrus.Println(data)
		return User{
			id,
			data["name"].(string),
			data["nik"].(string),
			proto.User_Role(role),
			data["username"].(string),
			data["password_hash"].(string),
		}, nil
	} else {
		logrus.Error(err)
		return User{}, err
	}
}

func (u User) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": u,
	})

	return token.SignedString([]byte("secret"))
}

func (u User) toJSON() (string, error) {
	data, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func NewUser(
	id uuid.UUID,
	name string,
	nik string,
	role proto.User_Role,
	username string,
	password string,
) (User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}
	return User{
		id,
		name,
		nik,
		role,
		username,
		string(bytes),
	}, nil
}
