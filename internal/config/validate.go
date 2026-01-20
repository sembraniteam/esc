package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var segmentRegex = regexp.MustCompile(`^[a-z0-9]+$`)

const segmentLen = 3

func Match(pattern, name string) bool {
	if pattern == name {
		return true
	}

	if pattern == "*" {
		return true
	}

	if before, ok := strings.CutSuffix(pattern, ".*"); ok {
		prefix := before
		return name == prefix || strings.HasPrefix(name, prefix+".")
	}

	return false
}

func ValidateConnectionName(name string) error {
	if name == "" {
		return errors.New("connection name cannot be empty")
	}

	segments := strings.Split(name, ".")

	if len(segments) > segmentLen {
		return fmt.Errorf(
			"connection name '%s' exceeds maximum of 3 segments",
			name,
		)
	}

	for _, s := range segments {
		if s == "" {
			return fmt.Errorf(
				"connection name '%s' contains empty segment",
				name,
			)
		}
		if !segmentRegex.MatchString(s) {
			return fmt.Errorf(
				"invalid segment '%s' in connection name '%s' (only lowercase alphanumeric allowed)",
				s,
				name,
			)
		}
	}

	return nil
}
