package issue579

//go:generate go run github.com/do87/stackit-client-generator/cmd/oapi-codegen --package=issue579 --generate=types,skip-prune --alias-types -o issue.gen.go spec.yaml
