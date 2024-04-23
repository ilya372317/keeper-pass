package app

import (
	"github.com/ilya372317/pass-keeper/internal/server/adapter/pgsqlrepo/userrepo"
	"github.com/ilya372317/pass-keeper/internal/server/config"
	"github.com/ilya372317/pass-keeper/internal/server/interceptor"
	"github.com/ilya372317/pass-keeper/internal/server/service/auth"
	"github.com/ilya372317/pass-keeper/internal/server/service/jwtmanager"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	pgsqlx *sqlx.DB

	conf config.Config
}

func NewContainer(conf config.Config, pgsqlx *sqlx.DB) *Container {
	return &Container{
		conf:   conf,
		pgsqlx: pgsqlx,
	}
}

func (c *Container) GetDefaultAuthService() *auth.Service {
	return auth.NewAuthService(c.GetJWTTokenManager(), c.GetPostgresqlUserRepo())
}

func (c *Container) GetPostgresqlUserRepo() *userrepo.Repository {
	return userrepo.New(c.pgsqlx)
}

func (c *Container) GetJWTTokenManager() *jwtmanager.JWTManager {
	return jwtmanager.New(c.conf.JWT.SecretKey, c.conf.JWT.TokenExpDuration)
}

func (c *Container) GetAuthInterceptor() *interceptor.AuthInterceptor {
	return interceptor.NewAuthInterceptor(c.GetJWTTokenManager(), c.GetPostgresqlUserRepo())
}
