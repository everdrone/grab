package instance

import (
	"path/filepath"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"
	"github.com/mitchellh/go-homedir"

	"github.com/hashicorp/hcl/v2"
)

func (s *Grab) ParseFlags() *hcl.Diagnostics {
	flags := &FlagsState{}

	flags.Force, _ = s.Command.Flags().GetBool("force")
	flags.Quiet, _ = s.Command.Flags().GetBool("quiet")
	flags.Strict, _ = s.Command.Flags().GetBool("strict")
	flags.DryRun, _ = s.Command.Flags().GetBool("dry-run")
	flags.Progress, _ = s.Command.Flags().GetBool("progress")
	flags.Verbosity, _ = s.Command.Flags().GetCount("verbose")
	flags.ConfigPath, _ = s.Command.Flags().GetString("config")

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
	if !filepath.IsAbs(s.Config.Global.Location) {
		expanded, err := homedir.Expand(s.Config.Global.Location)
		if err != nil {
			return &hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Could not expand home directory",
				Detail:   err.Error(),
			}}
		}
		if filepath.IsAbs(expanded) {
			s.Config.Global.Location = expanded
		} else {
			s.Config.Global.Location = filepath.Join(utils.Wd, s.Config.Global.Location)
		}
	}

	return &hcl.Diagnostics{}
}

func (s *Grab) ParseURLs(args []string) *hcl.Diagnostics {
	// gather urls from positional args
	args = utils.Unique(args)

	urls, diags := utils.GetURLsFromArgs(args)
	if diags.HasErrors() {
		return &diags
	}

	s.URLs = utils.Unique(urls)

	return &hcl.Diagnostics{}
}
