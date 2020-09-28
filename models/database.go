package models

import "time"

// SQLDataBase struct
type SQLDataBase struct {
	Server          string   `toml:"Server"`
	Database        string   `toml:"Database"`
	Port			int      `toml:"Port"`
	ApplicationName string   `toml:"ApplicationName"`
	MaxIdleConns    int      `toml:"MaxIdleConns"`
	MaxOpenConns    int      `toml:"MaxOpenConns"`
	ConnMaxLifetime duration `toml:"ConnMaxLifetime"`
	UserID          string
	Password        string
}

type Transaction struct {
	UserId  int       `json:"user_id"`
	Comment string    `json:"comment"`
	Amount  float64   `json:"amount"`
	Date    time.Time `json:"date"`
}

type Balance struct {
	UserId int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

type AmountCurrency struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
