package instance

import (
	"path/filepath"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/hcl/v2"
)

func (s *Grab) ParseFlags() {
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

	switch flags.Verbosity {
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 4:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	s.Flags = flags
}

func (s *Grab) ParseConfig() *hcl.Diagnostics {
	if s.Flags.ConfigPath == "" {
		log.Trace().Msg("no config file specified, resolving")

		resolved, err := config.Resolve("grab.hcl", utils.Wd)
		if err != nil {
			return &hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Could not resolve config file",
				Detail:   err.Error(),
			}}
		}

		log.Debug().Str("path", resolved).Msg("using config file")

		s.Flags.ConfigPath = resolved
	}

	log.Trace().Str("path", s.Flags.ConfigPath).Msg("parsing config file")

	// read file contents of config file
	fc, err := utils.Io.ReadFile(utils.Fs, s.Flags.ConfigPath)
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
	log.Trace().Msg("parsing arguments")

	args = utils.Unique(args)

	urls, diags := utils.GetURLsFromArgs(args)
	if diags.HasErrors() {
		return &diags
	}

	s.URLs = utils.Unique(urls)

	log.Trace().Strs("urls", s.URLs).Msgf("found %d %s", len(s.URLs), utils.Plural(len(s.URLs), "url", "urls"))

	return &hcl.Diagnostics{}
}
