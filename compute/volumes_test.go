package compute_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

var (
	getVolumeErrorType    = errors.New("Error executing Get request:")
	deleteVolumeErrorType = errors.New("Error executing Delete request:")
	createVolumeErrorType = errors.New("Error executing Create request:")
	updateVolumeErrorType = errors.New("Error executing Update request:")
)

func TestListVolume(t *testing.T) {
	computeClient := MockComputeClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) ([]*compute.Volume, error) {
		defer testutils.DeactivateClient()

		volumes, err := cc.Volumes().List(ctx, &compute.ListVolumeInput{
			State: "running",
		})
		if err != nil {
			return nil, err
		}
		return volumes, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/volumes", accountUrl), listVolumeSuccess)

		resp, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})
}

func TestGetVolume(t *testing.T) {
	computeClient := MockComputeClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) (*compute.Volume, error) {
		defer testutils.DeactivateClient()

		user, err := cc.Volumes().Get(ctx, &compute.GetVolumeInput{
			ID: "e435d72a-2498-8d49-a042-87b222a8b63f",
		})
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), getVolumeSuccess)

		resp, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), getVolumeEmpty)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), getVolumeBadDecode)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), getVolumeError)

		resp, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "Error executing Get request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestDeleteVolume(t *testing.T) {
	computeClient := MockComputeClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) error {
		defer testutils.DeactivateClient()

		return cc.Volumes().Delete(ctx, &compute.DeleteVolumeInput{
			ID: "e435d72a-2498-8d49-a042-87b222a8b63f",
		})
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("DELETE", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), deleteVolumeSuccess)

		err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("DELETE", fmt.Sprintf("/%s/volumes/%s", accountUrl, "1e435d72a-2498-8d49-a042-87b222a8b63f"), deleteVolumeError)

		err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "Error executing Delete request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestCreateVolume(t *testing.T) {
	computeClient := MockComputeClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) (*compute.Volume, error) {
		defer testutils.DeactivateClient()

		volume, err := cc.Volumes().Create(ctx, &compute.CreateVolumeInput{
			Name:     "test-volume",
			Size:     1000000,
			Type:     "tritonnfs",
			Networks: []string{"1537d72a-949a-2d89-7049-17b2f2a8b634"},
		})
		if err != nil {
			return nil, err
		}
		return volume, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/volumes", accountUrl), createVolumeSuccess)

		_, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/volumes", accountUrl), createVolumeError)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "Error executing Create request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestUpdateVolume(t *testing.T) {
	computeClient := MockComputeClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) error {
		defer testutils.DeactivateClient()

		return cc.Volumes().Update(ctx, &compute.UpdateVolumeInput{
			ID:   "e435d72a-2498-8d49-a042-87b222a8b63f",
			Name: "updated-name",
		})
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), updateVolumeSuccess)

		err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/volumes/%s", accountUrl, "e435d72a-2498-8d49-a042-87b222a8b63f"), updateVolumeError)

		err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "Error executing Update request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func listVolumeSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[{
    "id": "e435d72a-2498-8d49-a042-87b222a8b45d",
    "name": "my-volume",
    "owner_uuid": "ae35672a-9498-ed41-b017-82b221a8c63f",
    "type": "tritonnfs",
    "state": "ready",
    "networks": [
      "1537d72a-949a-2d89-7049-17b2f2a8b634"
    ]
  },
  {
    "id": "e435d72a-2498-8d49-a042-87b222a8b63f",
    "name": "my-volume2",
    "owner_uuid": "ae35672a-9498-ed41-b017-82b221a8c63f",
    "type": "tritonnfs",
    "state": "ready",
    "networks": [
      "1537d72a-949a-2d89-7049-17b2f2a8b634"
    ]
  }
]`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getVolumeSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "id": "e435d72a-2498-8d49-a042-87b222a8b63f",
    "name": "my-volume",
    "owner_uuid": "ae35672a-9498-ed41-b017-82b221a8c63f",
    "type": "tritonnfs",
    "state": "ready",
    "networks": [
      "1537d72a-949a-2d89-7049-17b2f2a8b634"
    ]
  }
`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getVolumeEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func getVolumeBadDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "id": "e435d72a-2498-8d49-a042-87b222a8b63f",
    "name": "my-volume",
    "owner_uuid": "ae35672a-9498-ed41-b017-82b221a8c63f",
    "type": "tritonnfs",
    "state": "ready",}`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getVolumeError(req *http.Request) (*http.Response, error) {
	return nil, getVolumeErrorType
}

func deleteVolumeSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 204,
		Header:     header,
	}, nil
}

func deleteVolumeError(req *http.Request) (*http.Response, error) {
	return nil, deleteVolumeErrorType
}

func createVolumeSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "id": "e435d72a-2498-8d49-a042-87b222a8b63f",
    "name": "my-volume",
    "owner_uuid": "ae35672a-9498-ed41-b017-82b221a8c63f",
    "type": "tritonnfs",
    "state": "ready",
    "networks": [
      "1537d72a-949a-2d89-7049-17b2f2a8b634"
    ]
  }
`)

	return &http.Response{
		StatusCode: 201,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func createVolumeError(req *http.Request) (*http.Response, error) {
	return nil, createVolumeErrorType
}

func updateVolumeSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 204,
		Header:     header,
	}, nil
}

func updateVolumeError(req *http.Request) (*http.Response, error) {
	return nil, updateVolumeErrorType
}
