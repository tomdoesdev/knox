package secrets

type SecretStore interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
}
