package cluster

import (
	"github.com/do87/stackit-client-generator/examples/ske-client/generated/cluster"
)

type CreateOrUpdateClusterResponse struct {
	cluster.ClientWithResponsesInterface
}

type DeleteClusterResponse struct {
	cluster.ClientWithResponsesInterface
}
