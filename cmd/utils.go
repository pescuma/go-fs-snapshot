package cli

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseAddr(addr string) (string, int, error) {
	ip := ""
	port := 0

	if addr == "" {
		return ip, port, nil
	}

	parts := strings.Split(addr, ":")
	if len(parts) > 2 {
		return "", 0, errors.Errorf("invalid address: %v", addr)
	}

	if parts[0] != "" {
		ip = parts[0]
	}
	if len(parts) == 2 && parts[1] != "" {
		var err error

		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, errors.Wrapf(err, "invalid address: %v", addr)
		}
	}

	return ip, port, nil
}
