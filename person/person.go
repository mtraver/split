package person

import (
	"fmt"
	"strconv"
	"strings"
)

type Person struct {
	Name  string
	Items []float64
}

func New(s string) (Person, error) {
	parts := strings.Split(s, ",")
	if len(parts) < 2 {
		return Person{}, fmt.Errorf("person: must have a name and at least one amount")
	}

	p := Person{
		Name:  parts[0],
		Items: make([]float64, 0, len(parts)-1),
	}

	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}

		amount, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return p, err
		}

		p.Items = append(p.Items, amount)
	}

	return p, nil
}

func (p Person) Total() float64 {
	var total float64
	for _, a := range p.Items {
		total += a
	}
	return total
}
