package main

import (
	"SMA-TP1/mail"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func FournisseurSimple(prixMin int) {
	comm := mail.Register(mail.Fournisseur)
	fmt.Println("FournisseurSimple lancé")

	addresses := comm.ListAgents()
	for len(addresses) == 0 {
		time.Sleep(100 * time.Millisecond)
		addresses = comm.ListAgents()
	}

	// calcul de l'offre
	offre := int(float64(prixMin) * 1.5)
	memoire := make(map[uuid.UUID]int)

	for _, addresse := range addresses {
		if addresse.AgentType == mail.Acheteur {
			fmt.Println("FournisseurSimple : Envoi d'une offre de base prix", offre)
			id := randomID()
			comm.Send(addresse, MessageOffre{
				ID:            id,
				Fournisseur:   comm.GetMyAddress(),
				Prix:          offre,
				Reduction:     false,
				TypeTransport: TransportTrain,
				Origin:        "Paris",
				Destination:   "Lyon",
			})
			memoire[id] = offre
		}
	}

	for {
		message := attenteMessage(comm)
		switch msg := message.(type) {
		case MessageAcceptation:
			fmt.Println("FournisseurSimple : offre acceptée avec le message : ", msg.Message)
			return
		case MessageRefus:
			fmt.Println("FournisseurSimple : offre refusée avec le message : ", msg.Message)
			return
		case MessageContreOffre:
			fmt.Println("FournisseurSimple : Contre offre reçue :", msg.Prix)
			contreOffre, possible := StrategieVendeurSimple(msg.Prix, memoire[msg.IDOffre], prixMin, msg.Round)
			if !possible {
				fmt.Println("FournisseurSimple : prix trop bas")
				comm.Send(msg.Interlocuteur, MessageRefus{
					IDOffre: msg.IDOffre,
					Message: "Prix trop bas",
				})
				return
			}
			fmt.Println("FournisseurSimple : Envoi contre offre :", contreOffre)
			comm.Send(msg.Interlocuteur, MessageContreOffre{
				IDOffre:       msg.IDOffre,
				Round:         msg.Round,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
			memoire[msg.IDOffre] = contreOffre
		}
	}
}

func AcheteurSimple(prixMax int, aggressivite int) {
	comm := mail.Register(mail.Acheteur)
	fmt.Println("AcheteurSimple lancé")

	for {
		message := attenteMessage(comm)
		switch msg := message.(type) {
		case MessageOffre:
			fmt.Println("Acheteur simple : Offre reçue :", msg.Prix)
			contreOffre := StrategieAcheteurSimple(msg.Prix, prixMax, 1, aggressivite)
			fmt.Println("Acheteur simple : Envoi contre offre :", contreOffre)
			comm.Send(msg.Fournisseur, MessageContreOffre{
				IDOffre:       msg.ID,
				Round:         1,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
		case MessageContreOffre:
			if msg.Round == 3 {
				if msg.Prix <= prixMax {
					fmt.Println("Acheteur simple : acceptation", msg.Prix)
					comm.Send(msg.Interlocuteur, MessageAcceptation{
						IDOffre: msg.IDOffre,
						Message: "J'accepte l'offre",
					})
				} else {
					fmt.Println("AcheteurSimple : refus", msg.Prix)
					comm.Send(msg.Interlocuteur, MessageRefus{
						IDOffre: msg.IDOffre,
						Message: "Pas d'accord trouvé",
					})
				}
				return
			}
			contreOffre := StrategieAcheteurSimple(msg.Prix, prixMax, msg.Round+1, aggressivite)
			fmt.Println("Acheteur simple : Envoi contre offre :", contreOffre)
			comm.Send(msg.Interlocuteur, MessageContreOffre{
				IDOffre:       msg.IDOffre,
				Round:         msg.Round + 1,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
		case MessageRefus:
			fmt.Println("AcheteurSimple : offre refusée avec le message : ", msg.Message)
		}
	}

}

func attenteMessage(comm *mail.Box) any {
	for {
		message, ok := comm.Receive()
		if ok {
			return message
		}
		time.Sleep(100 * time.Millisecond)
	}
}
