package DB

import "database/sql"

func ThrowError() (string, error) {
	return "success", sql.ErrNoRows
}
