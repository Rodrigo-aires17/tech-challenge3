package internal

import (
	"hash/fnv"
	"slices"
)

// Rule replica a lógica de segmentação do serviço `targeting` (percentual
// determinístico + allow/block lists) para decidir o valor final da flag.
type Rule struct {
	FlagID            string   `json:"flag_id"`
	AllowedUserIDs    []string `json:"allowed_user_ids,omitempty"`
	BlockedUserIDs    []string `json:"blocked_user_ids,omitempty"`
	RolloutPercentage int      `json:"rollout_percentage"`
}

func Evaluate(rule Rule, userID string) bool {
	if slices.Contains(rule.BlockedUserIDs, userID) {
		return false
	}
	if slices.Contains(rule.AllowedUserIDs, userID) {
		return true
	}
	if rule.RolloutPercentage <= 0 {
		return false
	}
	if rule.RolloutPercentage >= 100 {
		return true
	}
	return bucketFor(rule.FlagID, userID) < rule.RolloutPercentage
}

func bucketFor(flagID, userID string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(flagID + ":" + userID))
	return int(h.Sum32() % 100)
}
