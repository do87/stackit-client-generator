package issue579

//go:generate go run dev.azure.com/schwarzit/schwarzit.odj.core/_git/stackit-client-generator.git/cmd/oapi-codegen --package=issue579 --generate=types,skip-prune --alias-types -o issue.gen.go spec.yaml
