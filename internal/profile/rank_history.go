package profile

import "time"

type RankHistory struct {
	Region       string
	RankID       string
	PUUID        string
	Date         time.Time
	Tier         string
	Rank         string
	LeaguePoints int32
}
