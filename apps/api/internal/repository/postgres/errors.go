package postgres

import "strings"

// isDuplicateKey returns true when pgx reports a unique constraint violation (SQLSTATE 23505).
func isDuplicateKey(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}
