package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"lah-2022/backend/ent"
	"lah-2022/backend/utils"

	gAuth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type api struct {
	client *ent.Client
}

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:3500/auth/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/gmail.readonly",
	},
	Endpoint: google.Endpoint,
}

func RegisterRoutes(client *ent.Client, g *echo.Group) {
	a := api{client}
	g.GET("/login", a.handleRedirect)
	g.GET("/callback", a.handleCallback)
}

func (a api) handleRedirect(c echo.Context) error {
	email := c.FormValue("email")
	callback := c.FormValue("callback")
	if email == "" || callback == "" {
		return c.JSON(http.StatusBadRequest, utils.APIError{
			ErrorMessage: "bad parameters",
		})
	}

	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 32)
	rand.Read(b)
	stateNonce := base64.URLEncoding.EncodeToString(b)

	stateToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nonce":    stateNonce,
		"callback": callback,
	}).SignedString([]byte(os.Getenv("SIGNING_KEY")))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.APIError{
			ErrorMessage: "problem making state token",
		})
	}

	cookie := http.Cookie{Name: "oauthstate", Value: stateToken, Expires: expiration}
	c.SetCookie(&cookie)

	u := GoogleOauthConfig.AuthCodeURL(
		stateToken,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("login_hint", email),
		oauth2.SetAuthURLParam("prompt", "login"),
	)

	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (a api) handleCallback(c echo.Context) error {
	oauthState, _ := c.Cookie("oauthstate")

	if c.FormValue("state") != oauthState.Value {
		return c.JSON(http.StatusBadRequest, utils.APIError{
			ErrorMessage: "invalid state",
		})
	}
	stateClaims := ParseJWT(c.FormValue("state"))
	if stateClaims == nil || stateClaims["callback"].(string) == "" {
		return c.JSON(http.StatusBadRequest, utils.APIError{
			ErrorMessage: "invalid state",
		})
	}

	token, err := GoogleOauthConfig.Exchange(c.Request().Context(), c.FormValue("code"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.APIError{
			ErrorMessage: "invalid code",
		})
	}

	oauth2Service, err := gAuth.NewService(
		c.Request().Context(),
		option.WithTokenSource(GoogleOauthConfig.TokenSource(c.Request().Context(), token)),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.APIError{
			ErrorMessage: "problems.",
		})
	}
	userInfo, err := oauth2Service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.APIError{
			ErrorMessage: "unable to fetch account data from gmail",
		})
	}

	sessionToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           userInfo.Id,
		"email":         userInfo.Email,
		"refresh_token": token.RefreshToken,
	})

	tokenString, err := sessionToken.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.APIError{
			ErrorMessage: "unable to fetch account data from gmail",
		})
	}

	callback, err := url.Parse(stateClaims["callback"].(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.APIError{
			ErrorMessage: "invalid callback",
		})
	}

	callbackValues := callback.Query()

	callbackValues.Add("token", tokenString)
	callbackValues.Add("email", userInfo.Email)
	callback.RawQuery = callbackValues.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, callback.String())
}

func ParseJWT(token string) jwt.MapClaims {
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return nil
	}

	if claims, ok := tokenData.Claims.(jwt.MapClaims); ok && tokenData.Valid {
		return claims
	} else {
		return nil
	}
}
