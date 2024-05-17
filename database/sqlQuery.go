package database

import (
	"fmt"
	"strings"
)

// update
func BuildUpdateQuery(table string, updates map[string]interface{}, id string) (string, []interface{}) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIdx))
		args = append(args, value)
		argIdx++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(setClauses, ", "), argIdx)

	return query, args
}

// update or create
func BuildUpsertQuery(table string, updates map[string]interface{}, id string) (string, []interface{}) {
	columns := []string{}
	values := []string{}
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	for column, value := range updates {
		columns = append(columns, column)
		values = append(values, fmt.Sprintf("$%d", argIdx))
		setClauses = append(setClauses, fmt.Sprintf("%s = EXCLUDED.%s", column, column))
		args = append(args, value)
		argIdx++
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s, id) VALUES (%s, $%d) ON CONFLICT (id) DO UPDATE SET %s",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
		argIdx,
		strings.Join(setClauses, ", "),
	)

	return query, args
}
