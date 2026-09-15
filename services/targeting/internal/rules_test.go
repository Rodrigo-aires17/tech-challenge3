package internal

import "testing"

func TestEvaluateBlockedUserAlwaysFalse(t *testing.T) {
	rule := Rule{FlagID: "flag-1", BlockedUserIDs: []string{"user-1"}, RolloutPercentage: 100}
	if Evaluate(rule, "user-1") {
		t.Fatal("blocked user should never match")
	}
}

func TestEvaluateAllowedUserAlwaysTrue(t *testing.T) {
	rule := Rule{FlagID: "flag-1", AllowedUserIDs: []string{"user-2"}, RolloutPercentage: 0}
	if !Evaluate(rule, "user-2") {
		t.Fatal("allow-listed user should always match")
	}
}

func TestEvaluateRolloutBoundaries(t *testing.T) {
	rule0 := Rule{FlagID: "flag-1", RolloutPercentage: 0}
	if Evaluate(rule0, "any-user") {
		t.Fatal("0% rollout should never match")
	}

	rule100 := Rule{FlagID: "flag-1", RolloutPercentage: 100}
	if !Evaluate(rule100, "any-user") {
		t.Fatal("100% rollout should always match")
	}
}

func TestEvaluateIsDeterministic(t *testing.T) {
	rule := Rule{FlagID: "flag-1", RolloutPercentage: 50}
	first := Evaluate(rule, "stable-user")
	for i := 0; i < 20; i++ {
		if Evaluate(rule, "stable-user") != first {
			t.Fatal("evaluation should be deterministic for the same user/flag")
		}
	}
}
