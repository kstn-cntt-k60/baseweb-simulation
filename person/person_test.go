package person

import (
	"math/rand"
	"testing"
)

func TestPerson(t *testing.T) {
	rand.Seed(100)
	for i := 0; i < 10; i++ {
		person := Person{}
		err := GeneratePerson(&person)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(person)
	}
}
