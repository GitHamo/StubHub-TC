package infrastructure

import (
	"database/sql"
	"encoding/json"

	"github.com/githamo/stubhub-tc/internal/common/encryption"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	db     *sql.DB
	crypto encryption.Hasher
}

func NewMySQLRepository(db *sql.DB, crypto encryption.Hasher) *MySQLRepository {
	return &MySQLRepository{
		db:     db,
		crypto: crypto,
	}
}

func (r *MySQLRepository) FindByUUID(uuid string) (json.RawMessage, error) {
	var (
		path string
		data struct {
			Filename string
			Content  []byte
		}
	)

	firstQuery := `SELECT path FROM endpoints WHERE id = ?`

	err := r.db.QueryRow(firstQuery, uuid).Scan(&path)

	if err != nil {
		return nil, err
	}

	filename := r.crypto.Hash(path)

	secondQuery := `SELECT filename, content FROM stub_contents WHERE filename = ?`

	err = r.db.QueryRow(secondQuery, filename).Scan(
		&data.Filename,
		&data.Content,
	)

	if err != nil {
		return nil, err
	}

	return json.RawMessage(data.Content), nil
}
