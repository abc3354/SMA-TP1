package main

import (
	"fmt"
	"time"

	"github.com/arnopensource/SMA-NEGOCIATION-TP1/mail"

	"github.com/google/uuid"
)

const MaxRound = 5
const RoundEngageant = 3

func FournisseurSimple(prixMin int, strategie StrategieFournisseur, nom string) {
	comm := mail.Register(mail.Fournisseur)
	fmt.Println(nom, "lancé")

	addresses := comm.ListAgents()
	for len(addresses) == 0 {
		time.Sleep(100 * time.Millisecond)
		addresses = comm.ListAgents()
	}

	// calcul de l'offre
	offre := int(float64(prixMin) * 1.5)
	memoire := make(map[uuid.UUID]int)

	var listeOffres []MessageContreOffre
	engage := false

	for _, addresse := range addresses {
		if addresse.AgentType == mail.Acheteur {
			fmt.Println(nom, ": Envoi d'une offre de base prix", offre)
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
			fmt.Println(nom, ": offre acceptée avec le message : ", msg.Message)
			return
		case MessageRefus:
			fmt.Println(nom, ": offre refusée avec le message : ", msg.Message)
			return
		case MessageContreOffre:
			fmt.Println(nom, ": Contre offre reçue :", msg.Prix)
			contreOffre, possible := strategie(msg.Prix, memoire[msg.IDOffre], prixMin, msg.Round)
			if !possible {
				fmt.Println(nom, ": prix trop bas")
				comm.Send(msg.Interlocuteur, MessageRefus{
					IDOffre: msg.IDOffre,
					Message: "Prix trop bas",
				})
				return
			}
			if msg.Round > RoundEngageant && !engage {
				fmt.Println(nom, ": C'est l'heure d'envoyer une offre engageante !")
				listeOffres = append(listeOffres, msg)
				go func() {
					time.Sleep(time.Second)
					comm.Send(comm.GetMyAddress(), MessageTimer{})
				}()
				continue
			}
			fmt.Println(nom, ": Envoi contre offre :", contreOffre)
			comm.Send(msg.Interlocuteur, MessageContreOffre{
				IDOffre:       msg.IDOffre,
				Round:         msg.Round,
				Engageant:     msg.Round > RoundEngageant,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
			memoire[msg.IDOffre] = contreOffre
		case MessageTimer:
			if engage {
				continue
			}
			engage = true
			best := listeOffres[0]
			for _, offre := range listeOffres {
				if offre.Prix < best.Prix {
					best = offre
				}
			}
			fmt.Println(nom, "simple : Meilleure offre : ", best.Prix)
			contreOffre, possible := strategie(best.Prix, memoire[best.IDOffre], prixMin, best.Round)
			if !possible {
				fmt.Println(nom, ": prix trop bas")
				comm.Send(best.Interlocuteur, MessageRefus{
					IDOffre: best.IDOffre,
					Message: "Prix trop bas",
				})
				return
			}
			fmt.Println(nom, ": Envoi contre offre :", contreOffre)
			comm.Send(best.Interlocuteur, MessageContreOffre{
				IDOffre:       best.IDOffre,
				Round:         best.Round,
				Engageant:     true,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
		}
	}
}

func AcheteurSimple(prixMax int, aggressivite int, strategie StrategieAcheteur, nom string) {
	comm := mail.Register(mail.Acheteur)
	fmt.Println(nom, "lancé")

	var listeOffres []MessageContreOffre
	var engage = false

	for {
		message := attenteMessage(comm)
		switch msg := message.(type) {
		case MessageOffre:
			fmt.Println(nom, ": Offre reçue :", msg.Prix)
			contreOffre := strategie(msg.Prix, prixMax, 1, aggressivite)
			fmt.Println(nom, ": Envoi contre offre :", contreOffre)
			comm.Send(msg.Fournisseur, MessageContreOffre{
				IDOffre:       msg.ID,
				Round:         1,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
		case MessageContreOffre:
			fmt.Println(nom, ": Contre offre reçue :", msg.Prix)
			if msg.Engageant && !engage {
				fmt.Println(nom, ": C'est l'heure de s'engager !")
				listeOffres = append(listeOffres, msg)
				go func() {
					time.Sleep(time.Second)
					comm.Send(comm.GetMyAddress(), MessageTimer{})
				}()
				continue
			}
			if msg.Round == MaxRound {
				if msg.Prix <= prixMax {
					fmt.Println(nom, ": acceptation", msg.Prix)
					comm.Send(msg.Interlocuteur, MessageAcceptation{
						IDOffre: msg.IDOffre,
						Message: "J'accepte l'offre",
					})
				} else {
					fmt.Println(nom, ": refus", msg.Prix)
					if msg.Engageant {
						fmt.Println(nom, ": l'offre était engageante, pénalité appliquée")
					}
					comm.Send(msg.Interlocuteur, MessageRefus{
						IDOffre: msg.IDOffre,
						Message: "Pas d'accord trouvé",
					})
				}
				return
			}
			contreOffre := strategie(msg.Prix, prixMax, msg.Round+1, aggressivite)
			fmt.Println(nom, ": Envoi contre offre :", contreOffre)
			comm.Send(msg.Interlocuteur, MessageContreOffre{
				IDOffre:       msg.IDOffre,
				Round:         msg.Round + 1,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
		case MessageRefus:
			fmt.Println(nom, ": offre refusée avec le message : ", msg.Message)
		case MessageTimer:
			if engage {
				continue
			}
			engage = true
			best := listeOffres[0]
			for _, offre := range listeOffres {
				if offre.Prix < best.Prix {
					best = offre
				}
			}
			fmt.Println(nom, ": Meilleure offre : ", best.Prix)
			contreOffre := strategie(best.Prix, prixMax, best.Round+1, aggressivite)
			fmt.Println(nom, ": Envoi contre offre :", contreOffre)
			comm.Send(best.Interlocuteur, MessageContreOffre{
				IDOffre:       best.IDOffre,
				Round:         best.Round + 1,
				Prix:          contreOffre,
				Interlocuteur: comm.GetMyAddress(),
			})
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
