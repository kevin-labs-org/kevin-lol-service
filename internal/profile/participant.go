package profile

import "github.com/rank1zen/kevin/internal/store"

// TODO
//SummonerItem item = 10;
//int32 kills = 11;
//int32 deaths = 12;
//int32 assists = 13;
//float kill_participation = 14;
//int32 creep_score = 15;
//float creep_score_per_minute = 16;
//int32 damage_dealt = 17;
//int32 damage_taken = 18;
//int32 damage_delta_counterpart = 19;
//float damage_share = 20;
//int32 gold_earned = 21;
//int32 gold_delta_counterpart = 22;
//float gold_share = 23;
//int32 vision_score = 24;
//int32 control_wards_bought = 25;
//Rank rank = 26;

type Participant struct {
	PUUID            string
	MatchID          string
	TeamID           string
	ParticipantID    string
	ChampionID       string
	ChampionLevel    int
	TeamPosition     string
	SummonerSpellIDs []string
	RuneIDs          []string
}

func newParticipant(participant store.Participant) Participant {
	return Participant{
		PUUID:            participant.PUUID,
		MatchID:          participant.MatchID,
		TeamID:           participant.TeamID,
		ChampionID:       participant.ChampionID,
		ChampionLevel:    participant.ChampionLevel,
		TeamPosition:     participant.TeamPosition,
		SummonerSpellIDs: participant.SummonerIDs,
		RuneIDs:          participant.RuneIDs,
	}
}
