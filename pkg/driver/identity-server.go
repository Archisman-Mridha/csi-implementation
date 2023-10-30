package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Probe is invoked to verify that the plugin is in a healthy and ready state. If an unhealthy state
// is reported, via a non-success response, the CO may take action with the intent to bring the
// plugin to a healthy state. Such actions may include, but shall not be limited to, the following:
// (a) Restarting the plugin container or (b) Notifying the plugin supervisor.
func (c *CsiDriver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	response := &csi.ProbeResponse{Ready: wrapperspb.Bool(c.isReady)}
	return response, nil
}

func (c *CsiDriver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	response := &csi.GetPluginInfoResponse{Name: CSI_DRIVER_NAME}
	return response, nil
}

// GetPluginCapabilities allows the CO to query the supported capabilities of the Plugin as a whole:
// it is the grand sum of all capabilities of all instances of the Plugin software, as it is
// intended to be deployed. All instances of the same version (see vendor_version
// of GetPluginInfoResponse) of the Plugin shall return the same set of capabilities, regardless of
// both: (a) where instances are deployed on the cluster as well as; (b) which RPCs an instance is
// serving.
func (c *CsiDriver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	response := &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						// Indicates that the Plugin provides RPCs for the ControllerService.
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}
	return response, nil
}
