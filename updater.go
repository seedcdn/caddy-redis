package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
)

const defaultUpdateInterval = 5 * time.Minute

func (r *Redis) updateConfiguration() {
	if r.config.Adapter.SubscribeUpdates != "" {
		subscription := r.client.Subscribe(context.TODO(), r.config.Adapter.SubscribeUpdates)
		for {
			_, _ = subscription.ReceiveMessage(context.TODO())
			config, err := r.generateConfiguration()
			if err != nil {
				caddy.Log().Named("adapters.redis.updater").Debug(
					fmt.Sprintf("failed to get config: %v", err))
				return
			}
			if err = caddy.Load(config, false); err != nil {
				caddy.Log().Named("adapters.redis.updater").Debug(
					fmt.Sprintf("failed to load config: %v", err))
			}
		}
	}

	updateInterval, err := time.ParseDuration(r.config.Adapter.UpdateInterval)
	if err != nil {
		updateInterval = defaultUpdateInterval
		caddy.Log().Named("adapters.redis.updater").Debug(
			fmt.Sprintf("failed to parse update interval value: %v", err))
	}
	ticker := time.NewTicker(updateInterval)
	go func() {
		for range ticker.C {
			config, err := r.generateConfiguration()
			if err != nil {
				caddy.Log().Named("adapters.redis.updater").Debug(
					fmt.Sprintf("failed to get config: %v", err))
				return
			}
			if err = caddy.Load(config, false); err != nil {
				caddy.Log().Named("adapters.redis.updater").Debug(
					fmt.Sprintf("failed to load config: %v", err))
			}
		}
	}()
}
