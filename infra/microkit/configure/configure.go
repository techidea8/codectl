package configure

import "fmt"

type Configure interface {
	Provider() ProviderType
	Parse(ptr any) error
}

type ProviderType string

const ConsulProvider ProviderType = "consul"
const FileProvider ProviderType = "file"

func BuildPath(namespace, env, app, ctype, delim string) string {
	return fmt.Sprintf("%s%s%s%s%s.%s", namespace, delim, env, delim, app, ctype)
}
