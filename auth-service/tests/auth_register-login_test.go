package tests

import (
	"AuthService/tests/suite"
	"math/rand"
	authv1 "protos/gen/go/auth-service"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRegisterAndLogin(t *testing.T) {
	ctx, st := suite.New(t)

	email := strconv.Itoa(rand.Intn(100)) + "_test@example.com"
	password := "testpassword"

	responseReg, err := st.AuthService.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	if responseReg.UserId <= 0 {
		t.Fatalf("expected valid user ID, got: %d", responseReg.UserId)
	}

	responseLogin, err := st.AuthService.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    1,
	})

	if err != nil {
		t.Fatalf("failed to login user: %v", err)
	}

	token := responseLogin.GetToken()
	if token == "" {
		t.Fatalf("expected valid token, got empty string")
	}

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret1"), nil
	})

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("failed to extract claims from token")
	}

	if claims["email"].(string) != email {
		t.Fatalf("expected email claim to be %s, got %s", email, claims["email"])
	}

	if int(claims["app_id"].(float64)) != 1 {
		t.Fatalf("expected app_id claim to be 1, got %v", claims["app_id"])
	}

	if int64(claims["user_id"].(float64)) != responseReg.UserId {
		t.Fatalf("expected user_id claim to be %d, got %v", responseReg.UserId, claims["user_id"])
	}

	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		t.Fatalf("expected token to be valid, but it is expired")
	}
}

func TestRegisterWithExistingEmail(t *testing.T) {
	ctx, st := suite.New(t)

	email := "33_test@example.com"
	password := "testpassword"

	_, err := st.AuthService.Register(
		ctx, &authv1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)

	if err == nil {
		t.Fatalf("expected error when registering with existing email, got nil")
	}
}

func TestLoginWithExpiredToken(t *testing.T) {
	ctx, st := suite.New(t)

	email := "33_test@example.com"
	password := "testpassword"

	responseLogin, err := st.AuthService.Login(
		ctx, &authv1.LoginRequest{
			Email:    email,
			Password: password,
			AppId:    1,
		},
	)

	if err != nil {
		t.Fatalf("failed to login user: %v", err)
	}

	if responseLogin.GetToken() == "" {
		t.Fatalf("failed to get token from login response")
	}

	token := responseLogin.GetToken()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret1"), nil
	})

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("failed to extract claims from token")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().Add(time.Hour*2).Before(expirationTime) {
		t.Fatalf("expected token to be expired, but it is still valid")
	}
}
