package secrets

type SecretReader interface {
	Read(key string) (value string, err error)
}

type SecretWriter interface {
	Write(key, value string) error
}
type SecretReaderWriter interface {
	SecretReader
	SecretWriter
}
