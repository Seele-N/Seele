package types

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(heightAdjustment uint64) Minter {
	return Minter{
		HeightAdjustment: heightAdjustment,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(heightAdjustment uint64) Minter {
	return NewMinter(heightAdjustment)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(0)
}

// validate minter
func ValidateMinter(minter Minter) error {
	return nil
}
