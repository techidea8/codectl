package stringx

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func PKID() string {
	id, _ := gonanoid.New()
	return id
}
