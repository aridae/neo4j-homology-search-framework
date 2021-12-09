package dataaccess

type Repository interface {
	MergePrecedingKMers(kPlus1Mer []byte) error
}
