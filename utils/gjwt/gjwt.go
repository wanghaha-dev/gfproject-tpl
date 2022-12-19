package gjwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wanghaha-dev/gf/frame/g"
	"time"
)

var secret = g.Cfg().GetBytes("jwt.secret")

/*
example:
tokenString, err := GenerateToken(jwt.MapClaims{
	"foo":  "bar",
	"name": "ryan",
	"exp":  time.Now().Add(time.Second * 1).Unix(), // 过期时间, 单位时间戳秒
})
*/

// GenerateToken 生成 token
func GenerateToken(claims jwt.MapClaims) (string, error) {
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(g.Cfg().GetInt("jwt.expires"))).Unix()
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析 token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ok
		return claims, nil
	} else {
		return nil, err
	}
}
