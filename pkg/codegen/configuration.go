package codegen

import (
	"errors"
	"reflect"
)

type AdditionalImport struct {
	Alias   string `yaml:"alias,omitempty"`
	Package string `yaml:"package"`
}

// Configuration defines code generation customizations
type Configuration struct {
	PackageName       string               `yaml:"package"` // PackageName to generate
	Generate          GenerateOptions      `yaml:"generate,omitempty"`
	Compatibility     CompatibilityOptions `yaml:"compatibility,omitempty"`
	OutputOptions     OutputOptions        `yaml:"output-options,omitempty"`
	ImportMapping     map[string]string    `yaml:"import-mapping,omitempty"` // ImportMapping specifies the golang package path for each external reference
	AdditionalImports []AdditionalImport   `yaml:"additional-imports,omitempty"`
	Tidy              Tidy                 `yaml:"tidy,omitempty"`
}

type Tidy struct {
	Verbose   bool       `yaml:"verbose,omitempty"`
	Functions []TidyRule `yaml:"functions,omitempty"`
	Params    []TidyRule `yaml:"params,omitempty"`
	Schemas   []TidyRule `yaml:"schemas,omitempty"`
}

// Tidy is a list of rules for replacing strings in the openAPI file
type TidyRule struct {
	Replace string `yaml:"replace,omitempty"`
	With    string `yaml:"with,omitempty"`
	Prefix  bool   `yaml:"prefix,omitempty"`
	Suffix  bool   `yaml:"suffix,omitempty"`
	All     bool   `yaml:"all,omitempty"`
	Match   bool   `yaml:"match,omitempty"`
}

// GenerateOptions specifies which supported output formats to generate.
type GenerateOptions struct {
	ChiServer     bool `yaml:"chi-server,omitempty"`     // ChiServer specifies whether to generate chi server boilerplate
	EchoServer    bool `yaml:"echo-server,omitempty"`    // EchoServer specifies whether to generate echo server boilerplate
	GinServer     bool `yaml:"gin-server,omitempty"`     // GinServer specifies whether to generate gin server boilerplate
	GorillaServer bool `yaml:"gorilla-server,omitempty"` // GorillaServer specifies whether to generate Gorilla server boilerplate
	Strict        bool `yaml:"strict-server,omitempty"`  // Strict specifies whether to generate strict server wrapper
	Client        bool `yaml:"client,omitempty"`         // Client specifies whether to generate client boilerplate
	Models        bool `yaml:"models,omitempty"`         // Models specifies whether to generate type definitions
	EmbeddedSpec  bool `yaml:"embedded-spec,omitempty"`  // Whether to embed the swagger spec in the generated code
}

// CompatibilityOptions specifies backward compatibility settings for the
// code generator.
type CompatibilityOptions struct {
	// In the past, we merged schemas for `allOf` by inlining each schema
	// within the schema list. This approach, though, is incorrect because
	// `allOf` merges at the schema definition level, not at the resulting model
	// level. So, new behavior merges OpenAPI specs but generates different code
	// than we have in the past. Set OldMergeSchemas to true for the old behavior.
	// Please see https://dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/issues/531
	OldMergeSchemas bool `yaml:"old-merge-schemas,omitempty"`
	// Enum values can generate conflicting typenames, so we've updated the
	// code for enum generation to avoid these conflicts, but it will result
	// in some enum types being renamed in existing code. Set OldEnumConflicts to true
	// to revert to old behavior. Please see:
	// Please see https://dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/issues/549
	OldEnumConflicts bool `yaml:"old-enum-conflicts,omitempty"`
	// It was a mistake to generate a go type definition for every $ref in
	// the OpenAPI schema. New behavior uses type aliases where possible, but
	// this can generate code which breaks existing builds. Set OldAliasing to true
	// for old behavior.
	// Please see https://dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/issues/549
	OldAliasing bool `yaml:"old-aliasing,omitempty"`
	// When an object contains no members, and only an additionalProperties specification,
	// it is flattened to a map. Set
	DisableFlattenAdditionalProperties bool `yaml:"disable-flatten-additional-properties,omitempty"`
	// When an object property is both required and readOnly the go model is generated
	// as a pointer. Set DisableRequiredReadOnlyAsPointer to true to mark them as non pointer.
	// Please see https://dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/issues/604
	DisableRequiredReadOnlyAsPointer bool `yaml:"disable-required-readonly-as-pointer,omitempty"`
	// When set to true, always prefix enum values with their type name instead of only
	// when typenames would be conflicting.
	AlwaysPrefixEnumValues bool `yaml:"always-prefix-enum-values,omitempty"`
	// Our generated code for Chi has historically inverted the order in which Chi middleware is
	// applied such that the last invoked middleware ends up executing first in the Chi chain
	// This resolves the behavior such that middlewares are chained in the order they are invoked.
	// Please see https://dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/issues/786
	ApplyChiMiddlewareFirstToLast bool `yaml:"apply-chi-middleware-first-to-last,omitempty"`
}

// OutputOptions are used to modify the output code in some way.
type OutputOptions struct {
	SkipFmt       bool              `yaml:"skip-fmt,omitempty"`       // Whether to skip go imports on the generated code
	SkipPrune     bool              `yaml:"skip-prune,omitempty"`     // Whether to skip pruning unused components on the generated code
	IncludeTags   []string          `yaml:"include-tags,omitempty"`   // Only include operations that have one of these tags. Ignored when empty.
	ExcludeTags   []string          `yaml:"exclude-tags,omitempty"`   // Exclude operations that have one of these tags. Ignored when empty.
	UserTemplates map[string]string `yaml:"user-templates,omitempty"` // Override built-in templates from user-provided files

	ExcludeSchemas     []string `yaml:"exclude-schemas,omitempty"`      // Exclude from generation schemas with given names. Ignored when empty.
	ResponseTypeSuffix string   `yaml:"response-type-suffix,omitempty"` // The suffix used for responses types
	ClientTypeName     string   `yaml:"client-type-name,omitempty"`     // Override the default generated client type with the value

	SplitByTags    FileSplit            `yaml:"split-by-tags,omitempty"`   // FileSplit splits to separate files/dirs according to tags
	ExtendResponse []ExtendResponseItem `yaml:"extend-response,omitempty"` // FileSplit splits to separate files/dirs according to tags
	CustomDoer     CustomDoer           `yaml:"custom-doer,omitempty"`     // FileSplit splits to separate files/dirs according to tags
	Copy           []CopyItem           `yaml:"copy,omitempty"`            // Files to copy
}

type FileSplit struct {
	Verbose bool     `yaml:"verbose,omitempty"`
	Enabled bool     `yaml:"enabled,omitempty"`
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type ExtendResponseItem struct {
	Field       string   `yaml:"field,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Imports     []string `imports:"imports,omitempty"`
	ApplyTo     []string `yaml:"apply-to,omitempty"`
	Set         string   `yaml:"set,omitempty"`
}

type CustomDoer struct {
	Enabled  bool   `yaml:"enabled,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Import   string `imports:"import,omitempty"`
	Generics string `imports:"generics,omitempty"`
}

type CopyItem struct {
	From string     `yaml:"from,omitempty"`
	To   string     `yaml:"to,omitempty"`
	Tidy []TidyRule `imports:"tidy,omitempty"`
}

// UpdateDefaults sets reasonable default values for unset fields in Configuration
func (o Configuration) UpdateDefaults() Configuration {
	if reflect.ValueOf(o.Generate).IsZero() {
		o.Generate = GenerateOptions{
			EchoServer:   true,
			Models:       true,
			EmbeddedSpec: true,
		}
	}
	return o
}

// Validate checks whether Configuration represent a valid configuration
func (o Configuration) Validate() error {
	if o.PackageName == "" {
		return errors.New("package name must be specified")
	}

	// Only one server type should be specified at a time.
	nServers := 0
	if o.Generate.ChiServer {
		nServers++
	}
	if o.Generate.EchoServer {
		nServers++
	}
	if o.Generate.GinServer {
		nServers++
	}
	if nServers > 1 {
		return errors.New("only one server type is supported at a time")
	}
	return nil
}
