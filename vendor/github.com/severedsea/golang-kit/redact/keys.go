package redact

import "strings"

// Defaults
var (
	// defaultSecretKeys is the secret keys the redact package infers from
	defaultSecretKeys = []string{
		"password",
	}
	// defaultUINKeys is the masked UIN keys the redact package infers from
	defaultUINKeys = []string{}
	// defaultMaskedKeys is the masked keys the redact package infers from
	defaultMaskedKeys = []string{}
)

// Keys is the struct representation for redacted keys
//
// Please note that when accessing the Secrets and Masked fields directly,
// keys should be stripped off of "_" and converted to lowercase.
//
// Sanitising of keys is done before evaluation by stripping "_" and converting
// to lowercase
//
// Example:
//
//	{
//	 "some_key": "this should be redacted if I provide `somekey` OR `some_key` in my lookup table"
//	 "somekey": "this won't be redacted if I provide `some_key` in my lookup table"
//	}
//
// In the event that there are same key names in any of the key groups (i.e. "password" in Secrets and Hashed),
// the redact logic will apply based on this precedence:
// 1. Secrets
// 2. Hashed
// 3. Masked
type Keys struct {
	Secrets []string
	UIN     []string
	Masked  []string
}

// AddSecrets appends the list of keys in the argument to the secret key list
func (k Keys) AddSecrets(values ...string) Keys {
	if k.Secrets == nil {
		k.Secrets = make([]string, 0)
	}
	k.Secrets = append(k.Secrets, values...)

	return k
}

// AddMasked appends the list of keys in the argument to the masked key list
func (k Keys) AddMasked(values ...string) Keys {
	if k.Masked == nil {
		k.Masked = make([]string, 0)
	}
	k.Masked = append(k.Masked, values...)

	return k
}

// IsEmpty evaluates if Keys is empty
func (k Keys) IsEmpty() bool {
	return len(k.Secrets)+len(k.UIN)+len(k.Masked) == 0
}

// DefaultKeys returns the default redacted keys
func DefaultKeys() Keys {
	return Keys{
		Secrets: defaultSecretKeys,
		UIN:     defaultUINKeys,
		Masked:  defaultMaskedKeys,
	}
}

func sanitizeKey(k string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ToLower(k), "_", ""))
}
