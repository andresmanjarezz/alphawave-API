package codegenerator

import (
	"fmt"
	"math/rand"
	"time"
)

type CodeGenerator struct {
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (g *CodeGenerator) GenerateUniqueCode() string {
	rand.Seed(time.Now().UnixNano())

	code := rand.Intn(90000) + 10000

	return fmt.Sprintf("%05d", code)
}
