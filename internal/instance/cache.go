package instance

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/everdrone/grab/internal/config"
	"github.com/everdrone/grab/internal/net"
	"github.com/everdrone/grab/internal/utils"
	"github.com/hashicorp/hcl/v2"
)

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

func removePathFromURL(str string) (*url.URL, error) {
	base, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	// get only the base (scheme://host)
	base.Path = ""
	base.RawPath = ""
	base.RawQuery = ""
	base.Fragment = ""
	base.RawFragment = ""

	return base, nil
}

func (s *Grab) BuildAssetCache() *hcl.Diagnostics {
	var diags *hcl.Diagnostics

	for siteIndex, site := range s.Config.Sites {
		for _, pageUrl := range site.URLs {
			// we already checked this url before, so we can skip the error
			base, _ := removePathFromURL(pageUrl)

			options := net.MergeFetchOptionsChain(s.Config.Global.Network, site.Network)

			s.Log(2, hcl.Diagnostics{{
				Severity: utils.DiagInfo,
				Summary:  "Fetching page",
				Detail:   pageUrl,
			}})

			// MARK: - get the page body

			body, err := net.Fetch(pageUrl, options)
			if err != nil {
				diags = &hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch page",
					Detail:   fmt.Sprintf("%s: %s", pageUrl, err.Error()),
				}}

				// if we are in strict mode we need to return immediately
				if s.Flags.Strict {
					return diags
				} else {
					continue
				}
			}

			// MARK: - get the destination path (subdirectory)

			var subdirectory string
			if site.Subdirectory != nil {
				// we have a subdirectory block

				var source string
				if site.Subdirectory.From == "url" {
					source = pageUrl
				} else {
					source = body
				}

				subDirs, err := utils.GetCaptures(s.RegexCache[site.Subdirectory.Pattern], false, site.Subdirectory.Capture, source)
				if err != nil {
					return &hcl.Diagnostics{{
						Severity: hcl.DiagError,
						Summary:  "Failed to get subdirectory",
						Detail:   err.Error(),
					}}
				}

				if len(subDirs) > 0 {
					// do not append if the path is absolute
					if filepath.IsAbs(subDirs[0]) {
						subdirectory = subDirs[0]
					} else {
						subdirectory = filepath.Join(s.Config.Global.Location, site.Name, subDirs[0])
					}
				}
			} else {
				// we have no subdirectory block, just use the site name
				subdirectory = filepath.Join(s.Config.Global.Location, site.Name)
			}

			// MARK: - loop through the asset blocks

			for assetIndex, asset := range site.Assets {
				// match against body
				if s.RegexCache[asset.Pattern].MatchString(body) {
					findAll := false
					if asset.FindAll != nil {
						findAll = *asset.FindAll
					}

					// get capture groups
					captures, err := utils.GetCaptures(s.RegexCache[asset.Pattern], findAll, asset.Capture, body)
					if err != nil {
						return &hcl.Diagnostics{{
							Severity: hcl.DiagError,
							Summary:  "Failed to get captures",
							Detail:   fmt.Sprintf("%s: %s", pageUrl, err.Error()),
						}}
					}

					// remove duplicates
					captures = utils.Unique(captures)

					// MARK: - Transform url

					// TODO: we should change the config schema to store transforms as a map
					// where the key is the transform label, so we don't end up looping through an array
					transformUrl := utils.Filter(asset.Transforms, func(t config.TransformConfig) bool {
						return t.Name == "url"
					})

					if len(transformUrl) > 0 {
						// we have a transform url block
						t := transformUrl[0]
						for i, src := range captures {
							captures[i] = s.RegexCache[t.Pattern].ReplaceAllString(src, t.Replace)
						}
					}

					// MARK: - Transform destination filename

					transformFilename := utils.Filter(asset.Transforms, func(t config.TransformConfig) bool {
						return t.Name == "filename"
					})

					destinations := make(map[string]string, 0)

					if len(transformFilename) > 0 {
						// we have a transform filename block
						t := transformFilename[0]
						for _, src := range captures {
							fileName := s.RegexCache[t.Pattern].ReplaceAllString(src, t.Replace)

							// NOTE: the result of "transform filename" could be an absolute path!
							//       so we should not append if absolute
							if filepath.IsAbs(fileName) {
								destinations[src] = fileName
							} else {
								destinations[src] = filepath.Join(subdirectory, fileName)
							}

							// unescape the filename to write on disk
							unescaped, err := url.QueryUnescape(destinations[src])
							if err != nil {
								return &hcl.Diagnostics{{
									Severity: hcl.DiagError,
									Summary:  "Failed to unescape filename",
									Detail:   fmt.Sprintf("%s: %s", fileName, err.Error()),
								}}
							}

							destinations[src] = unescaped
						}
					} else {
						// we don't have any transform filename blocks
						for _, src := range captures {
							// simply get the filename from the url path
							fileName := filepath.Base(src)
							destinations[src] = filepath.Join(subdirectory, fileName)

							// unescape the filename to write on disk
							unescaped, err := url.QueryUnescape(destinations[src])
							if err != nil {
								return &hcl.Diagnostics{{
									Severity: hcl.DiagError,
									Summary:  "Failed to unescape filename",
									Detail:   fmt.Sprintf("%s: %s", fileName, err.Error()),
								}}
							}

							destinations[src] = unescaped
						}
					}

					// MARK: - loop through the map to check for relative urls

					resolvedDestinations := make(map[string]string, 0)
					for src, dst := range destinations {
						parsed, err := url.Parse(src)
						if err != nil {
							return &hcl.Diagnostics{{
								Severity: hcl.DiagError,
								Summary:  "Failed to parse url",
								Detail:   fmt.Sprintf("%s: %s", src, err.Error()),
							}}
						}

						// if path is still relative, append it to the scheme://domain.name of the page
						if !parsed.IsAbs() {
							resolved, err := base.Parse(src)
							if err != nil {
								return &hcl.Diagnostics{{
									Severity: hcl.DiagError,
									Summary:  "Failed to resolve relative url",
									Detail:   fmt.Sprintf("%s: %s", src, err.Error()),
								}}
							}

							s.Log(3, hcl.Diagnostics{{
								Severity: utils.DiagInfo,
								Summary:  "Resolved relative url",
								Detail:   fmt.Sprintf("%s -> %s", src, resolved.String()),
							}})

							resolvedDestinations[resolved.String()] = dst
						} else {
							// nothing to do, the url is already absolute
							resolvedDestinations[src] = dst
						}
					}

					// initialize the map if nil
					if s.Config.Sites[siteIndex].Assets[assetIndex].Downloads == nil {
						s.Config.Sites[siteIndex].Assets[assetIndex].Downloads = make(map[string]string, 0)
					}

					// add the destinations to the asset
					for src, dst := range resolvedDestinations {
						s.Config.Sites[siteIndex].Assets[assetIndex].Downloads[src] = dst
					}

					// is this site going to perform downloads?
					// if len(resolvedDestinations) > 0 {
					// 	s.Config.Sites[siteIndex].HasMatches = true
					// }

					s.TotalAssets += int64(len(resolvedDestinations))
				}
			}

			// MARK: - Indexing

			// store the url and the timestamp by default
			infoMap := make(map[string]string, 0)
			infoMap["url"] = pageUrl
			infoMap["timestamp"] = time.Now().UTC().Format(time.RFC3339Nano)

			// loop through index blocks
			for _, info := range site.Infos {
				key := info.Name

				if s.RegexCache[info.Pattern].MatchString(body) {
					captures, err := utils.GetCaptures(s.RegexCache[info.Pattern], false, info.Capture, body)
					if err != nil {
						return &hcl.Diagnostics{{
							Severity: hcl.DiagError,
							Summary:  "Failed to get capture",
							Detail:   fmt.Sprintf("%s: %s", pageUrl, err.Error()),
						}}
					}

					if len(captures) > 0 {
						infoMap[key] = captures[0]
					}
				}
			}

			if s.Config.Sites[siteIndex].InfoMap == nil {
				s.Config.Sites[siteIndex].InfoMap = make(config.InfoCacheMap, 0)
			}

			s.Config.Sites[siteIndex].InfoMap[subdirectory] = infoMap
		}
	}

	return &hcl.Diagnostics{}
}
