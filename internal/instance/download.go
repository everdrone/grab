package instance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/everdrone/grab/internal/net"
	"github.com/everdrone/grab/internal/utils"

	"github.com/hashicorp/hcl/v2"
)

func (s *Grab) Download() *hcl.Diagnostics {
	if s.Flags.DryRun {
		for _, site := range s.Config.Sites {
			// do not print anything if the site does not download anything?
			// if !site.HasMatches {
			// 	continue
			// }

			s.Command.Printf("site: %s\n", site.Name)

			for location, infoMap := range site.InfoMap {
				s.Command.Printf("info: %s:\n", location)
				s.Command.Print(utils.FormatMap(infoMap, " : ", true))
			}

			for _, asset := range site.Assets {
				relativeDownloads := make(map[string]string, 0)
				for src, dst := range asset.Downloads {
					rel, _ := filepath.Rel(s.Config.Global.Location, dst)
					relativeDownloads[src] = rel
				}

				s.Command.Printf("asset \"%s\": %s\n", asset.Name, utils.Plural(len(asset.Downloads), "asset", "assets"))
				s.Command.Print(utils.FormatMap(relativeDownloads, " → ", false))
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

			s.Log(2, hcl.Diagnostics{{
				Severity: utils.DiagInfo,
				Summary:  "Indexing info",
				Detail:   dst,
			}})

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

				// check if file exists
				performWrite := true
				if exists, err := utils.AFS.Exists(dst); err != nil || exists {
					performWrite = false
				}

				// if force or file does not exist, write to disk
				if s.Flags.Force || performWrite {
					s.Log(2, hcl.Diagnostics{{
						Severity: utils.DiagInfo,
						Summary:  "Downloading asset",
						Detail:   fmt.Sprintf("%s -> %s", src, dst),
					}})

					if err := net.Download(src, dst, options); err != nil {
						// FIXME: this should be appended and returned when the function finishes if not in strict mode
						diags := &hcl.Diagnostics{{
							Severity: hcl.DiagError,
							Summary:  "Failed to download asset",
							Detail:   fmt.Sprintf("%s: %s", src, err.Error()),
						}}

						// return now if we are in strict mode
						if s.Flags.Strict {
							return diags
						} else {
							s.Log(1, *diags)
						}
					}
				} else {
					s.Log(1, hcl.Diagnostics{{
						Severity: hcl.DiagWarning,
						Summary:  "File already exists",
						Detail:   fmt.Sprintf("Skipping download of %s.%s into %s", site.Name, asset.Name, dst),
					}})
				}
			}
		}
	}

	return &hcl.Diagnostics{}
}
