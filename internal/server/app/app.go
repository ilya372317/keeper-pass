package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ilya372317/pass-keeper/internal/server/adapter/filerepo/file"
	"github.com/ilya372317/pass-keeper/internal/server/config"
	"github.com/ilya372317/pass-keeper/internal/server/service/keyring"
	"github.com/ilya372317/pass-keeper/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	timeForGetKeyring        = time.Second * 5
	timeForCreateMinIOBucket = time.Second * 5
)

type App struct {
	c    *Container
	conf config.Config
}

func New(configPath, masterKey string) (*App, error) {
	conf, err := config.New(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed create new config: %w", err)
	}
	logger.InitMust()

	app := App{
		conf: conf,
	}
	pgsqlxConnect, err := app.newPgSqlxConnect(conf.MainDB)
	if err != nil {
		return nil, fmt.Errorf("failed make connection with postgresql db: %w", err)
	}

	minIOClient, err := app.newMinIOConnection(conf.MinIO)
	if err != nil {
		return nil, fmt.Errorf("failed cerate minio connection: %w", err)
	}
	app.c = NewContainer(conf, pgsqlxConnect, minIOClient)
	ctx, stop := context.WithTimeout(context.Background(), timeForGetKeyring)
	defer stop()
	kring := keyring.New([]byte(masterKey), app.c.GetKeyRepository())
	if err = kring.InitGeneralKey(ctx); err != nil {
		return nil, fmt.Errorf("failed init general key: %w", err)
	}
	app.c.SetKeyring(kring)
	return &app, nil
}

func (a *App) newMinIOConnection(conf config.MinIOConfig) (*minio.Client, error) {
	client, err := minio.New(conf.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Login, conf.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed create minio client instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeForCreateMinIOBucket)
	defer cancel()
	bucketExist, err := client.BucketExists(ctx, file.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed check if main minio bucket exists: %w", err)
	}
	if !bucketExist {
		if err = client.MakeBucket(ctx, file.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed create main minio bucket: %w", err)
		}
	}

	return client, nil
}

func (a *App) newPgSqlxConnect(cfg config.SQLConfig) (*sqlx.DB, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("host=%s port=%s ", cfg.Host, cfg.Port))
	builder.WriteString(fmt.Sprintf("user=%s password=%s ", cfg.User, cfg.Password))
	builder.WriteString(fmt.Sprintf("dbname=%s ", cfg.DBName))
	builder.WriteString(fmt.Sprintf("timezone=%s ", cfg.Timezone))
	builder.WriteString("sslmode=disable ")

	params := builder.String()

	db, err := sqlx.Open("pgx", params)
	if err != nil {
		return nil, fmt.Errorf("failed open pgsql connection. invalid config data: %w", err)
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}
