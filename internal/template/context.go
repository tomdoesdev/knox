package template

import (
	"os"

	"github.com/tomdoesdev/knox/internal/secrets"
	"github.com/tomdoesdev/knox/pkg/errs"
)

type Context struct {
	secretStore secrets.SecretReader
}

func NewContext(secretStore secrets.SecretReader) *Context {
	return &Context{
		secretStore: secretStore,
	}
}

func (ctx *Context) Secret(key string) (string, error) {
	value, err := ctx.secretStore.ReadSecret(key)
	if err != nil {
		return "", errs.Wrap(err, TemplateSecretNotFoundCode, "secret %q not found", key)
	}
	return value, nil
}

func (ctx *Context) Env(key string) string {
	return os.Getenv(key)
}

func (ctx *Context) Default(key, fallback string) string {
	value, err := ctx.secretStore.ReadSecret(key)
	if err != nil {
		return fallback
	}
	return value
}

func (ctx *Context) FuncMap() map[string]interface{} {
	return map[string]interface{}{
		"Secret":  ctx.Secret,
		"Env":     ctx.Env,
		"Default": ctx.Default,
	}
}
