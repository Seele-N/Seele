package params

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (

	// Day time of day
	Day = 24 * time.Hour
	// DefaultBondDenom Default Bond Denom
	DefaultBondDenom = "snp"
	// DefaultMintDenom Default Mint Denom
	DefaultMintDenom = "seele"
)

// staking
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	DefaultUnbondingTime = 21 * Day

	// DefaultMaxValidators Default maximum number of bonded validators
	DefaultMaxValidators uint32 = 21
)

// Crisis
var (
	// DefaultCrisisConstantFee Default Crisis Constant Fee
	DefaultCrisisConstantFee = sdk.NewInt(1000)
)

// gov
var (
	// DefaultGovMinDepositAmount Default Gov Min Deposit Amount
	DefaultGovMinDepositAmount = sdk.NewInt(10000000)
)

// gov
const (
	// DefaultPeriod Default period for deposits & voting
	DefaultPeriod = 14 * Day
	// VotingPeriod Voting Period
	VotingPeriod = 7 * Day
)

// mint
const (
	BlocksPerYear = 6311520
)
