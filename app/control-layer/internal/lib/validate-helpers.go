package lib

import (
	"errors"
	"net/url"
	"strconv"
)

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
