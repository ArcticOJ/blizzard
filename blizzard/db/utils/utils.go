package utils

import (
	"fmt"
	"github.com/uptrace/bun"
)

func Cascade(query *bun.CreateTableQuery, col, table, tcol string) error {
	query.ForeignKey(fmt.Sprintf(`("%s") REFERENCES "%s" ("%s") ON DELETE CASCADE`, col, table, tcol))
	return nil
}
