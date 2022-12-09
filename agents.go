package main

import (
	"SMA-TP1/mail"
	"fmt"
	"time"
)

func FournisseurSimple() {
	comm := mail.Register(mail.Fournisseur)
	fmt.Println("FournisseurSimple lancé")

	addresses := comm.ListAgents()
	for len(addresses) == 0 {
		time.Sleep(100 * time.Millisecond)
		addresses = comm.ListAgents()
	}

	for _, addresse := range addresses {
		if addresse.AgentType == mail.Acheteur {
			fmt.Println("Envoi d'une offre de base")
			comm.Send(addresse, MessageOffre{
				ID:            randomID(),
				Fournisseur:   comm.GetMyAddress(),
				Prix:          100,
				Reduction:     false,
				TypeTransport: TransportTrain,
				Origin:        "Paris",
				Destination:   "Lyon",
			})
		}
	}

	message := attenteMessage(comm)
	fmt.Println(">", message)
	switch msg := message.(type) {
	case MessageAcceptation:
		fmt.Println("FournisseurSimple : offre acceptée avec le message : ", msg.Message)
		return
	}

}

func AcheteurSimple(prixDesire int) {
	comm := mail.Register(mail.Acheteur)
	fmt.Println("AcheteurSimple lancé")

	for {
		message := attenteMessage(comm)
		switch msg := message.(type) {
		case MessageOffre:
			fmt.Println("AcheteurSimple : l'offre est ok")
			if msg.Prix < prixDesire {
				comm.Send(msg.Fournisseur, MessageAcceptation{
					IDOffre: msg.ID,
					Message: "J'accepte l'offre",
				})
				return
			} else {
				fmt.Println("Offre n'accepté pas: Prix plus grand que le désiré.")
			}
		}
	}

}

func attenteMessage(comm *mail.Box) any {
	for {
		message, ok := comm.Receive()
		if ok {
			fmt.Println(message)
			return message
		}
		time.Sleep(100 * time.Millisecond)
	}
}
