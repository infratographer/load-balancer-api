package graphapi

import (
	"regexp"
	"strings"
	"text/template"

	"go.infratographer.com/x/gidx"
)

// validateGidx validates a gidx.PrefixedID
func validateGidx(gid gidx.PrefixedID) error {
	if _, err := gidx.Parse(gid.String()); err != nil {
		return err
	}

	id := strings.TrimSpace(gid.String())
	if len(id) > 0 {
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(id) {
			return ErrInvalidCharacters
		}
	} else {
		return ErrFieldEmpty
	}

	return nil
}

// sanitizeField sanitizes a field string
func sanitizeField(field string) string {
	s := strings.TrimSpace(field)
	s = template.HTMLEscapeString(s)

	re := regexp.MustCompile(`\r\n|[\r\n\v\f\x{0085}\x{2028}\x{2029}]`)

	s = re.ReplaceAllString(s, " ")

	return s
}
