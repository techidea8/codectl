package stringx

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func PKID() string {
	u1 := uuid.NewV1()
	return strings.Join(strings.Split(u1.String(), "-"), "")
}
