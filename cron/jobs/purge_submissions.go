package jobs

import (
	"blizzard/cache/stores"
	"blizzard/db"
	"blizzard/db/models/contest"
	"blizzard/logger"
	"context"
)

func PurgeSubmissions(ctx context.Context) {
	var sub []contest.Submission
	if e := db.Database.NewSelect().Model(&sub).Column("id").Where("result IS ? AND submitted_at < NOW() - INTERVAL '30 MINUTE'", nil).Scan(ctx); e != nil {
		logger.Logger.Error().Err(e).Msg("could not query for staled submissions")
		return
	}
	var toPurge []contest.Submission
	for i := range sub {
		if !stores.Pending.IsPending(ctx, sub[i].ID) {
			toPurge = append(toPurge, sub[i])
		}
	}
	if len(toPurge) > 0 {
		if _, e := db.Database.NewDelete().Model(&toPurge).WherePK().Returning("NULL").Exec(ctx); e != nil {
			logger.Logger.Error().Err(e).Msg("could not purge staled submissions")
			return
		}
	}
}
