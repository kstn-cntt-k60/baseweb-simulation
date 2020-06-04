package product

import (
	"math/rand"
	"testing"
)

func TestProduct(t *testing.T) {
	rand.Seed(200)

	for i := 0; i < 10; i++ {
		p := Product{}
		err := GenerateProduct(&p)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(p)
	}
}
