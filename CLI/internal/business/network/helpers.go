package network

import (
	"encoding/json"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
)

type nodeIPPayload struct {
	NodeID   string `json:"nodeId"`
	IP       string `json:"ip"`
	SubnetID string `json:"subnetId"`
}

func validateStorageServiceID(rt *meta.Runtime) error {
	return inputs.RequireNonEmpty("storage-svc-id", rt.String("storage-svc-id"))
}

func validateSubnetID(rt *meta.Runtime) error {
	return inputs.RequireNonEmpty("subnet-id", rt.String("subnet-id"))
}

func validateNodeID(rt *meta.Runtime) error {
	return inputs.RequireNonEmpty("node-id", rt.String("node-id"))
}

func buildNodeIPPayload(rt *meta.Runtime) ([]byte, error) {
	return json.Marshal(nodeIPPayload{
		NodeID:   rt.String("node-id"),
		IP:       rt.String("ip"),
		SubnetID: rt.String("subnet-id"),
	})
}
