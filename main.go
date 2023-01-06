package main

import (
	"github.com/arnopensource/SMA-NEGOCIATION-TP1/mail"
	"time"
)

func main() {
	ready := make(chan mail.Atom)
	go mail.Start(ready)
	<-ready

	go AcheteurSimple(250, 1, StrategieAcheteurSimple)
	go FournisseurSimple(200, StrategieVendeurSimple)
	go FournisseurSimple(250, StrategieVendeurSimple)
	go FournisseurSimple(200, StrategieVendeurBizarre)

	for {
		time.Sleep(1 * time.Hour)
	}
}

// Mail
// - demander la liste des agents
// - envoyer un mail à un agent
// - s'enregistrer

// todo sauvegarder les messages dans l'agent
