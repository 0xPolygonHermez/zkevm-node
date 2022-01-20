package runtime

import (
	"math/big"
)

// Params are all the set of params for the chain
type Params struct {
	Forks   *Forks `json:"forks"`
	ChainID int    `json:"chainID"`
}

// Forks specifies when each fork is activated
type Forks struct {
	Homestead      *Fork `json:"homestead,omitempty"`
	Byzantium      *Fork `json:"byzantium,omitempty"`
	Constantinople *Fork `json:"constantinople,omitempty"`
	Petersburg     *Fork `json:"petersburg,omitempty"`
	Istanbul       *Fork `json:"istanbul,omitempty"`
	EIP150         *Fork `json:"EIP150,omitempty"`
	EIP158         *Fork `json:"EIP158,omitempty"`
	EIP155         *Fork `json:"EIP155,omitempty"`
}

func (f *Forks) active(ff *Fork, block uint64) bool {
	if ff == nil {
		return false
	}
	return ff.Active(block)
}

// IsHomestead checks Homestead fork is being used
func (f *Forks) IsHomestead(block uint64) bool {
	return f.active(f.Homestead, block)
}

// IsByzantium checks Byzantium fork is being used
func (f *Forks) IsByzantium(block uint64) bool {
	return f.active(f.Byzantium, block)
}

// IsConstantinople checks Constantinople fork is being used
func (f *Forks) IsConstantinople(block uint64) bool {
	return f.active(f.Constantinople, block)
}

// IsPetersburg checks Petersburg fork is being used
func (f *Forks) IsPetersburg(block uint64) bool {
	return f.active(f.Petersburg, block)
}

// IsEIP150 checks EIP150 fork is being used
func (f *Forks) IsEIP150(block uint64) bool {
	return f.active(f.EIP150, block)
}

// IsEIP158 checks EIP158 fork is being used
func (f *Forks) IsEIP158(block uint64) bool {
	return f.active(f.EIP158, block)
}

// IsEIP155 checks EIP155 fork is being used
func (f *Forks) IsEIP155(block uint64) bool {
	return f.active(f.EIP155, block)
}

// At returns the active fork
func (f *Forks) At(block uint64) ForksInTime {
	return ForksInTime{
		Homestead:      f.active(f.Homestead, block),
		Byzantium:      f.active(f.Byzantium, block),
		Constantinople: f.active(f.Constantinople, block),
		Petersburg:     f.active(f.Petersburg, block),
		Istanbul:       f.active(f.Istanbul, block),
		EIP150:         f.active(f.EIP150, block),
		EIP158:         f.active(f.EIP158, block),
		EIP155:         f.active(f.EIP155, block),
	}
}

// Fork is the current fork
type Fork uint64

// NewFork creates a new fork
func NewFork(n uint64) *Fork {
	f := Fork(n)

	return &f
}

// Active checks if the fork is active
func (f Fork) Active(block uint64) bool {
	return block >= uint64(f)
}

// Int returns the int value for a given fork
func (f Fork) Int() *big.Int {
	return big.NewInt(int64(f))
}

// ForksInTime contains all available forks
type ForksInTime struct {
	Homestead,
	Byzantium,
	Constantinople,
	Petersburg,
	Istanbul,
	EIP150,
	EIP158,
	EIP155 bool
}

// AllForksEnabled creates all available forks
var AllForksEnabled = &Forks{
	Homestead:      NewFork(0),
	EIP150:         NewFork(0),
	EIP155:         NewFork(0),
	EIP158:         NewFork(0),
	Byzantium:      NewFork(0),
	Constantinople: NewFork(0),
	Petersburg:     NewFork(0),
	Istanbul:       NewFork(0),
}
