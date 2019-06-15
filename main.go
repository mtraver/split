package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/mtraver/split/person"
)

const (
	nameSeparator = "_"
)

var (
	filePath string
)

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}

func parsePeople(people []string) (map[string]float64, error) {
	m := make(map[string]float64)

	for _, s := range people {
		if s == "" {
			continue
		}

		p, err := person.New(s)
		if err != nil {
			return m, err
		}

		subnames := strings.Split(p.Name, nameSeparator)
		total := p.Total()
		for _, name := range subnames {
			if _, ok := m[name]; !ok {
				m[name] = 0
			}
			m[name] += total / float64(len(subnames))
		}
	}

	return m, nil
}

func init() {
	flag.StringVar(&filePath, "f", "", "path to CSV file")
	flag.Usage = func() {
		message := `usage: split total [-f file] [person [person ...]]

Positional arguments (required):
  total
	the the total amount paid, including everything split
	amongst the group such as tax, tip, bag fee, etc.

Options:
`

		fmt.Fprintf(flag.CommandLine.Output(), message)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	total, err := strconv.ParseFloat(flag.Args()[0], 64)
	if err != nil {
		fatalf("Total amount must be a number: %v", err)
	}

	var people map[string]float64
	if filePath == "" {
		if flag.NArg() < 2 {
			fmt.Fprintf(flag.CommandLine.Output(), "one or more people required")
			flag.Usage()
			os.Exit(1)
		}

		var err error
		if people, err = parsePeople(flag.Args()[1:flag.NArg()]); err != nil {
			fatalf("Error during parsing: %v", err)
		}
	} else {
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			fatalf("Error reading file: %v", err)
		}

		if people, err = parsePeople(strings.Split(string(b), "\n")); err != nil {
			fatalf("Error during parsing: %v", err)
		}
	}

	var subtotal float64
	for _, amount := range people {
		subtotal += amount
	}

	if subtotal > total {
		fatalf("Subtotal ($%.2f) is greater than total ($%.2f)", subtotal, total)
	}

	// Compute the amount each person owes.
	payments := make([]person.Person, 0, len(people))
	for name, amount := range people {
		payments = append(payments, person.Person{
			Name:  name,
			Items: []float64{(amount / subtotal) * total},
		})
	}
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].Name < payments[j].Name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Subtotal\t$%.2f\t\n", subtotal)
	fmt.Fprintf(w, "Tax, tip, etc.\t$%.2f\t\n", total-subtotal)
	w.Flush()

	fmt.Println("")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, p := range payments {
		fmt.Fprintf(w, "%s\t$%.2f\t\n", p.Name, p.Items[0])
	}
	w.Flush()
}
