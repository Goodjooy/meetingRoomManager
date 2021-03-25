package server

import (
	"net/http"
	"time"
	"todo-web/err"
	"todo-web/server/IOC"
	"todo-web/server/manage"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	ID                 uint   `json:"id"`
	Password           string `json:"paswd"`
	Name               string `json:"name"`
	jwt.StandardClaims `json:"-"`
}
type JWTParm struct {
	Token IOC.Value `ioc:"from:cookie;to:string;name:token;default:"`
}

var secret = []byte("1141451919")
var ignorePath = []interface{}{
	"/api/user/login",
	"/api/user/sign-up",
}
var effectTime = 5 * time.Hour

func GenerateToken(claims *UserClaims) (string, err.Exception) {
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()

	sign, erro := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secret)
	if erro != nil {
		return "", err.GenerateTokenFailure(erro.Error())
	}
	return sign, err.NoExcetion
}

func JWTVerifyIOC(parm JWTParm, Req *http.Request, context *IOC.ConxtextSeter) {
	var token string
	if manage.Contains(Req.RequestURI, ignorePath) {
		return
	}
	parm.Token.Get(&token)

	user, e := parseToken(token)
	context.Set("err", e)
	context.Set("token", *user)
}

func parseToken(tokenStrng string) (*UserClaims, err.Exception) {
	token, erro := jwt.ParseWithClaims(tokenStrng, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if erro != nil {
		return nil, err.TokenFailure(erro.Error())
	}

	claim, ok := token.Claims.(*UserClaims)

	if !ok {
		return nil, err.TokenInvaild("user " + claim.Name + " over Time")
	}
	return claim, err.NoExcetion
}

func Reflesh(old string) (string, err.Exception) {
	jwt.TimeFunc = func() time.Time { return time.Unix(0, 0) }

	token, e := jwt.ParseWithClaims(old, UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if e != nil {
		return "", err.TokenFailure(e.Error())
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", err.TokenInvaild("user " + claims.Name + " over Time")
	}
	jwt.TimeFunc = time.Now

	claims.StandardClaims.ExpiresAt = time.Now().Add(effectTime).Unix()

	return GenerateToken(claims)
}
