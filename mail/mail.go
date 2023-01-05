package mail

import (
	"fmt"

	"github.com/google/uuid"
)

type Atom struct{}

type Message any

type agent struct {
	aType AgentType
	inbox []Message
}

type staticData struct {
	agents map[uuid.UUID]agent

	commands chan Command // chan = tube
}

type Command any

type RegisterCommand struct {
	responder chan uuid.UUID
	agentType AgentType
}

type SendCommand struct {
	address Address
	message Message
}

type ReceiveCommand struct {
	agentID   uuid.UUID
	responder chan any
}

type ListCommand struct {
	agentID   uuid.UUID
	responder chan []Address
}

type QuitCommand struct{}

var data staticData

func Start(ready chan Atom) {
	data = staticData{
		agents:   make(map[uuid.UUID]agent),
		commands: make(chan Command),
	}
	ready <- Atom{}
	for {
		cmd := <-data.commands
		switch cmd := cmd.(type) {
		case RegisterCommand:
			id := getUUID()
			data.agents[id] = agent{
				aType: cmd.agentType,
				inbox: nil,
			}
			cmd.responder <- id
		case SendCommand:
			id := cmd.address.id
			agent := data.agents[id]
			agent.inbox = append(agent.inbox, cmd.message)
			data.agents[id] = agent
		case QuitCommand:
			fmt.Println("Quitting")
		case ReceiveCommand:
			agent := data.agents[cmd.agentID]
			if len(agent.inbox) == 0 {
				cmd.responder <- Atom{}
				continue
			}
			message := agent.inbox[0]
			agent.inbox = agent.inbox[1:]
			data.agents[cmd.agentID] = agent
			cmd.responder <- message
		case ListCommand:
			var listAgent []Address
			for id, agent := range data.agents {
				if id == cmd.agentID {
					continue
				}
				listAgent = append(listAgent, Address{
					AgentType: agent.aType,
					id:        id,
				})
			}
			cmd.responder <- listAgent
		default:
			panic("unknown cmd")
		}
	}
}

type AgentType int

const (
	Fournisseur AgentType = iota
	Acheteur    AgentType = iota
)

func Register(agentType AgentType) *Box {
	responder := make(chan uuid.UUID)

	data.commands <- RegisterCommand{
		responder: responder,
		agentType: agentType,
	}

	return &Box{
		Address{
			AgentType: agentType,
			id:        <-responder,
		},
	}
}

type Box struct {
	Address
}

type Address struct {
	AgentType AgentType
	id        uuid.UUID
}

func (box *Box) Send(address Address, message Message) {
	data.commands <- SendCommand{
		address: address,
		message: message,
	}
}

func (box *Box) Receive() (Message, bool) {
	responder := make(chan any)

	data.commands <- ReceiveCommand{
		responder: responder,
		agentID:   box.id,
	}

	response := <-responder
	switch result := response.(type) {
	case Atom:
		return Message(0), false
	case Message:
		return result, true
	default:
		panic("unknown")
	}
}

func (box *Box) ListAgents() []Address {
	responder := make(chan []Address)

	data.commands <- ListCommand{
		agentID:   box.id,
		responder: responder,
	}

	return <-responder
}

func (box *Box) GetMyAddress() Address {
	return box.Address
}

func getUUID() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id
}

func Quit() {
	data.commands <- QuitCommand{}
}
