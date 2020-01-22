package unstable_api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/beevik/guid"

	uuid "github.com/satori/go.uuid"
)

type Credentials struct {
	Password string
	Username string
}

func Signin(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	// Get the expected password from our in memory map
	expectedUser, ok := Accounts[strings.ToUpper(creds.Username)]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedUser.Pwd != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}

	// Create a new random session token
	sessionToken := uuid.NewV4()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	/***************************** b, _ := httputil.DumpRequest(r, false) */ /////////
	//////println(string(b), "END")
	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	session_token := Token{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: time.Now().Add(120 * time.Minute),
	}
	cache.Put(expectedUser, session_token)

	result, err := json.Marshal(session_token)
	return result, err
}

// AuthenticateToken to verify
func AuthenticateToken(r *http.Request) (gid []guid.Guid, err error) {
	c, err := GetToken(r.Header.Get("Token"))
	if err != nil {
		// If the TOKEN not present on header, return an unauthorized status
		return
	}
	token := c.Value

	response, err := cache.Get(token)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		return
	}

	// compares token Expire time and destroy token if Expired
	if response.Token.Expires.Before(time.Now()) {
		cache.Destruct(response.Token.Value)
		return
	}
	return response.User.Claims, err
}
