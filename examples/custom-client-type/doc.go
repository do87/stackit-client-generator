package customclienttype

// This is an example of how to add a prefix to the name of the generated Client struct
// See https://github.com/do87/stackit-client-generator/issues/785 for why this might be necessary

//go:generate go run github.com/do87/stackit-client-generator/cmd/oapi-codegen -config cfg.yaml api.yaml
