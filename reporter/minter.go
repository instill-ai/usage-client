package reporter

// Minter a interface that creates Hashcash stamp
type Minter interface {
	Mint(string) (string, error)
}
