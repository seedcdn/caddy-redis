package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
)

const defaultUpdateTimer = 5 * time.Minute

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

	updateTimer, err := time.ParseDuration(r.config.Adapter.UpdateTimer)
	if err != nil {
		updateTimer = defaultUpdateTimer
		caddy.Log().Named("adapters.redis.updater").Debug(
			fmt.Sprintf("failed to parse update timer value: %v", err))
	}
	ticker := time.NewTicker(updateTimer)
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
