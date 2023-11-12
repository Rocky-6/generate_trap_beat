package repository

import (
	"context"

	"github.com/Rocky-6/trap/model"
)

type DBRepository interface {
	Scan(ctx context.Context) ([]model.ChordInfomation, error)
}
