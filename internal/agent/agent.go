package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/jaeyoung050/resource-checker/internal/config"
	"github.com/jaeyoung050/resource-checker/internal/metrics"
)

func Run(ctx context.Context, cfg config.AgentConfig, nodeName string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer client.Close()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			snapshot, err := collectSnapshot(nodeName, cfg.DiskPath)
			if err != nil {
				return err
			}
			payload, err := json.Marshal(snapshot)
			if err != nil {
				return err
			}
			channel := fmt.Sprintf("%s:%s", cfg.ChannelPrefix, nodeName)
			if err := client.Publish(ctx, channel, payload).Err(); err != nil {
				return err
			}
		}
	}
}

func collectSnapshot(nodeName, diskPath string) (metrics.Snapshot, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return metrics.Snapshot{}, err
	}
	memory, err := mem.VirtualMemory()
	if err != nil {
		return metrics.Snapshot{}, err
	}
	diskUsage, err := disk.Usage(diskPath)
	if err != nil {
		return metrics.Snapshot{}, err
	}

	cpuPercent := 0.0
	if len(cpuPercents) > 0 {
		cpuPercent = cpuPercents[0]
	}

	return metrics.Snapshot{
		Node:      nodeName,
		CPU:       cpuPercent,
		Memory:    memory.UsedPercent,
		Disk:      diskUsage.UsedPercent,
		Timestamp: time.Now().Unix(),
	}, nil
}
