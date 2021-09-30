package ibm

import (
	"context"
	"errors"
	"fmt"
	gohttp "net/http"
	"os"
	"strings"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/strfmt"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func configApiKey(_ context.Context, d *plugin.QueryData) (string, error) {

	// First, prefer the most specific env var
	apiKey := os.Getenv("IBMCLOUD_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	// Second, fall back to the common short form
	apiKey = os.Getenv("IC_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	// Then, use config
	ibmConfig := GetConfig(d.Connection)
	if ibmConfig.APIKey != nil {
		apiKey = *ibmConfig.APIKey
	}
	if apiKey != "" {
		return apiKey, nil
	}

	// No key, cannot proceed
	return "", errors.New("api_key must be configured")
}

func connect(ctx context.Context, d *plugin.QueryData) (*session.Session, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*session.Session), nil
	}

	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}

	conf := &bluemix.Config{
		BluemixAPIKey: apiKey,
	}

	conn, err := session.New(conf)
	if err != nil {
		return nil, err
	}

	err = authenticateAPIKey(conn)
	if err != nil {
		return nil, err
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn, nil
}

//UserConfig ...
type UserConfig struct {
	userID      string
	userEmail   string
	userAccount string
	cloudName   string `default:"bluemix"`
	cloudType   string `default:"public"`
	generation  int    `default:"2"`
}

func fetchUserDetails(sess *session.Session, generation int) (*UserConfig, error) {
	config := sess.Config
	user := UserConfig{}
	var bluemixToken string
	if strings.HasPrefix(config.IAMAccessToken, "Bearer") {
		bluemixToken = config.IAMAccessToken[7:len(config.IAMAccessToken)]
	} else {
		bluemixToken = config.IAMAccessToken
	}
	token, err := jwt.Parse(bluemixToken, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	//TODO validate with key
	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		return &user, err
	}
	claims := token.Claims.(jwt.MapClaims)
	if email, ok := claims["email"]; ok {
		user.userEmail = email.(string)
	}
	user.userID = claims["id"].(string)
	user.userAccount = claims["account"].(map[string]interface{})["bss"].(string)
	iss := claims["iss"].(string)
	if strings.Contains(iss, "https://iam.cloud.ibm.com") {
		user.cloudName = "bluemix"
	} else {
		user.cloudName = "staging"
	}
	user.cloudType = "public"
	user.generation = generation
	return &user, nil
}

func authenticateAPIKey(sess *session.Session) error {
	config := sess.Config
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return err
	}
	return tokenRefresher.AuthenticateAPIKey(config.BluemixAPIKey)
}

func resourceInterfaceDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "Title of the resource."
	}
	return ""
}

// Transform to ensure a string array. It's better to use transform.EnsureStringArray, but
// as of v0.2.6 it does not support a *string argument.
func ensureStringArray(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value != nil {
		switch v := d.Value.(type) {
		case []string:
			return v, nil
		case *string:
			return []string{*v}, nil
		case string:
			return []string{v}, nil
		default:
			str := fmt.Sprintf("%v", d.Value)
			return []string{str}, nil
		}
	}
	return nil, nil
}

func ensureTimestamp(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	t := d.Value.(*strfmt.DateTime)
	plugin.Logger(ctx).Warn("ensureTimestamp", "d.Value", t)
	return t.String(), nil
}

func crnToAccountID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	crn := d.Value.(*string)
	crnParts := strings.Split(*crn, ":")
	accountIDPart := crnParts[6]
	if accountIDPart == "" {
		return "", nil
	}
	aParts := strings.Split(accountIDPart, "/")
	accountID := aParts[1]
	if accountID == "" {
		return "", nil
	}
	return accountID, nil
}
