package profile

import "time"

type MatchHistory struct {
	Region                 string
	MatchID                string
	PUUID                  string
	Date                   time.Time
	Duration               time.Duration
	Version                string
	WinnerID               int
	TeamID                 int
	ParticipantID          int
	ChampionID             int
	ChampionLevel          int
	TeamPosition           string
	SummonerSpell          []string
	RunePage               []string
	Item                   []string
	Kills                  int
	Deaths                 int
	Assists                int
	KillParticipation      float32
	CreepScore             int
	CreepScorePerMinute    float32
	DamageDealt            int
	DamageTaken            int
	DamageDeltaCounterpart int
	DamageShare            float32
	GoldEarned             int
	GoldDeltaCounterpart   int
	GoldShare              float32
	VisionScore            int
	ControlWardsBought     int
	Rank                   *Rank
}
