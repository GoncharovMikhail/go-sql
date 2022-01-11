package authority

import (
	"context"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
)

type AuthorityRepository interface {
	Save(ctx context.Context, ae *entity.AuthorityEntity) (*entity.AuthorityEntity, error)
	FindAllByUsername(ctx context.Context, username string) ([]*entity.AuthorityEntity, bool)
}
