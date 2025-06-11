package customstaking

import "time"

type Params struct {
	ValidatorUnbondingTime time.Duration
	DelegatorUnbondingTime time.Duration
}

// Example default values.
func DefaultParams() Params {
	return Params{
		ValidatorUnbondingTime: time.Hour * 24 * 30, // 30 days
		DelegatorUnbondingTime: time.Hour * 24 * 7,  // 7 days
	}
}
