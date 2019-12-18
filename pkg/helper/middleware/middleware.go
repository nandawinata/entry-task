package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nandawinata/entry-task/pkg/common/redis"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
	"github.com/nandawinata/entry-task/pkg/helper/middleware/constants"
)

type TokenPayload struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type TokenClaims struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJwt(payload TokenPayload) (*string, error) {
	expireToken := time.Now().AddDate(0, 0, 1).Unix()
	claimsToken := TokenClaims{
		ID:       payload.ID,
		Username: payload.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "entry-task",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsToken)
	tokenString, err := token.SignedString([]byte(constants.SECRET_KEY))

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	newToken := "Bearer " + tokenString

	return &newToken, nil
}

func ValidateJwt(req *http.Request) (*TokenPayload, error) {
	tokenRaw := req.Header.Get("authorization")

	if len(tokenRaw) == 0 {
		return nil, eh.NewError(http.StatusBadRequest, "Token is required")
	}

	newToken := strings.TrimPrefix(tokenRaw, "Bearer ")

	var jwtAuthResult *TokenPayload
	redisService := redis.New()
	keyID := fmt.Sprintf(constants.REDIS_KEY, newToken)
	err := redisService.Get(keyID, &jwtAuthResult)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if jwtAuthResult != nil {
		return jwtAuthResult, nil
	}

	token, err := jwt.Parse(newToken, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, eh.NewError(http.StatusBadRequest, fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(constants.SECRET_KEY), nil
	})

	var tokenClaims jwt.MapClaims

	if token != nil && err == nil {
		tokenClaims = token.Claims.(jwt.MapClaims)
		if tokenClaims.VerifyExpiresAt(time.Now().Unix(), false) == false {
			return nil, eh.NewError(http.StatusUnauthorized, "Token is expired")
		}

		id := tokenClaims["id"].(float64)

		token := &TokenPayload{
			ID:       uint64(id),
			Username: tokenClaims["username"].(string),
		}

		redisService.Set(keyID, token, time.Minute)

		return token, nil
	}

	return nil, eh.NewError(http.StatusUnauthorized, "Token is invalid")
}
