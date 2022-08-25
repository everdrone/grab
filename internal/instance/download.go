package instance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/everdrone/grab/internal/net"
	"github.com/everdrone/grab/internal/utils"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/hcl/v2"
)

func (s *Grab) Download() error {
	if s.Flags.DryRun {
		log.Warn().Msg("dry run, not downloading")

		for _, site := range s.Config.Sites {
			// do not print anything if the site does not download anything?
			// if !site.HasMatches {
			// 	continue
			// }

			for location, infoMap := range site.InfoMap {
				log.Info().Fields(infoMap).Str("site", site.Name).Str("location", location).Msg("indexing")
			}

			for _, asset := range site.Assets {
				for src, dst := range asset.Downloads {
					rel, _ := filepath.Rel(s.Config.Global.Location, dst)
					log.Info().Str("source", src).Str("destination", rel).Str("site", site.Name).Str("asset", asset.Name).Msg("downloading")
				}

			}
		}

		return &hcl.Diagnostics{}
	}

	for _, site := range s.Config.Sites {

		// MARK: - Download info file

		for subdirectory, infoMap := range site.InfoMap {
			// create directory
			if err := utils.AFS.MkdirAll(subdirectory, os.ModePerm); err != nil {
				return &hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to create directory",
					Detail:   fmt.Sprintf("%s: %s", subdirectory, err.Error()),
				}}
			}

			marshaled, err := json.MarshalIndent(infoMap, "", "  ")
			if err != nil {
				return &hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to marshal info",
					Detail:   fmt.Sprintf("%+v: %s", infoMap, err.Error()),
				}}
			}

			dst := filepath.Join(subdirectory, "_info.json")

			log.Info().Str("destination", dst).Msg("indexing")

			if err := utils.AFS.WriteFile(dst, marshaled, os.ModePerm); err != nil {
				return &hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to write info file",
					Detail:   fmt.Sprintf("%s: %s", dst, err.Error()),
				}}
			}
		}

		// MARK: - Download asset files

		for _, asset := range site.Assets {
			for src, dst := range asset.Downloads {
				// create directory
				dir := filepath.Dir(dst)
				if err := utils.AFS.MkdirAll(dir, os.ModePerm); err != nil {
					return &hcl.Diagnostics{{
						Severity: hcl.DiagError,
						Summary:  "Failed to create directory",
						Detail:   fmt.Sprintf("%s: %s", dir, err.Error()),
					}}
				}

				options := net.MergeFetchOptionsChain(s.Config.Global.Network, site.Network, asset.Network)

				log.Debug().Str("site", site.Name).Str("asset", asset.Name).Str("source", src).Interface("options", options).Msg("network options")

				// check if file exists
				performWrite := true
				if exists, err := utils.AFS.Exists(dst); err != nil || exists {
					performWrite = false
				}

				// if force or file does not exist, write to disk
				if s.Flags.Force || performWrite {
					log.Info().Str("source", src).Str("destination", strings.TrimPrefix(dst, s.Config.Global.Location)).Msg("downloading")

					if err := net.Download(src, dst, options); err != nil {
						// return now if we are in strict mode
						if s.Flags.Strict {
							log.Err(err).Str("source", src).Str("destination", strings.TrimPrefix(dst, s.Config.Global.Location)).Msg("failed to download asset")
							return err
						} else {
							log.Err(err).Str("source", src).Str("destination", strings.TrimPrefix(dst, s.Config.Global.Location)).Msg("failed to download asset")
						}
					}
				} else {
					log.Warn().Str("destination", strings.TrimPrefix(dst, s.Config.Global.Location)).Msg("file already exists")
				}
			}
		}
	}

	return nil
}
