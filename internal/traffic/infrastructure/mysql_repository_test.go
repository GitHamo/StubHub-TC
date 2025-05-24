package infrastructure_test

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/githamo/stubhub-tc/internal/traffic/infrastructure"
	"github.com/stretchr/testify/assert"
)

type MockCrypto struct{}

func (f *MockCrypto) Hash(input string) string {
	return "hashed_path_value"
}

func TestMySQLRepositoryFindByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	uuid := "abc-123"
	path := "/test/path"
	filename := "hashed_path_value"
	content := []byte(`{"message":"hello world"}`)

	mock.ExpectQuery("SELECT path FROM endpoints WHERE id = ?").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"path"}).AddRow(path))

	mock.ExpectQuery("SELECT filename, content FROM stub_contents WHERE filename = ?").
		WithArgs(filename).
		WillReturnRows(sqlmock.NewRows([]string{"filename", "content"}).AddRow(filename, content))

	repo := infrastructure.NewMySQLRepository(db, &MockCrypto{})

	actual, err := repo.FindByUUID(uuid)

	assert.NoError(t, err)
	assert.Equal(t, json.RawMessage(content), actual)
	assert.NoError(t, mock.ExpectationsWereMet())
}
