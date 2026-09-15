package internal

import (
	"hash/fnv"
	"slices"
)

// Rule descreve como uma flag deve ser direcionada para um usuário:
// lista de usuários permitidos/bloqueados e/ou rollout percentual.
type Rule struct {
	FlagID            string   `json:"flag_id"`
	AllowedUserIDs    []string `json:"allowed_user_ids,omitempty"`
	BlockedUserIDs    []string `json:"blocked_user_ids,omitempty"`
	RolloutPercentage int      `json:"rollout_percentage"` // 0-100
}

// Evaluate decide se o usuário está dentro do público-alvo da regra.
// O bucket percentual é determinístico (mesmo usuário sempre cai no mesmo bucket
// para uma dada flag), calculado via hash FNV-1a de flagID+userID.
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
