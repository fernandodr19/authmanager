package authenticator

import (
	"github.com/fernandodr19/authenticator/pkg/config"
	"github.com/fernandodr19/authenticator/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/authenticator/pkg/gateway/repositories"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	Accounts accounts.Usecase
}

func BuildApp(dbConn *pgxpool.Pool, cfg *config.Config) (*App, error) {
	accRepo := repositories.NewAccountRepository(dbConn)

	return &App{
		Accounts: accounts.NewAccountsUsecase(accRepo),
	}, nil
}
