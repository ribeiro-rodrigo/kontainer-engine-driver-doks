package digitalocean

import "github.com/rancher/kontainer-engine/types"

func getCreateOptions() *types.DriverFlags {

	builder := flagBuilder()

	builder(
		"display-name",
		types.StringType,
		"The displayed name of the cluster in the Rancher UI",
		nil,
	)

	builder(
		"cluster-name",
		types.StringType,
		"Default name to something meaningful to you",
		nil,
	)

	builder(
		"auto-upgraded",
		types.BoolType,
		"Automatically updates Kubernetes version",
		&types.Default{
			DefaultBool: false,
		},
	)

	builder(
		"region-slug",
		types.StringType,
		"Your Kubernetes cluster will be located in a single datacenter.",
		&types.Default{
			DefaultString: "nyc3",
		},
	)

	builder(
		"version-slug",
		types.StringType,
		"Kubernetes version",
		nil,
	)

	builder(
		"node-pool-name",
		types.StringType,
		"Name of pool of worker nodes",
		nil,
	)

	builder(
		"node-pool-autoscale",
		types.BoolType,
		"Enables auto scaling group",
		&types.Default{
			DefaultBool: false,
		},
	)

	builder(
		"node-count",
		types.IntType,
		"The desired number of worker nodes",
		&types.Default{
			DefaultInt: 2,
		},
	)

	builder(
		"node-min",
		types.IntType,
		"The minimum number of worker nodes",
		nil,
	)

	builder(
		"node-max",
		types.IntType,
		"The maximum number of worker nodes",
		nil,
	)

	builder(
		"node-type",
		types.StringType,
		"The type of machine to use for worker nodes",
		&types.Default{
			DefaultString: "s-2vcpu-2gb",
		},
	)

	builder(
		"tags",
		types.StringSliceType,
		"Optional tags to your cluster",
		nil,
	)

	return builder(
		"vpc-id",
		types.StringType,
		"Cluster VPC",
		nil,
	)
}

func flagBuilder() func(name, typ, usage string, def *types.Default) *types.DriverFlags {
	flags := &types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}

	return func(name, typ, usage string, def *types.Default) *types.DriverFlags {
		flags.Options[name] = &types.Flag{
			Type:    typ,
			Usage:   usage,
			Default: def,
		}

		return flags
	}
}

func makeFlag(name, typ, usage string, def *types.Default) *types.Flag {
	return &types.Flag{
		Type:    typ,
		Usage:   usage,
		Default: def,
	}
}
