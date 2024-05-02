package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilya372317/pass-keeper/internal/server/adapter/pgsqlrepo/datarepo"
	"github.com/ilya372317/pass-keeper/internal/server/adapter/pgsqlrepo/keyrepo"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ilya372317/pass-keeper/internal/server/adapter/pgsqlrepo/userrepo"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db       *sqlx.DB
	userRepo *userrepo.Repository
	keyRepo  *keyrepo.Repository
	dataRepo *datarepo.Repository
)

type keysFields struct {
	Key   string `db:"key"`
	Nonce string `db:"nonce"`
}

func TestMain(m *testing.M) {
	database, pool, resource, err := makeTestConnection("../../db/migrations")
	if err != nil {
		log.Fatal(err)
		return
	}
	db = database
	userRepo = userrepo.New(database)
	keyRepo = keyrepo.New(database)
	dataRepo = datarepo.New(database)
	m.Run()
	if err = closeTestConnection(database, pool, resource); err != nil {
		log.Fatal(err)
		return
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	type want struct {
		err  bool
		user *domain.User
	}
	tests := []struct {
		name   string
		arg    string
		want   want
		fields []domain.User
	}{
		{
			name: "success found case with one user in storage",
			arg:  "email1",
			want: want{
				err: false,
				user: &domain.User{
					Email:          "email1",
					HashedPassword: "123",
					Salt:           "salt",
				},
			},
			fields: []domain.User{
				{
					Email:          "email1",
					HashedPassword: "123",
					Salt:           "salt",
				},
			},
		},
		{
			name: "success found case with multiply users in storage",
			arg:  "email1",
			want: want{
				err: false,
				user: &domain.User{
					Email:          "email1",
					HashedPassword: "pass1",
					Salt:           "salt1",
				},
			},
			fields: []domain.User{
				{
					Email:          "email1",
					HashedPassword: "pass1",
					Salt:           "salt1",
				},
				{
					Email:          "email2",
					HashedPassword: "pass2",
					Salt:           "salt2",
				},
				{
					Email:          "email3",
					HashedPassword: "pass3",
					Salt:           "salt3",
				},
			},
		},
		{
			name: "user not found case",
			arg:  "email1",
			want: want{
				err: true,
			},
			fields: []domain.User{
				{
					Email:          "email2",
					HashedPassword: "pass2",
					Salt:           "salt2",
				},
			},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillUsersTable(t, tt.fields)

			got, err := userRepo.GetUserByEmail(ctx, tt.arg)

			if tt.want.err {
				require.Error(t, err)
				require.ErrorIs(t, err, sql.ErrNoRows)
				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.user.Email, got.Email)
			assert.Equal(t, tt.want.user.HashedPassword, got.HashedPassword)
			assert.Equal(t, tt.want.user.Salt, got.Salt)
			clearUsersTable(t)
		})
	}
}

func TestUserRepository_SaveUser(t *testing.T) {
	tests := []struct {
		name     string
		argument *domain.User
		wantErr  bool
	}{
		{
			name: "save user success case",
			argument: &domain.User{
				Email:          "email",
				HashedPassword: "password",
				Salt:           "salt",
			},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userRepo.SaveUser(ctx, tt.argument)
			require.NoError(t, err)
			var savedUser domain.User

			err = db.Get(&savedUser, "SELECT * FROM users WHERE email = $1", tt.argument.Email)

			require.NoError(t, err)
			assert.Equal(t, tt.argument.Email, savedUser.Email)
			assert.Equal(t, tt.argument.Salt, savedUser.Salt)
			assert.Equal(t, tt.argument.HashedPassword, savedUser.HashedPassword)

			clearUsersTable(t)
		})
	}
}

func TestUserRepository_HasUser(t *testing.T) {
	tests := []struct {
		name     string
		argument string
		data     []domain.User
		want     bool
	}{
		{
			name:     "success case with empty storage",
			argument: "email",
			want:     false,
		},
		{
			name:     "success case with filled storage",
			argument: "email1",
			data: []domain.User{
				{
					Email:          "email2",
					HashedPassword: "123",
					Salt:           "salt",
				},
			},
			want: false,
		},
		{
			name:     "has user case",
			argument: "email",
			data: []domain.User{
				{
					Email:          "email",
					HashedPassword: "123",
					Salt:           "salt",
				},
			},
			want: true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillUsersTable(t, tt.data)
			got, err := userRepo.HasUser(ctx, tt.argument)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			clearUsersTable(t)
		})
	}
}

func TestKeyRepository_SaveKey(t *testing.T) {
	const expectedKeyNumber = 1
	tests := []struct {
		name     string
		fields   []keysFields
		argument *domain.Keys
	}{
		{
			name: "success save case with filled storage",
			fields: []keysFields{
				{
					Key:   "key1",
					Nonce: "nonce1",
				},
				{
					Key:   "key2",
					Nonce: "nonce2",
				},
				{
					Key:   "key3",
					Nonce: "nonce3",
				},
				{
					Key:   "key4",
					Nonce: "nonce4",
				},
			},
			argument: &domain.Keys{
				Key:       "key5",
				Nonce:     "nonce5",
				IsCurrent: true,
			},
		},
		{
			name:   "success case with empty storage",
			fields: nil,
			argument: &domain.Keys{
				Key:   "key1",
				Nonce: "nonce1",
			},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillKeysTable(t, tt.fields)
			err := keyRepo.SaveKey(ctx, tt.argument)
			require.NoError(t, err)
			var currentKeyNumber int
			err = db.Get(&currentKeyNumber, "SELECT COUNT(*) FROM keys WHERE is_current = true")
			require.NoError(t, err)
			assert.Equal(t, expectedKeyNumber, currentKeyNumber)
			clearKeysTable(t)
		})
	}
}

func TestKeyRepository_GetKey(t *testing.T) {
	type want struct {
		err   bool
		key   string
		nonce string
	}
	tests := []struct {
		name string
		data []keysFields
		want want
	}{
		{
			name: "get Key success case",
			data: []keysFields{
				{
					Key:   "key1",
					Nonce: "nonce1",
				},
			},
			want: want{
				err:   false,
				key:   "key1",
				nonce: "nonce1",
			},
		},
		{
			name: "key not found case",
			data: nil,
			want: want{
				err: true,
			},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillKeysTable(t, tt.data)

			got, err := keyRepo.GetKey(ctx)

			if tt.want.err {
				require.Error(t, err)
				require.ErrorIs(t, err, sql.ErrNoRows)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.key, got.Key)
			assert.Equal(t, tt.want.nonce, got.Nonce)
			clearKeysTable(t)
		})
	}
}

func TestDataRepository_SaveData(t *testing.T) {
	user := getOrCreateUser(t)
	type want struct {
		result domain.Data
		err    bool
	}
	tests := []struct {
		name string
		arg  domain.Data
		want want
	}{
		{
			name: "success case with empty storage",
			arg: domain.Data{
				Payload:        `payload`,
				Metadata:       `{"url":"test"}`,
				PayloadNonce:   "123",
				CryptoKeyNonce: "123",
				CryptoKey:      "123",
				UserID:         int(user.ID),
				Kind:           domain.KindLoginPass,
			},
			want: want{
				result: domain.Data{
					Payload:        `payload`,
					Metadata:       `{"url":"test"}`,
					PayloadNonce:   "123",
					CryptoKeyNonce: "123",
					CryptoKey:      "123",
					UserID:         int(user.ID),
					Kind:           domain.KindLoginPass,
				},
				err: false,
			},
		},
		{
			name: "invalid user id case",
			arg: domain.Data{
				Payload:        `{"login":"password"}`,
				Metadata:       `{"url":"test"}`,
				PayloadNonce:   "123",
				CryptoKeyNonce: "123",
				CryptoKey:      "123",
				UserID:         0,
				Kind:           domain.KindLoginPass,
			},
			want: want{
				result: domain.Data{},
				err:    true,
			},
		},
		{
			name: "invalid metadata case",
			arg: domain.Data{
				Payload:        `{"login":"password"}`,
				Metadata:       `invalid-metadata`,
				PayloadNonce:   "123",
				CryptoKeyNonce: "123",
				CryptoKey:      "123",
				UserID:         int(user.ID),
				Kind:           domain.KindLoginPass,
			},
			want: want{
				result: domain.Data{},
				err:    true,
			},
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dataRepo.SaveData(ctx, tt.arg)
			if tt.want.err {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			got := getLastInsertData(t)

			assert.Equal(t, tt.want.result.Payload, got.Payload)
			assert.Equal(t, tt.want.result.CryptoKeyNonce, got.CryptoKeyNonce)
			assert.Equal(t, tt.want.result.Metadata, got.Metadata)
			assert.Equal(t, tt.want.result.CryptoKey, got.CryptoKey)
			assert.Equal(t, tt.want.result.UserID, got.UserID)
			assert.Equal(t, tt.want.result.Kind, got.Kind)
			clearDataRecordsTable(t)
		})
	}
}

func getOrCreateUser(t *testing.T) domain.User {
	t.Helper()
	var user domain.User
	err := db.Get(&user, "SELECT * FROM users WHERE id = (SELECT MAX(id) FROM users)")
	if errors.Is(err, sql.ErrNoRows) {
		fillUsersTable(t, []domain.User{{
			Email:          "1@gmail.com",
			HashedPassword: "123",
			Salt:           "123",
		}})
		err = db.Get(&user, "SELECT * FROM users WHERE id = (SELECT MAX(id) FROM users)")
		require.NoError(t, err)
	}

	return user
}

func getLastInsertData(t *testing.T) *domain.Data {
	t.Helper()
	var result domain.Data
	err := db.Get(&result, "SELECT * FROM data_records WHERE id = (SELECT max(id) FROM data_records)")
	require.NoError(t, err)

	return &result
}

func clearDataRecordsTable(t *testing.T) {
	t.Helper()
	_, err := db.Exec("DELETE FROM data_records WHERE  id > 0")
	require.NoError(t, err)
}

func fillKeysTable(t *testing.T, data []keysFields) {
	t.Helper()
	for _, key := range data {
		_, err := db.NamedExec("INSERT INTO keys (key, nonce) VALUES (:key, :nonce)", key)
		require.NoError(t, err)
	}
}

func clearKeysTable(t *testing.T) {
	t.Helper()
	_, err := db.Exec("DELETE FROM keys WHERE id > 0")
	require.NoError(t, err)
}

func fillUsersTable(t *testing.T, users []domain.User) {
	t.Helper()
	if len(users) == 0 {
		return
	}
	_, err := db.NamedExec(
		"INSERT INTO users (email, hashed_password, salt) VALUES (:email, :hashed_password, :salt)",
		users,
	)
	require.NoError(t, err)
}

func clearUsersTable(t *testing.T) {
	t.Helper()
	_, err := db.Exec("DELETE FROM users WHERE id > 0")
	require.NoError(t, err)
}

func makeTestConnection(migrationPath string) (*sqlx.DB, *dockertest.Pool, *dockertest.Resource, error) {
	var db *sqlx.DB

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	resource, err := pool.
		Run("postgres", "15", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=pass_test"})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	port := resource.GetPort("5432/tcp")
	connectionString := fmt.Sprintf(
		"host=localhost port=%s user=postgres password=secret dbname=pass_test sslmode=disable",
		port,
	)

	if err = pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("pgx", connectionString)
		if err != nil {
			return fmt.Errorf("failed open test connection: %w", err)
		}
		pingErr := db.Ping()
		if pingErr != nil {
			return fmt.Errorf("failed ping test db: %w", pingErr)
		}
		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	if migrationErr := runMigrations(db.DB, migrationPath); migrationErr != nil {
		if err = closeTestConnection(db, pool, resource); err != nil {
			return nil,
				nil,
				nil,
				fmt.Errorf("failed close pool conections on failed migration err: %w: %w", err, migrationErr)
		}
		return nil, nil, nil, fmt.Errorf("failed run migrations on test database: %w", migrationErr)
	}

	return db, pool, resource, nil
}

func closeTestConnection(db *sqlx.DB, pool *dockertest.Pool, resource *dockertest.Resource) error {
	_ = db.Close()
	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("failed purge docker resource: %w", err)
	}

	return nil
}

func runMigrations(db *sql.DB, migrationPath string) error {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("failed init postgres driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath,
		"metrics", driver)
	if err != nil {
		return fmt.Errorf("failed get migration instance: %w", err)
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed run migrations: %w", err)
		}
	}

	return nil
}
