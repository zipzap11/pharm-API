package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductNameShortener(t *testing.T) {
	name := "Diatomix Pembasmi Kutu Pada Hewan Peliharaan & Ternak 25 g"
	t.Log("len = ", len(name))
	convName := productNameShortener(name)
	t.Log("len after = ", len(convName))
	assert.LessOrEqual(t, len(convName), 50)
}