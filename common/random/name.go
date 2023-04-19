package random

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

func Name() string {
	c := gofakeit.NewCrypto()
	return fmt.Sprintf("%s %s %s %s %s", c.Name(), c.Name(), c.Name(), c.Name(), c.Name())
}
