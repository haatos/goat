package goat_test

import (
	"testing"

	"github.com/haatos/goat"
	"github.com/haatos/goat/providers/faux"
	"github.com/stretchr/testify/assert"
)

func Test_UseProviders(t *testing.T) {
	a := assert.New(t)

	provider := &faux.Provider{}
	goat.UseProviders(provider)
	a.Equal(len(goat.GetProviders()), 1)
	a.Equal(goat.GetProviders()[provider.Name()], provider)
	goat.ClearProviders()
}

func Test_GetProvider(t *testing.T) {
	a := assert.New(t)

	provider := &faux.Provider{}
	goat.UseProviders(provider)

	p, err := goat.GetProvider(provider.Name())
	a.NoError(err)
	a.Equal(p, provider)

	_, err = goat.GetProvider("unknown")
	a.Error(err)
	a.Equal(err.Error(), "no provider for unknown exists")
	goat.ClearProviders()
}
