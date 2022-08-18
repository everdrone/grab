package net

import (
	"github.com/everdrone/grab/internal/config"
)

type FetchOptions struct {
	Headers map[string]string
	Timeout int
	Retries int
}

func ClampDefaults(options *FetchOptions) {
	// TODO: We should disable the timeout if <= 0
	if options.Timeout <= 0 {
		options.Timeout = 3000
	}

	if options.Retries <= 0 {
		options.Retries = 1
	}

	if options.Headers == nil {
		options.Headers = make(map[string]string, 0)
	}
}

func MergeFetchOptionsChain(root *config.RootNetworkConfig, configs ...*config.NetworkConfig) *FetchOptions {
	options := &FetchOptions{}

	if root != nil {
		if root.Timeout != nil {
			options.Timeout = *root.Timeout
		}

		if root.Retries != nil {
			options.Retries = *root.Retries
		}

		if root.Headers != nil {
			options.Headers = *root.Headers
		}
	}

	for _, config := range configs {
		if config != nil {
			if config.Inherit != nil && !*config.Inherit {
				// if the config says not to inherit, reset the options object
				options = &FetchOptions{}
			}

			if config.Timeout != nil {
				options.Timeout = *config.Timeout
			}

			if config.Retries != nil {
				options.Retries = *config.Retries
			}

			if config.Headers != nil {
				if options.Headers == nil {
					options.Headers = make(map[string]string, 0)
				}
				for k, v := range *config.Headers {
					options.Headers[k] = v
				}
			}
		}
	}

	ClampDefaults(options)

	return options
}
