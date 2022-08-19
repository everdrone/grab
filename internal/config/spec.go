package config

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var ConfigSpec = &hcldec.ObjectSpec{
	"global": &hcldec.BlockSpec{
		TypeName: "global",
		Required: true,
		Nested:   GlobalSpec,
	},
	"sites": &hcldec.BlockTupleSpec{
		TypeName: "site",
		MinItems: 1,
		Nested:   SiteSpec,
	},
}

var GlobalSpec = &hcldec.ObjectSpec{
	"location": &hcldec.AttrSpec{
		Name:     "location",
		Required: true,
		Type:     cty.String,
	},
	"network": &hcldec.BlockSpec{
		TypeName: "network",
		Required: false,
		Nested:   RootNetworkSpec,
	},
}

var RootNetworkSpec = &hcldec.ObjectSpec{
	"timeout": &hcldec.AttrSpec{
		Name:     "timeout",
		Required: false,
		Type:     cty.Number,
	},
	"retries": &hcldec.AttrSpec{
		Name:     "retries",
		Required: false,
		Type:     cty.Number,
	},
	"headers": &hcldec.AttrSpec{
		Name:     "headers",
		Required: false,
		Type:     cty.Map(cty.String),
	},
}

var SiteSpec = &hcldec.ObjectSpec{
	"name": &hcldec.BlockLabelSpec{
		Index: 0,
		Name:  "name",
	},
	"test": &hcldec.AttrSpec{
		Name:     "test",
		Type:     cty.String,
		Required: true,
	},
	"network": &hcldec.BlockSpec{
		TypeName: "network",
		Required: false,
		Nested:   NetworkSpec,
	},
	// NOTE: must have at least one "asset" or at least one "info"
	"assets": &hcldec.BlockTupleSpec{
		TypeName: "asset",
		MinItems: 0,
		Nested:   AssetSpec,
	},
	"infos": &hcldec.BlockTupleSpec{
		TypeName: "info",
		MinItems: 0,
		Nested:   InfoSpec,
	},
	"subdirectory": &hcldec.BlockSpec{
		TypeName: "subdirectory",
		Required: false,
		Nested:   SubdirectorySpec,
	},
}

var NetworkSpec = &hcldec.ObjectSpec{
	"inherit": &hcldec.AttrSpec{
		Name:     "inherit",
		Required: false,
		Type:     cty.Bool,
	},
	"timeout": &hcldec.AttrSpec{
		Name:     "timeout",
		Required: false,
		Type:     cty.Number,
	},
	"retries": &hcldec.AttrSpec{
		Name:     "retries",
		Required: false,
		Type:     cty.Number,
	},
	"headers": &hcldec.AttrSpec{
		Name:     "headers",
		Required: false,
		Type:     cty.Map(cty.String),
	},
}

var AssetSpec = &hcldec.ObjectSpec{
	"name": &hcldec.BlockLabelSpec{
		Index: 0,
		Name:  "name",
	},
	"pattern": &hcldec.AttrSpec{
		Name:     "pattern",
		Type:     cty.String,
		Required: true,
	},
	"capture": &hcldec.AttrSpec{
		Name:     "capture",
		Type:     cty.String,
		Required: true,
	},
	"find_all": &hcldec.AttrSpec{
		Name:     "find_all",
		Type:     cty.Bool,
		Required: false,
	},
	"network": &hcldec.BlockSpec{
		TypeName: "network",
		Required: false,
		Nested:   NetworkSpec,
	},
	"transforms": &hcldec.BlockTupleSpec{
		TypeName: "transform",
		MinItems: 0,
		MaxItems: 2,
		Nested:   TransformSpec,
	},
	// TODO: allow setting a subdirectory for the asset
}

var InfoSpec = &hcldec.ObjectSpec{
	"name": &hcldec.BlockLabelSpec{
		Index: 0,
		Name:  "name",
	},
	"pattern": &hcldec.AttrSpec{
		Name:     "pattern",
		Type:     cty.String,
		Required: true,
	},
	"capture": &hcldec.AttrSpec{
		Name:     "capture",
		Type:     cty.String,
		Required: true,
	},
	// TODO: allow using find_all to save arrays of strings
}

// must validate that there is only one "url" and only one "filename"
var TransformSpec = &hcldec.ObjectSpec{
	"name": &hcldec.BlockLabelSpec{
		Index: 0,
		Name:  "name",
	},
	"pattern": &hcldec.AttrSpec{
		Name:     "pattern",
		Type:     cty.String,
		Required: true,
	},
	"replace": &hcldec.AttrSpec{
		Name:     "replace",
		Type:     cty.String,
		Required: true,
	},
}

var SubdirectorySpec = &hcldec.ObjectSpec{
	"pattern": &hcldec.AttrSpec{
		Name:     "pattern",
		Type:     cty.String,
		Required: true,
	},
	"capture": &hcldec.AttrSpec{
		Name:     "capture",
		Type:     cty.String,
		Required: true,
	},
	"from": &hcldec.AttrSpec{
		Name:     "from",
		Type:     cty.String,
		Required: true,
	},
}
