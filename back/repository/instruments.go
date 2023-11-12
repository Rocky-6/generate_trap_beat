package repository

import "context"

type InstrumentsRepository interface {
	MakeSMF(ctx context.Context) ([]byte, error)
}
