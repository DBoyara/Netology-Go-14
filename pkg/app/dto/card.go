package dto

import (
	"github.com/DBoyara/Netology-Go-14/pkg/card"
	"time"
)

type CardDTO struct {
	Id      int64       `json:"id"`
	Number  int64       `json:"number"`
	Balance int64       `json:"balance"`
	Issuer  string      `json:"issuer"`
	Type    string      `json:"type"`
	OwnerId card.UserID `json:"owner_id"`
	Status  string      `json:"status"`
	Created time.Time   `json:"created"`
}

type ErrDTO struct {
	Err string `json:"error"`
}

type TransactionDTO struct {
	Id          int64     `json:"id"`
	CardId      int64     `json:"card_id"`
	Amount      int64     `json:"amount"`
	Created     time.Time `json:"created"`
	Status      string    `json:"status"`
	MccId       int64     `json:"mcc_id"`
	Description string    `json:"description"`
	IconId      int64     `json:"icon_id"`
}

type MostOftenBought struct {
	MccId       int64  `json:"mcc_id"`
	Count       int64  `json:"count"`
	Description string `json:"description"`
}

type MostPaid struct {
	MccId       int64  `json:"mcc_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}
