package main

import (
	"github.com/arnopensource/SMA-NEGOCIATION-TP1/mail"

	"github.com/google/uuid"
)

type MessageOffre struct {
	ID            uuid.UUID
	Fournisseur   mail.Address
	Prix          int
	Reduction     bool
	TypeTransport TypeTransport
	Origin        string
	Destination   string
}

type MessageAcceptation struct {
	IDOffre uuid.UUID
	Message string
}

type MessageRefus struct {
	IDOffre uuid.UUID
	Message string
}

type MessageContreOffre struct {
	IDOffre       uuid.UUID
	Round         int
	Prix          int
	Interlocuteur mail.Address
}

type TypeTransport int

const (
	TransportAvion TypeTransport = iota
	TransportTrain
)

func randomID() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id
}
