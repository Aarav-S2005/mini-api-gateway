package plugin

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
)

type IpFilterPlugin struct{}

type IpFilterConfig struct {
	BlackListedIps []string `json:"black_listed_ips"`
	WhiteListedIps []string `json:"white_listed_ips"`
}

func (i *IpFilterPlugin) Name() string {
	return "ip-filter"
}

func (i *IpFilterPlugin) Validate(config map[string]interface{}) error {
	var cfg IpFilterConfig
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}
	for _, ip := range cfg.BlackListedIps {
		parsedIp := net.ParseIP(ip)
		if parsedIp == nil {
			return errors.New("invalid ip address: " + ip)
		}
	}
	for _, ip := range cfg.WhiteListedIps {
		parsedIp := net.ParseIP(ip)
		if parsedIp == nil {
			return errors.New("invalid ip address: " + ip)
		}
	}
	return nil
}

func (i *IpFilterPlugin) CreateMiddleware(config map[string]interface{}) (MiddlewareFunc, error) {
	var cfg IpFilterConfig
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	blacklist := make(map[string]struct{})
	whitelist := make(map[string]struct{})
	for _, ip := range cfg.BlackListedIps {
		blacklist[ip] = struct{}{}
	}
	for _, ip := range cfg.WhiteListedIps {
		whitelist[ip] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, err := getClientIP(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ip := net.ParseIP(host)
			if ip == nil {
				http.Error(w, "invalid ip address", http.StatusBadRequest)
				return
			}
			clientIP := ip.String()
			if _, ok := blacklist[clientIP]; ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if _, ok := whitelist[clientIP]; len(whitelist) > 0 || !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}, nil
}

func getClientIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0], nil
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}
