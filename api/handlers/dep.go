package handlers

import (
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	DB      *pgxpool.Pool
	Mailer  *services.MailService
	Cacher  *services.CacheService
	Batcher *services.BatchService
}

func (d *Dependencies) NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		db:     d.DB,
		mailer: d.Mailer,
		cacher: d.Cacher,
	}
}

func (d *Dependencies) NewUsersHandler() *UsersHandler {
	return &UsersHandler{
		db: d.DB,
	}
}

func (d *Dependencies) NewTickerHandler() *TickerHandler {
	return &TickerHandler{
		db: d.DB,
	}
}

func (d *Dependencies) NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		db:      d.DB,
		mailer:  d.Mailer,
		batcher: d.Batcher,
	}
}
