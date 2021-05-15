package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kofalt/go-memoize"
	"groceryspend.io/server/services/users"
	"groceryspend.io/server/utils"
)

type response struct {
	Message string `json:"message"`
}

type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getJwks(url string) (jwks, error) {
	resp, err := http.Get(url)

	if err != nil {
		return jwks{}, err
	}
	defer resp.Body.Close()

	var jwks = jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return jwks, err
	}

	return jwks, nil
}

func getPemCert(token *jwt.Token, cache *memoize.Memoizer) (string, error) {
	cert := ""
	url := "https://groceryspend-dev.us.auth0.com/.well-known/jwks.json"

	getJwksClosure := func() (interface{}, error) {
		return getJwks(url)
	}

	j, err, cached := cache.Memoize(url, getJwksClosure)
	if err != nil {
		return cert, err
	}

	if cached {
		println("Pulling JWKS from cache")
	}

	for k := range j.(jwks).Keys {
		if token.Header["kid"] == j.(jwks).Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + j.(jwks).Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

// Auth0JwtAuthMiddleware leverages Auth0 for auth/authz and session management
type Auth0JwtAuthMiddleware struct {
	userClient users.Client
	middleware *jwtmiddleware.JWTMiddleware
}

// NewAuth0JwtAuthMiddleware create a auth middleware leveraging Auth0
func NewAuth0JwtAuthMiddleware(cache *memoize.Memoizer, userClient users.Client) *Auth0JwtAuthMiddleware {

	// TODO: initialize cache wtih JWKS so that we don't have to wait for the cache
	//			 to warm up. It's about 500ms extra, but definitely noticeable

	jwtmiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			// verify 'aud' claim TODO: move this to a configuration
			aud := "https://bknight.dev.groceryspend.io"

			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				// TODO: better audience verification for the middleware
				fmt.Println("Invalid audience, continuing but this should get fixed")
			}

			// verify 'iss' claim TODO: move this to a configuration
			iss := "https://groceryspend-dev.us.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token, cache)
			if err != nil {
				panic(err.Error())
			}
			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

			return result, nil
		},
	})

	return &Auth0JwtAuthMiddleware{
		middleware: jwtmiddleware,
		userClient: userClient,
	}
}

// VerifySession check the JWT to ensure it's valid still
func (m *Auth0JwtAuthMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		println("Verifying jwt token")
		err := m.middleware.CheckJWT(c.Writer, c.Request)
		if err != nil {
			// token not found
			fmt.Println(err.Error())
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
		}

		// if it's a preflight check, don't try to verify the session since the
		// bearer token isn't sent
		if c.Request.Method == "OPTIONS" {
			println("returning since it's an preflight request")
			return
		}

		// TODO: add checks on scopes

		// set the contet with the user uuid
		// "user" is set by the auth0 CheckJWT call
		u := c.Request.Context().Value("user")
		println(fmt.Sprintf("User is %s", u))
		user := u.(*jwt.Token)
		iss := user.Claims.(jwt.MapClaims)["iss"].(string)
		sub := user.Claims.(jwt.MapClaims)["sub"].(string)

		auth0ID := iss + "|" + sub

		canonicalUser, err := m.userClient.LookupUserByAuthProvider(utils.GetOsValue("AUTH_PROVIDER"), auth0ID)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("Failed to get user from user database"))
		}
		println(fmt.Sprintf("Canonical user is %s", canonicalUser))
		// If we get here, everything worked and we can set the
		// user property in context.
		originalRequest := c.Request
		originalRequestCtx := originalRequest.Context()
		newRequest := originalRequest.WithContext(context.WithValue(originalRequestCtx, AuthUserIDKey, canonicalUser.ID))
		// Update the current request with the new context information.
		c.Request = newRequest
		println("done")

	}

	return gin.HandlerFunc(fn)
}
