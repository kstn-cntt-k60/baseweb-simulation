package product

import (
	"testing"

	"github.com/bxcodec/faker/v3"
)

type Product struct {
	Name string `faker:"name"`
}

func fibonacci(n int) int {
	if n <= 1 {
		return 1
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}
}

func TestProduct(t *testing.T) {
	for i := 0; i < 10; i++ {
		p := Product{}
		err := faker.FakeData(&p)
		if err != nil {
			t.Error(err)
		}
		t.Log(p)
	}
}

func TestTung(t *testing.T) {
	for i := 0; i < 10; i++ {
		p := Product{}
		err := faker.FakeData(&p)
		if err != nil {
			t.Error(err)
		}
		t.Log(p)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fibonacci(30)
	}
}
