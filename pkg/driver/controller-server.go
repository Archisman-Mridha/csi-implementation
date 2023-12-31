package driver

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/digitalocean/godo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ControllerGetCapabilities allows the CO to check the supported capabilities of controller service
// provided by this Plugin.
func (c *CsiDriver) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	supportedCapabilityTypes := []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	}

	var supportedCapabilities []*csi.ControllerServiceCapability
	for _, capabilityType := range supportedCapabilityTypes {
		supportedCapabilities = append(supportedCapabilities, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: capabilityType,
				},
			},
		})
	}

	response := &csi.ControllerGetCapabilitiesResponse{
		Capabilities: supportedCapabilities,
	}
	return response, nil
}

// CreateVolume will be invoked by the CO to provision a new volume on behalf of a user (to be
// consumed as either a block device or a mounted filesystem). This operation MUST be idempotent.
func (c *CsiDriver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	fmt.Println("CreateVolume invoked")

	volumeName := req.Name
	if volumeName == "" {
		return nil, status.Error(codes.InvalidArgument, "Volume name must be specified")
	}

	minStorageCapacity := req.CapacityRange.RequiredBytes

	// volumeCapabilities consists of volume access modes and the volume binding mode.
	volumeCapabilities := req.VolumeCapabilities
	if len(volumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume access modes must be specified")
	}

	volume, _, err := c.doStorageService.CreateVolume(context.Background(), &godo.VolumeCreateRequest{
		Name:          volumeName,
		Region:        c.region,
		SizeGigaBytes: minStorageCapacity / (1024 * 1024 * 1024),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed provisioning volume : %v", err))
	}

	response := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volume.ID,
			CapacityBytes: volume.SizeGigaBytes * (1024 * 1024 * 1024),
		},
	}
	return response, nil
}

func (c *CsiDriver) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, nil
}

func (c *CsiDriver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, nil
}

func (c *CsiDriver) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, nil
}

func (c *CsiDriver) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, nil
}

func (c *CsiDriver) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, nil
}

func (c *CsiDriver) ControllerModifyVolume(ctx context.Context, req *csi.ControllerModifyVolumeRequest) (*csi.ControllerModifyVolumeResponse, error) {
	return nil, nil
}
