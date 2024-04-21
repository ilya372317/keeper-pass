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
	db   *sqlx.DB
	repo *userrepo.Repository
)

func TestMain(m *testing.M) {
	database, pool, resource, err := makeTestConnection("../../db/migrations")
	if err != nil {
		log.Fatal(err)
		return
	}
	db = database
	repo = userrepo.New(database)
	m.Run()
	if err = closeTestConnection(database, pool, resource); err != nil {
		log.Fatal(err)
		return
	}
}

func TestRepository_GetUserByEmail(t *testing.T) {
	type want struct {
		err  bool
		user domain.User
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
				user: domain.User{
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
				user: domain.User{
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

			got, err := repo.GetUserByEmail(ctx, tt.arg)

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

func fillUsersTable(t *testing.T, users []domain.User) {
	t.Helper()
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