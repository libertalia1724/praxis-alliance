package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	"golang.org/x/exp/slices"
)

var (
	// Rounder is used to round up small errors due to fix point math
	Rounder = math.LegacyNewDecWithPrec(1, 2)
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		RewardDelayTime:       time.Hour * 24 * 7,
		TakeRateClaimInterval: time.Minute * 5,
		LastTakeRateClaimTime: time.Time{},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

type RewardHistories []RewardHistory

func NewRewardHistories(r []RewardHistory) RewardHistories {
	return r
}

func (r RewardHistories) GetIndexByAlliance(alliance string) (ris RewardHistories) {
	for _, rh := range r {
		// If alliance is empty, it means it was the legacy reward history that does not have a specific alliance
		// Used to handle the old implementation of reward history
		if rh.Alliance == alliance || rh.Alliance == "" {
			ris = append(ris, rh)
		}
	}
	return ris
}

func (r RewardHistories) GetIndexByDenom(denom string, alliance string) (ri *RewardHistory, found bool) {
	idx := slices.IndexFunc(r, func(e RewardHistory) bool {
		return e.Denom == denom && e.Alliance == alliance
	})
	if idx < 0 {
		return &RewardHistory{}, false
	}
	return &r[idx], true
}

func ValidatePositiveDuration(t time.Duration) error {
	if t < 0 {
		return fmt.Errorf("duration must be positive: %d", t)
	}
	return nil
}
