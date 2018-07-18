package main

import (
	"fmt"

	"github.com/trustmaster/goflow"
)

type Greeter struct {
	flow.Component
	Name <-chan string
	Res  chan<- string
}

func (g *Greeter) OnName(name string) {
	greeting := fmt.Sprintf("Hello, %s!", name)
	g.Res <- greeting
}

type Printer struct {
	flow.Component
	Line <-chan string
}

func (p *Printer) OnLine(line string) {
	fmt.Println(line)
}

type GreetingApp struct {
	flow.Graph
}

func NewGreetingApp() *GreetingApp {
	n := &GreetingApp{}
	n.InitGraphState() 

	n.Add(&Greeter{}, "greeter")
	n.Add(&Printer{}, "printer")

	n.Connect("greeter", "Res", "printer", "Line")

	n.MapInPort("In", "greeter", "Name")
	return n
}

func main() {
	net := NewGreetingApp()
	in := make(chan string)
	net.SetInPort("In", in)
	flow.RunNet(net)

	in <- "John"
	in <- "Boris"
	in <- "Hanna"

	close(in)
	<-net.Wait()
}
