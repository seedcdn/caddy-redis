package redis

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/redis/go-redis/v9"
)

var _ caddyconfig.Adapter = (*Redis)(nil)

func init() {
	caddyconfig.RegisterAdapter("redis", &Redis{})
}

func (r *Redis) Adapt(raw []byte, _ map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	r.config = &Config{
		Adapter: AdapterConfig{
			Address:     "localhost:6379",
			UpdateTimer: "5m", Prefix: "caddy.",
		},
	}

	err := json.Unmarshal(raw, &r.config)
	if err != nil {
		return nil, nil, err
	}

	caddy.Log().Named("adapters.redis.log").Info(string(raw))
	caddy.Log().Named("adapters.redis.log").Info(r.config.Adapter.Address)
	caddy.Log().Named("adapters.redis.log").Info(r.config.Adapter.Password)
	r.client = redis.NewClient(&redis.Options{
		DB:       r.config.Adapter.Database,
		Addr:     r.config.Adapter.Address,
		Password: r.config.Adapter.Password,
	})

	var config []byte
	config, err = r.generateConfiguration()
	if err != nil {
		return nil, nil, err
	}

	go r.updateConfiguration()
	return config, nil, err
}
