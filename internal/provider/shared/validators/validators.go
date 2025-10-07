package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// StringLengthBetween validates that a string length is between min and max
func StringLengthBetween(min, max int) validator.String {
	return &stringLengthBetweenValidator{
		min: min,
		max: max,
	}
}

type stringLengthBetweenValidator struct {
	min int
	max int
}

func (v *stringLengthBetweenValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d characters", v.min, v.max)
}

func (v *stringLengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d characters", v.min, v.max)
}

func (v *stringLengthBetweenValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	length := len(value)

	if length < v.min || length > v.max {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid String Length",
			fmt.Sprintf("String length must be between %d and %d characters, got %d", v.min, v.max, length),
		)
	}
}

// StringMatches validates that a string matches a regular expression
func StringMatches(pattern string, message string) validator.String {
	regex := regexp.MustCompile(pattern)
	return &stringMatchesValidator{
		regex:   regex,
		pattern: pattern,
		message: message,
	}
}

type stringMatchesValidator struct {
	regex   *regexp.Regexp
	pattern string
	message string
}

func (v *stringMatchesValidator) Description(ctx context.Context) string {
	if v.message != "" {
		return v.message
	}
	return fmt.Sprintf("string must match pattern: %s", v.pattern)
}

func (v *stringMatchesValidator) MarkdownDescription(ctx context.Context) string {
	if v.message != "" {
		return v.message
	}
	return fmt.Sprintf("string must match pattern: `%s`", v.pattern)
}

func (v *stringMatchesValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !v.regex.MatchString(value) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid String Format",
			v.Description(ctx),
		)
	}
}

// StringNoLeadingTrailingSpaces validates that a string has no leading or trailing spaces
func StringNoLeadingTrailingSpaces() validator.String {
	return &stringNoLeadingTrailingSpacesValidator{}
}

type stringNoLeadingTrailingSpacesValidator struct{}

func (v *stringNoLeadingTrailingSpacesValidator) Description(ctx context.Context) string {
	return "string must not have leading or trailing spaces"
}

func (v *stringNoLeadingTrailingSpacesValidator) MarkdownDescription(ctx context.Context) string {
	return "string must not have leading or trailing spaces"
}

func (v *stringNoLeadingTrailingSpacesValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	trimmed := strings.TrimSpace(value)

	if value != trimmed {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid String Format",
			"String must not have leading or trailing spaces",
		)
	}
}

// StringOneOf validates that a string is one of the provided values
func StringOneOf(values ...string) validator.String {
	return &stringOneOfValidator{
		values: values,
	}
}

type stringOneOfValidator struct {
	values []string
}

func (v *stringOneOfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string must be one of: %s", strings.Join(v.values, ", "))
}

func (v *stringOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string must be one of: `%s`", strings.Join(v.values, "`, `"))
}

func (v *stringOneOfValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	for _, validValue := range v.values {
		if value == validValue {
			return
		}
	}

	response.Diagnostics.AddAttributeError(
		request.Path,
		"Invalid String Value",
		fmt.Sprintf("String must be one of: %s, got: %s", strings.Join(v.values, ", "), value),
	)
}

// StringIsURL validates that a string is a valid URL
func StringIsURL() validator.String {
	return StringMatches(
		`^https?://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(:[0-9]+)?(/.*)?$`,
		"string must be a valid URL (http or https)",
	)
}

// StringIsEmail validates that a string is a valid email address
func StringIsEmail() validator.String {
	return StringMatches(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"string must be a valid email address",
	)
}

// IAMResourceName validates IAM resource names
func IAMResourceName() validator.String {
	return StringMatches(
		`^[a-zA-Z0-9._-]+$`,
		"IAM resource name must contain only letters, numbers, dots, underscores, and hyphens",
	)
}

// IAMPermission validates IAM permission format
func IAMPermission() validator.String {
	return StringMatches(
		`^[a-zA-Z][a-zA-Z0-9]*\.[a-zA-Z][a-zA-Z0-9]*\.[a-zA-Z][a-zA-Z0-9]*$`,
		"IAM permission must be in format 'service.resource.action' (e.g., 'iam.groups.list')",
	)
}

// MemberIdentifier validates member identifier format (user:email or group:name)
func MemberIdentifier() validator.String {
	return StringMatches(
		`^(user:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}|group:[a-zA-Z0-9._-]+)$`,
		"member identifier must be in format 'user:email@domain.com' or 'group:groupname'",
	)
}
