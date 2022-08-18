package instance

import (
	"path/filepath"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"

	"github.com/hashicorp/hcl/v2"
)

func (s *Grab) ParseFlags() *hcl.Diagnostics {
	var err error
	flags := &FlagsState{}

	flags.Force, err = s.Command.Flags().GetBool("force")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `force` flag",
			Detail:   err.Error(),
		}}
	}

	flags.Quiet, err = s.Command.Flags().GetBool("quiet")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `quiet` flag",
			Detail:   err.Error(),
		}}
	}

	flags.Strict, err = s.Command.Flags().GetBool("strict")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `strict` flag",
			Detail:   err.Error(),
		}}
	}

	flags.DryRun, err = s.Command.Flags().GetBool("dry-run")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `dry-run` flag",
			Detail:   err.Error(),
		}}
	}

	flags.Progress, err = s.Command.Flags().GetBool("progress")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `progress` flag",
			Detail:   err.Error(),
		}}
	}

	flags.Verbosity, err = s.Command.Flags().GetCount("verbose")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `verbose` flag",
			Detail:   err.Error(),
		}}
	}

	flags.ConfigPath, err = s.Command.Flags().GetString("config")
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the `config` flag",
			Detail:   err.Error(),
		}}
	}

	// if both quiet and verbose are set, quiet wins
	if flags.Quiet {
		flags.Verbosity = 0
	} else {
		// make verbosity 1 indexed (just add one)
		flags.Verbosity++
	}

	s.Flags = flags
	return &hcl.Diagnostics{}
}

func (s *Grab) ParseConfig() *hcl.Diagnostics {
	if s.Flags.ConfigPath == "" {
		resolved, err := config.Resolve("grab.hcl", utils.Wd)
		if err != nil {
			return &hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Could not resolve config file",
				Detail:   err.Error(),
			}}
		}

		s.Flags.ConfigPath = resolved
	}

	// read file contents of config file
	fc, err := utils.AFS.ReadFile(s.Flags.ConfigPath)
	if err != nil {
		return &hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Could not read config file",
			Detail:   err.Error(),
		}}
	}

	// parse config and get regexCache
	config, _, regexCache, diags := config.Parse(fc, s.Flags.ConfigPath)
	if diags.HasErrors() {
		return &diags
	}

	s.Config = config
	s.RegexCache = regexCache

	// get global location
	if filepath.IsAbs(s.Config.Global.Location) {
		s.GlobalLocation = s.Config.Global.Location
	} else {
		s.GlobalLocation = filepath.Join(utils.Wd, s.Config.Global.Location)
	}

	return &hcl.Diagnostics{}
}

func (s *Grab) ParseURLs(args []string) *hcl.Diagnostics {
	// gather urls from positional args
	urls, diags := utils.GetURLsFromArgs(args)
	if diags.HasErrors() {
		return &diags
	}

	s.URLs = urls
	return &hcl.Diagnostics{}
}

func (s *Grab) BuildSiteCache() {
	for _, url := range s.URLs {
		for i, site := range s.Config.Sites {
			if s.RegexCache[site.Test].MatchString(url) {
				if s.Config.Sites[i].URLs == nil {
					s.Config.Sites[i].URLs = make([]string, 0)
				}

				s.Config.Sites[i].URLs = append(s.Config.Sites[i].URLs, url)
				break
			}
		}
	}
}
