package views

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 宣告JWT結構
type Claims struct {
	User_id uint
	S_id string
	jwt.StandardClaims
}

// jwt key
var jwt_secret = []byte(os.Getenv("JWT_KEY"))

// get token
func Jwt_token_get(claims Claims) (string, error) {
	token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := token_claim.SignedString(jwt_secret)
	return token, err
}

// validate JWT
func AuthRequired(c *gin.Context) (cl *Claims, err error) {
	// 取token
	token, err := c.Cookie("Authorization")
	if err != nil {
		return nil ,err
	}

	// 看是否有設定jwt_secret
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		if jwt_secret == nil {
			err:= errors.New("can't find jw_secret")
			return nil, err
		}
		return jwt_secret, nil
	})
	
	// token mistake guide:
	// parse and validate token for six things:
	// validationErrorMalformed => token is malformed
	// validationErrorUnverifiable => token could not be verified because of signing problems
	// validationErrorSignatureInvalid => signature validation failed
	// validationErrorExpired => exp validation failed
	// validationErrorNotValidYet => nbf validation failed
	// validationErrorIssuedAt => iat validation failed
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors & jwt.ValidationErrorUnverifiable != 0{
				message = "token could not be verified because of signing problems"
			} else if ve.Errors & jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors & jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors & jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		err := errors.New(message)
		return nil, err
	}

	// 驗證成功回傳資料
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	// 失敗
	} else {
		err := errors.New("token data wrong")
		return nil, err
	}
}