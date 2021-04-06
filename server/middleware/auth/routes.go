package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://groceryspend-dev.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

type AuthMiddleware interface {
	VerifySession() gin.HandlerFunc
}

// PassthroughMiddleware allows all traffic and does no checks
type PassthroughMiddleware struct {
}

func NewPassthroughAuthMiddleware() *PassthroughMiddleware {
	return &PassthroughMiddleware{}
}

func (p *PassthroughMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}

// Auth0JwtAuthMiddleware leverages Auth0 for auth/authz and session management
type Auth0JwtAuthMiddleware struct {
	middleware *jwtmiddleware.JWTMiddleware
}

func NewAuth0JwtAuthMiddleware() *Auth0JwtAuthMiddleware {

	jwtmiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			// verify 'aud' claim TODO: move this to a configuration
			aud := "https://bknight.dev.groceryspend.io"

			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}

			// verify 'iss' claim TODO: move this to a configuration
			iss := "https://groceryspend-dev.us.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
	})

	return &Auth0JwtAuthMiddleware{
		middleware: jwtmiddleware,
	}
}

func (m *Auth0JwtAuthMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		err := m.middleware.CheckJWT(c.Writer, c.Request)
		if err != nil {
			// token not found
			fmt.Println(err.Error())
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
		}

		// TODO: add checks on scopes
	}

	return gin.HandlerFunc(fn)
}

func GetUserIdFromSession(r http.Request) string {
	return r.URL.User.String()
}
