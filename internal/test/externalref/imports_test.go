package externalref

import (
	"testing"

	packageA "dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/internal/test/externalref/packageA"
	packageB "dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/internal/test/externalref/packageB"
	"github.com/stretchr/testify/require"
)

func TestParameters(t *testing.T) {
	b := &packageB.ObjectB{}
	_ = Container{
		ObjectA: &packageA.ObjectA{ObjectB: b},
		ObjectB: b,
	}
}

func TestGetSwagger(t *testing.T) {
	_, err := packageB.GetSwagger()
	require.Nil(t, err)

	_, err = packageA.GetSwagger()
	require.Nil(t, err)

	_, err = GetSwagger()
	require.Nil(t, err)
}
