package digitalocean

import "github.com/digitalocean/godo"

type State struct {
	Name        string
	Tags        []string
	AutoUpgrade bool
	RegionSlug  string
	VPCID       string
	VersionSlug string
	NodePool    *godo.KubernetesNodePoolCreateRequest
}
