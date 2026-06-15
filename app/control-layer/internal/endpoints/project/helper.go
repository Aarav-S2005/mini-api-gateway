package project

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateGatewayAPIKey() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	key := strings.ToLower(
		base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b),
	)
	return "gw_" + key, nil
}

func GetProjectAndUserID(r *http.Request) (bson.ObjectID, bson.ObjectID, error) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		return bson.ObjectID{0}, bson.ObjectID{}, err
	}
	projectID, err := GetIdFromEndpoint(r, "id")
	if err != nil {
		return bson.ObjectID{}, bson.ObjectID{}, err
	}
	return projectID, userID, nil
}

func ValidateBackendURL(backend string) error {
	u, err := url.Parse(backend)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	if u.User != nil {
		return errors.New("userinfo is not allowed")
	}
	if u.Hostname() == "" {
		return errors.New("host is required")
	}
	if portStr := u.Port(); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			return errors.New("invalid port")
		}
	}
	if u.RawQuery != "" {
		return errors.New("query parameters are not allowed")
	}
	if u.Fragment != "" {
		return errors.New("fragments are not allowed")
	}
	if u.Path != "" && u.Path != "/" {
		return errors.New("backend URL must not contain a path")
	}

	return nil
}

func GetIdFromEndpoint(r *http.Request, key string) (bson.ObjectID, error) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, key))
	if err != nil {
		return bson.ObjectID{}, err
	}
	return id, nil
}
