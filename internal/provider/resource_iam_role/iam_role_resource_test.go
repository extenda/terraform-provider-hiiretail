package resource_iam_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestIamRoleResourceSchema_Basics(t *testing.T) {
	ctx := context.Background()
	sch := IamRoleResourceSchema(ctx)
	require.NotNil(t, sch.Attributes)

	// Expected attributes
	expected := []string{"id", "name", "permissions", "fixed_bindings", "return_aliases"}
	for _, k := range expected {
		_, ok := sch.Attributes[k]
		require.Truef(t, ok, "expected attribute %s in schema", k)
	}

	// id should be a string attribute with validators
	idAttr, ok := sch.Attributes["id"].(schema.StringAttribute)
	require.True(t, ok, "id attribute should be StringAttribute")
	require.NotEmpty(t, idAttr.Validators, "id attribute should have validators (regex)")

	// permissions should be a list of strings
	permAttr, ok := sch.Attributes["permissions"].(schema.ListAttribute)
	require.True(t, ok, "permissions should be ListAttribute")
	require.Equal(t, types.StringType, permAttr.ElementType)
}
