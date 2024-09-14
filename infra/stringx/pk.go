package stringx

import (
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	uuid "github.com/satori/go.uuid"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func PKID(size ...int) string {
	l := 32
	if len(size) > 0 {
		l = size[0]
	}
	id, _ := gonanoid.Generate(chars, l)
	return id
}

func UUID() string {
	u1 := uuid.NewV1()
	return strings.ReplaceAll(u1.String(), "-", "")
}
