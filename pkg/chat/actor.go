package chat

import (
	"fmt"
)

const ActorChanSize = 16

type (
	Action     func()
	ActionChan chan Action
)

type stop struct{}

type Actor interface {
	Run()
	Stop()
}

type ActorBase struct {
	Name       string
	actQ       ActionChan
	stopSignal chan struct{}
	stopped    bool
}

func NewActorBase(name string) *ActorBase {
	return &ActorBase{
		Name:       name,
		actQ:       make(ActionChan, ActorChanSize),
		stopSignal: make(chan struct{}, 1),
	}
}

func (a *ActorBase) Run() {
	fmt.Printf("Actor %s started\n", a.Name)
	defer fmt.Printf("Actor %s run loop stopped\n", a.Name)

	for {
		select {
		case <-a.stopSignal:
			fmt.Printf("Actor %s got stop signal\n", a.Name)
			return
		case act := <-a.actQ:
			act()
		}
	}
}

func (a *ActorBase) StopActor() {
	if a.stopped {
		return
	}
	a.stopSignal <- stop{}
	fmt.Printf("Actor %s stopping\n", a.Name)
	a.stopped = true
}

type Director struct {
	actors []Actor
}

func NewDirector(actors ...Actor) *Director {
	return &Director{
		actors: actors,
	}
}

func (d *Director) Run() {
	for _, actor := range d.actors {
		go actor.Run()
	}
}

func (d *Director) Stop() {
	for _, actor := range d.actors {
		actor.Stop()
	}
}
