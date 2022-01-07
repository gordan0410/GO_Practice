package tools

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 宣告JWT結構
type Claims struct {
	User_id uint
	S_id    string
	jwt.StandardClaims
}

// get token
func Jwt_token_init(c *gin.Context, claims Claims) (string, error) {
	// get configs
	configs, err := Get_configs(c)
	if err != nil {
		return "", err
	}

	// set jwt maxage
	now := time.Now()
	claims.StandardClaims.ExpiresAt = now.Add(time.Duration(configs.Jwt.JwtMaxage) * time.Second).Unix()

	// set jwt key
	jwt_key := []byte(configs.Jwt.JwtKey)
	token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := token_claim.SignedString(jwt_key)
	if err != nil {
		return "", err
	}
	return token, nil
}

// validate JWT
func AuthRequired(c *gin.Context, token string) (cl *Claims, err error) {
	// get configs
	configs, err := Get_configs(c)
	if err != nil {
		return nil, err
	}

	// load jw_config
	jwt_key := []byte(configs.Jwt.JwtKey)

	// 看是否有設定jwt_secret
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		if jwt_key == nil {
			err := errors.New("can't find jw_secret")
			return nil, err
		}
		return jwt_key, nil
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
