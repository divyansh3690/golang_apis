package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type CustomClains struct {
	Email     string
	Full_name string
	Role      string
	User_id   int
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}

func LoadENV(value string) (key string, err error) {
	err = godotenv.Load()
	if err != nil {
		return "", err
	}
	key = os.Getenv(value)
	return key, nil
}

func GenerateToken(email string, full_name string, role string, user_id int) (signerdAccessToken string, signedRefreshToken string, err error) {
	// i faced an error here : note here it cant be like we give Email : Email and then just provide value of jwt.Regitration value, either give key value pair or just the values not both .
	accessClaims := &CustomClains{
		email,
		full_name,
		role,
		user_id,
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)), // i have made expire time as 15 for now. change it later
		},
	}

	refreshClaims := &RefreshClaims{
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)), // i have made expire time as 15 for now. change it later
		},
	}

	// access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signingKey, err := LoadENV("ACCESS_SIGNING_KEY")
	if err != nil {
		return "", "", err
	}
	signerdAccessToken, err = accessToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	// refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signingKey, err = LoadENV("REFRESH_SIGNING_KEY")
	if err != nil {
		return "", "", fmt.Errorf("error occured while creating access token:%v", err)
	}
	signedRefreshToken, err = refreshToken.SignedString([]byte(signingKey)) //this takes interface not string ->error i faced
	if err != nil {
		return "", "", fmt.Errorf("error occured while creating refresh token:%v", err)
	}

	return signerdAccessToken, signedRefreshToken, nil
}

// CREATE A UPDATE MONGODB funcitn
