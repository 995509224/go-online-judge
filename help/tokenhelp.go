package help

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"hykoj/models"
	"time"
)

type User struct {
	Name     string `json:"Name"`
	Identity string `json:"identity"`
	Isroot   int    `json:"isadmin"`
	jwt.StandardClaims
}

var screat = []byte("hykyyds")

func Gettoken(data *models.UserBasic) (string, error) {
	temp := &User{
		Name:     data.Name,
		Identity: data.Identity,
		Isroot:   data.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "kobe",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, temp)
	return token.SignedString(screat)
}
func Parsetoken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &User{}, func(token *jwt.Token) (interface{}, error) {
		return screat, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*User); ok && token.Valid {
		fmt.Println(claims)
		return claims, nil
	}
	return nil, errors.New("invalid token")

}
