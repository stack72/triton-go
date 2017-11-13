package compute_test

import (
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

const accountUrl = "testing"

func MockComputeClient() *compute.ComputeClient {
	return &compute.ComputeClient{
		Client: testutils.NewMockClient(testutils.MockClientInput{
			AccountName: accountUrl,
		}),
	}
}
