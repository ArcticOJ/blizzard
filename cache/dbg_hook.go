package cache

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/redis/go-redis/v9"
	"net"
)

type DebugHook struct {
	Name string
}

func (d DebugHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		logger.Blizzard.Debug().Str("network", network).Str("addr", addr).Msgf("dialing '%s'", d.Name)
		return next(ctx, network, addr)
	}
}

func (d DebugHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if len(cmd.String()) < 32 {
			logger.Blizzard.Debug().Str("name", cmd.Name()).Str("cmd", cmd.String()).Err(cmd.Err()).Msgf("running cmd with '%s'", d.Name)
		}
		return next(ctx, cmd)
	}
}

func (d DebugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmd []redis.Cmder) error {
		var cmds []string
		for _, c := range cmd {
			cmds = append(cmds, c.String())
		}
		logger.Blizzard.Debug().Strs("cmds", cmds).Msgf("running pipeline with '%s'", d.Name)
		return next(ctx, cmd)
	}
}
