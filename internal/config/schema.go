package config

import "regexp"

type Config struct {
	Global GlobalConfig `hcl:"global,block"`
	Sites  []SiteConfig `hcl:"site,block"`
}

type GlobalConfig struct {
	Location string             `hcl:"location"`
	Network  *RootNetworkConfig `hcl:"network,block"`
}

type RootNetworkConfig struct {
	Timeout *int               `hcl:"timeout"`
	Retries *int               `hcl:"retries"`
	Headers *map[string]string `hcl:"headers"`
}

type SiteConfig struct {
	Name         string              `hcl:"name,label"`
	Test         string              `hcl:"test"`
	Network      *NetworkConfig      `hcl:"network,block"`
	Subdirectory *SubdirectoryConfig `hcl:"subdirectory,block"`
	Assets       []AssetConfig       `hcl:"asset,block"`
	Infos        []InfoConfig        `hcl:"info,block"`
	// computed
	URLs       []string
	InfoMap    InfoCacheMap // location -> info -> value
	HasMatches bool
}

type SubdirectoryConfig struct {
	Pattern string `hcl:"pattern"`
	Capture string `hcl:"capture"`
	From    string `hcl:"from"`
}

type AssetConfig struct {
	Name       string            `hcl:"name,label"`
	Pattern    string            `hcl:"pattern"`
	Capture    string            `hcl:"capture"`
	FindAll    *bool             `hcl:"find_all"`
	Network    *NetworkConfig    `hcl:"network,block"`
	Transforms []TransformConfig `hcl:"transform,block"`
	// computed
	Downloads map[string]string
}

type InfoConfig struct {
	Name    string `hcl:"name,label"`
	Pattern string `hcl:"pattern"`
	Capture string `hcl:"capture"`
}

type NetworkConfig struct {
	Inherit *bool              `hcl:"inherit"`
	Timeout *int               `hcl:"timeout"`
	Retries *int               `hcl:"retries"`
	Headers *map[string]string `hcl:"headers"`
}

type TransformConfig struct {
	Name    string `hcl:"name,label"`
	Pattern string `hcl:"pattern"`
	Replace string `hcl:"replace"`
}

type RegexCacheMap map[string]*regexp.Regexp

type InfoCacheMap map[string]map[string]string
