package profile

import (
	"context"

	"github.com/rank1zen/kevin/internal/riot"
)

type RiotClient interface {
	GetSummoner(ctx context.Context, region string, puuid string) (*riot.Summoner, error)
	GetLeagueEntriesByPUUID(ctx context.Context, region string, puuid string) (riot.LeagueList, error)
	GetAccountByRiotID(ctx context.Context, region string, name string, tag string) (*riot.Account, error)
	GetAccountByPUUID(ctx context.Context, region string, puuid string) (*riot.Account, error)
}
