package sql

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
)

type PostgresAuthorityRepository struct {
	Db *sql.DB
}

func (postgresAuthorityRepository *PostgresAuthorityRepository) Save(ctx context.Context, ae *entity.AuthorityEntity) (*entity.AuthorityEntity, error) {
	err := squirrel.
		Insert("authority").
		Columns("name").
		Values(ae.Name).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(postgresAuthorityRepository.Db).
		ScanContext(ctx, &ae.Id, &ae.Name)
	if err != nil {
		return nil, err
	}
	return ae, err
}

func (postgresAuthorityRepository *PostgresAuthorityRepository) FindAllByUsername(ctx context.Context, username string) ([]*entity.AuthorityEntity, bool) {
	var authorityNames []string
	err := squirrel.
		Select("name").
		From("authority").
		Join("user_authority USING (authority_id)").
		Join("user USING (user_id)").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(postgresAuthorityRepository.Db).
		ScanContext(ctx, authorityNames)
	if err != nil {
		return nil, false
	}
	var res []*entity.AuthorityEntity
	for _, name := range authorityNames {
		_ = append(res, &entity.AuthorityEntity{Name: name})
	}
	return res, true
}
