package env

import (
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = time.StampMicro
}

func osGet(key string) (value string, err error) {
	if value = os.Getenv(key); value == "" {
		return "", fmt.Errorf("env variable " + key + " is empty")
	}
	return value, nil
}

func Get(key string) (value string, err error) {
	return osGet(key)
}

func GetOrEmpty(key string) (value string) {
	if value, err := osGet(key); err != nil {
		return ""
	} else {
		return value
	}
}

func MustGet(key string) (value string) {
	if value, err := osGet(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func osGetBool(key string) (bool, error) {
	if value, err := osGet(key); err != nil {
		return false, err
	} else {
		return value == "1", nil
	}
}

func MustGetBool(key string) (value bool) {
	if value, err := osGetBool(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func osGetBigInt(key string) (*big.Int, error) {
	if value, err := osGet(key); err != nil {
		return big.NewInt(0), err
	} else {
		if res, ok := big.NewInt(0).SetString(value, 10); !ok {
			return big.NewInt(0), fmt.Errorf("env variable " + key + " is not a valid integer")
		} else {
			return res, nil
		}
	}
}

func MustGetBigInt(key string) (value *big.Int) {
	if value, err := osGetBigInt(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func MustGetUrl(key string) (value *url.URL) {
	if value, err := osGetUrl(key); err != nil {
		panic(err)
	} else {
		return value
	}
}
func osGetUrl(key string) (value *url.URL, err error) {
	if value, err = url.Parse(os.Getenv(key)); err != nil {
		return nil, err
	}
	return value, nil
}

func MustGetStrings(key string) []string {
	return strings.Split(MustGet(key), "|")
}
