package infrastructure

import (
	"database/sql"
	"encoding/json"

	"github.com/githamo/stubhub-tc/internal/common/encryption"
	"github.com/githamo/stubhub-tc/internal/traffic/domain"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	db     *sql.DB
	crypto *encryption.Helper
}

func NewMySQLRepository(db *sql.DB, crypto *encryption.Helper) *MySQLRepository {
	return &MySQLRepository{
		db:     db,
		crypto: crypto,
	}
}

func (r *MySQLRepository) FindByUUID(uuid string) (*domain.TrafficContentData, error) {
	var endpoint domain.TrafficEndpointData

	firstQuery := `SELECT path FROM endpoints WHERE id = ?`

	err := r.db.QueryRow(firstQuery, uuid).Scan(&endpoint.Path)

	if err != nil {
		return nil, err
	}

	hashedFilename := r.crypto.Hash(endpoint.Path)

	var data domain.TrafficContentData
	var content []byte

	secondQuery := `SELECT filename, content FROM stub_contents WHERE filename = ?`

	err = r.db.QueryRow(secondQuery, hashedFilename).Scan(
		&data.Filename,
		&content,
	)

	if err != nil {
		return nil, err
	}

	data.Content = json.RawMessage(content)

	return &data, nil
}
