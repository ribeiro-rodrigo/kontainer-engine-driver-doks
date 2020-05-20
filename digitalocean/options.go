package digitalocean

import "github.com/rancher/kontainer-engine/types"

func getCreateOptions() (*types.DriverFlags, error) {

	flags := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}

	flags.Options["display-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The displayed name of the cluster in the Rancher UI",
	}

	flags.Options["cluster-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Default name to something meaningful to you",
	}

	flags.Options["auto-upgraded"] = &types.Flag{
		Type:  types.BoolType,
		Usage: "Automatically updates Kubernetes version",
	}

	flags.Options["region-slug"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Your Kubernetes cluster will be located in a single datacenter.",
	}

	flags.Options["version-slug"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Kubernetes version",
	}

	flags.Options["node-pool-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Name of pool of worker nodes",
	}
	flags.Options["node-pool-autoscale"] = &types.Flag{
		Type:  types.BoolType,
		Usage: "Enables auto scaling group",
	}
	flags.Options["node-count"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The desired number of worker nodes",
	}

	flags.Options["node-min"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The minimum number of worker nodes",
	}

	flags.Options["node-max"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The maximum number of worker nodes",
	}

	flags.Options["node-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The type of machine to use for worker nodes ",
	}

	flags.Options["tags"] = &types.Flag{
		Type:  types.StringSliceType,
		Usage: "Optional tags to your cluster",
	}

	flags.Options["vpc-id"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Cluster VPC",
	}

	return &flags, nil
}
