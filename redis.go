package redis

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	config *Config
	client *redis.Client
}

func (r *Redis) loadServers() map[string]*caddyhttp.Server {
	var servers = make(map[string]*caddyhttp.Server)
	keys, err := r.client.Keys(context.TODO(), r.config.Adapter.Prefix+"*").Result()
	if err != nil {
		return nil
	}
	for _, key := range keys {
		srv := strings.TrimPrefix(key, r.config.Adapter.Prefix)
		servers[srv] = &caddyhttp.Server{}
		raw := r.client.Get(context.TODO(), key).Val()
		_ = json.Unmarshal([]byte(raw), servers[srv])
	}

	return servers
}

func (r *Redis) generateConfiguration() ([]byte, error) {
	var config = caddy.Config{}
	admin, _ := json.Marshal(r.config.Admin)
	logging, _ := json.Marshal(r.config.Logging)
	storage, _ := json.Marshal(r.config.Storage)
	_ = json.Unmarshal(admin, &config.Admin)
	_ = json.Unmarshal(logging, &config.Logging)
	_ = json.Unmarshal(storage, &config.StorageRaw)

	var storageRaw json.RawMessage
	_ = json.Unmarshal(storageRaw, &storageRaw)
	config.StorageRaw = storageRaw

	apps := caddyhttp.App{}
	apps.Servers = r.loadServers()

	var warnings []caddyconfig.Warning
	config.AppsRaw = make(caddy.ModuleMap)
	for module, app := range r.config.Apps {
		config.AppsRaw[module] = caddyconfig.JSON(&app, &warnings)
	}
	config.AppsRaw["http"] = caddyconfig.JSON(&apps, &warnings)
	for _, warning := range warnings {
		caddy.Log().Named("adapters.redis.loader").Warn(warning.String())
	}

	return json.Marshal(config)
}
