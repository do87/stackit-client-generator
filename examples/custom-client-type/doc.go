package customclienttype

// This is an example of how to add a prefix to the name of the generated Client struct
// See https://dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/issues/785 for why this might be necessary

//go:generate go run dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/cmd/oapi-codegen -config cfg.yaml api.yaml
