package main

import (
	"SMA-TP1/mail"
	"fmt"
)

func main() {
	ready := make(chan mail.Atom)
	go mail.Start(ready)
	<-ready

	alice := mail.Register(mail.Fournisseur)
	bob := mail.Register(mail.Acheteur)

	addresses := bob.ListAgents()
	bob.Send(addresses[0], 42)

	msg, ok := alice.Receive()
	fmt.Println(msg, ok)

	msg, ok = bob.Receive()
	fmt.Println(msg, ok)

}

// Mail
// - demander la liste des agents
// - envoyer un mail à un agent
// - s'enregistrer

// todo sauvegarder les messages dans l'agent
