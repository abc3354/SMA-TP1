package main

import (
	"SMA-TP1/mail"
)

func main() {
	ready := make(chan mail.Atom)
	go mail.Start(ready)
	<-ready

	go AcheteurSimple(200)
	go FournisseurSimple()

	select {}
}

// Mail
// - demander la liste des agents
// - envoyer un mail à un agent
// - s'enregistrer

// todo sauvegarder les messages dans l'agent
