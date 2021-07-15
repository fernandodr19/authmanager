package app

import (
	"github.com/fernandodr19/authmanager/pkg/config"
	"github.com/fernandodr19/authmanager/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/authmanager/pkg/gateway/encrypter"
	"github.com/fernandodr19/authmanager/pkg/gateway/repositories"
	"github.com/jackc/pgx/v4"
)

// App contains application's usecases
type App struct {
	Accounts *accounts.AccountsUsecase
}

// BuildApp builds application
func BuildApp(dbConn *pgx.Conn, cfg *config.Config, tokenGenerator accounts.TokenGenerator) (*App, error) {
	accRepo := repositories.NewAccountRepository(dbConn)

	return &App{
		Accounts: accounts.NewAccountsUsecase(accRepo, tokenGenerator, encrypter.Encrypter{}),
	}, nil
}
