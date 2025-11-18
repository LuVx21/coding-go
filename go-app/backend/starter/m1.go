package starter

import "fmt"

type Binding struct{}

func NewBinding() *Binding {
	return &Binding{}
}

func (a *Binding) Hello(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
