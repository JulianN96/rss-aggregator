package auth

import (
	"errors"
	"net/http"
	"strings"
)

//GetAPIKey extracts an API Key from headers of HTTP Req
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == ""{
		return "", errors.New("no authentication info found")
	}
	
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("auth header format not valid")
	}
	if vals[0] != "ApiKey"{
		return "", errors.New("auth header format not valid")
	}

	return vals[1], nil
}