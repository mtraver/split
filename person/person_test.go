package person

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name  string
		Str   string
		Valid bool
		Want  Person
	}{
		{
			name:  "1_item",
			Str:   "Aeolus,10.0",
			Valid: true,
			Want: Person{
				Name:  "Aeolus",
				Items: []float64{10.0},
			},
		},
		{
			name:  "2_item",
			Str:   "Aether,10.0,7.5",
			Valid: true,
			Want: Person{
				Name:  "Aether",
				Items: []float64{10.0, 7.5},
			},
		},
		{
			name:  "0_items",
			Str:   "Hades",
			Valid: false,
			Want:  Person{},
		},
		{
			name:  "trailing_comma",
			Str:   "Dionysus,5.25,",
			Valid: true,
			Want: Person{
				Name:  "Dionysus",
				Items: []float64{5.25},
			},
		},
		{
			name:  "empty_string",
			Str:   "",
			Valid: false,
			Want:  Person{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := New(c.Str)

			if c.Valid && err != nil {
				t.Errorf("Expected valid, got error: %v", err)
				return
			}

			if !c.Valid && err == nil {
				t.Errorf("Expected invalid, got nil error")
				return
			}

			if !reflect.DeepEqual(got, c.Want) {
				t.Errorf("Got %v, want %v", got, c.Want)
			}
		})
	}
}

func TestTotal(t *testing.T) {
	cases := []struct {
		name   string
		person Person
		want   float64
	}{
		{
			name:   "zero_value",
			person: Person{},
			want:   0.0,
		},
		{
			name: "non_zero",
			person: Person{
				Name:  "Athena",
				Items: []float64{15.47, 7.00, 9.50},
			},
			want: 31.97,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.person.Total()
			if got != c.want {
				t.Errorf("Got %v, want %v", got, c.want)
			}
		})
	}
}
