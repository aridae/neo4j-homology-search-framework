package dataaccess

type Repository interface {
	MergekPlus1Mer(kPlus1Mer []byte) error
}
