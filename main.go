package main

import (
	"github.com/arnopensource/SMA-NEGOCIATION-TP1/mail"
	"time"
)

func main() {
	ready := make(chan mail.Atom)
	go mail.Start(ready)
	<-ready

	go AcheteurSimple(250, 1, StrategieAcheteurSimple, "Mike")
	go AcheteurSimple(300, 1, StrategieAcheteurSimple, "Mattis")
	go AcheteurSimple(400, 2, StrategieAcheteurSimple, "Marc")
	go FournisseurSimple(200, StrategieVendeurSimple, "Charles")
	go FournisseurSimple(250, StrategieVendeurSimple, "Caroline")
	go FournisseurSimple(200, StrategieVendeurBizarre, "Corentin")

	for {
		time.Sleep(1 * time.Hour)
	}
}

// Mail
// - demander la liste des agents
// - envoyer un mail à un agent
// - s'enregistrer

// todo sauvegarder les messages dans l'agent
