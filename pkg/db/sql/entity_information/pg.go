package entity_information

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/Masterminds/squirrel"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"log"
)

const (
	count = "count(*) as caunt"
)

func IsNew(ctx context.Context, tableName, idColumnName string, idValue interface{}, tx *sql.Tx) (bool, errors.Errors, *sql.Tx) {
	if idValue == nil {
		return true,
			nil,
			tx
	}
	idValue = ResolveIdType(idValue)
	var countOfRowsWithSpecifiedIdValue int
	err := squirrel.
		Select(count).
		From(tableName).
		Where(squirrel.Eq{idColumnName: idValue}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &countOfRowsWithSpecifiedIdValue)
	if err != nil || countOfRowsWithSpecifiedIdValue != 0 {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(err, tx)
		return false,
			errorz,
			tx
	}
	return true,
		nil,
		tx
}

func ResolveIdType(idValue interface{}) interface{} {
	switch idValue.(type) {
	case uuid.UUID:
		return idValue.(uuid.UUID).UUID.String()
	}
	log.Println("unresolvable type: %T", idValue)
	panic(idValue)
}
