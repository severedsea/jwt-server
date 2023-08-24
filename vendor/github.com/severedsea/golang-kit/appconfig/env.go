package appconfig

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// LoadFromEnv loads all the app config from env into the context
func LoadFromEnv(ctx context.Context) (context.Context, error) {
	value, err := fromEnv()
	if err != nil {
		return ctx, err
	}

	return Set(ctx, value), nil
}

func fromEnv() (Config, error) {
	value := map[string]interface{}{}

	for _, it := range os.Environ() {
		parts := strings.Split(it, "=") // split by = sign
		k := parts[0]
		v := parts[1]
		if len(parts) > 2 {
			for _, it := range parts[2:] {
				v += ("=" + it)
			}
		}

		switch v {
		case "true":
			value[k] = true
		case "false":
			value[k] = false
		default:
			intV, err := strconv.Atoi(v)
			if err == nil {
				value[k] = intV
			} else {
				value[k] = v
			}
		}
	}

	// Convert to Config
	b, err := json.Marshal(value)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
