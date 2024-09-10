package configure

type Configure interface {
	Provider() ProviderType
	Parse(ptr any) error
}

type ProviderType string

const ConsulProvider ProviderType = "consul"
const FileProvider ProviderType = "file"
