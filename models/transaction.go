package models

import "time"

type Transaction struct {
	Id        string
	AccountId string
	Nonce     uint64
	Payload   string
	Metadata  string
	CreatedOn time.Time
}

func NewTransaction(id string, acccountId string, nonce uint64, payload string, metadata string) *Transaction {
	return &Transaction{
		Id:        id,
		AccountId: acccountId,
		Nonce:     nonce,
		Payload:   payload,
		Metadata:  metadata,
		CreatedOn: time.Now(),
	}
}

func CopyTransaction(other *Transaction) *Transaction {
	return &Transaction{
		Id:        other.Id,
		AccountId: other.AccountId,
		Nonce:     other.Nonce,
		Payload:   other.Payload,
		Metadata:  other.Metadata,
		CreatedOn: other.CreatedOn,
	}
}
