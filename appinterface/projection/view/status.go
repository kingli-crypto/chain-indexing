package view

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-com/chainindex/appinterface/rdb"
)

type Status struct {
	rdb *rdb.Handle
}

func NewStatus(handle *rdb.Handle) *Status {
	return &Status{
		handle,
	}
}

func (view *Status) Insert(statusid string, statusvalue string) error {
	var err error

	var sql string
	suffix := sq.Expr("on conflict(key) do update set value=?", 12)
	sql, _, err = view.rdb.StmtBuilder.Insert(
		"view_status",
	).Columns(
		"key",
		"value",
	//).Values("?", "?").Suffix(fmt.Sprintf("on conflict(key) do update set value='%s'", statusvalue)).ToSql()
	).Values("?", "?").SuffixExpr(suffix).ToSql()

	if err != nil {
		return fmt.Errorf("error building blocks insertion sql: %v: %w", err, rdb.ErrBuildSQLStmt)
	}
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

	fmt.Printf("SQL=%v\n", sql)
	result, err := view.rdb.Exec(sql, statusid, statusvalue, statusvalue)
	if err != nil {
		panic(err)
		return fmt.Errorf("error inserting view_status into the table: %v: %w", err, rdb.ErrWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error inserting view_status into the table: no rows inserted: %w", rdb.ErrWrite)
	}

	return nil
}
func (view *Status) FindBy(statusid string) (string, error) {
	var err error

	selectStmtBuilder := view.rdb.StmtBuilder.Select(
		"value",
	).From("view_status").Where("key = ?", statusid)

	sql, sqlArgs, err := selectStmtBuilder.ToSql()
	fmt.Printf("sql=%s sql args=%s", sql, sqlArgs)

	if err != nil {
		return "", fmt.Errorf("error building status selection sql: %v: %w", err, rdb.ErrPrepare)
	}

	var found string
	if err = view.rdb.QueryRow(sql, sqlArgs...).Scan(
		&found,
	); err != nil {
		if err == rdb.ErrNoRows {
			return "", rdb.ErrNoRows
		}
		return "", fmt.Errorf("error scanning block row: %v: %w", err, rdb.ErrQuery)
	}

	return found, nil
}
