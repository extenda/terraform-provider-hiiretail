package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func runValidateString(v validator.String, val string) (hasErr bool) {
	var req validator.StringRequest
	req.Path = path.Root("test")
	req.ConfigValue = types.StringValue(val)
	var resp validator.StringResponse
	v.ValidateString(context.Background(), req, &resp)
	return resp.Diagnostics.HasError()
}

func TestStringLengthBetween(t *testing.T) {
	v := StringLengthBetween(3, 5)
	if runValidateString(v, "ab") == false {
		t.Fatalf("expected error for too short")
	}
	if runValidateString(v, "abcd") == true {
		t.Fatalf("expected no error for valid length")
	}
	if runValidateString(v, "abcdef") == false {
		t.Fatalf("expected error for too long")
	}
}

func TestStringMatchesAndNoSpaces(t *testing.T) {
	v := StringMatches(`^[0-9]+$`, "must be digits")
	if runValidateString(v, "abc") == false {
		t.Fatalf("expected regex mismatch to produce error")
	}
	if runValidateString(StringNoLeadingTrailingSpaces(), " a ") == false {
		t.Fatalf("expected leading/trailing spaces to be invalid")
	}
}

func TestStringOneOf(t *testing.T) {
	v := StringOneOf("a", "b", "c")
	if runValidateString(v, "d") == false {
		t.Fatalf("expected error for value not in set")
	}
	if runValidateString(v, "b") == true {
		t.Fatalf("expected no error for allowed value")
	}
}
