package config

import (
	"fmt"
	"regexp"

	"github.com/everdrone/grab/internal/context"
	"github.com/everdrone/grab/internal/utils"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func ValidateSpec(body *hcl.Body, ctx *hcl.EvalContext) hcl.Diagnostics {
	_, diags := hcldec.Decode(*body, ConfigSpec, ctx)
	if diags.HasErrors() {
		return diags
	}

	return nil
}

func ValidateConfig(root *hclsyntax.Body, ctx *hcl.EvalContext) hcl.Diagnostics {
	var diags hcl.Diagnostics

	sites := utils.Filter(root.Blocks, func(b *hclsyntax.Block) bool {
		return b.Type == "site"
	})

	for _, site := range sites {
		// validate that there is at least one "asset" or at least one "info" block inside every "site" block
		assets := utils.Filter(site.Body.Blocks, func(b *hclsyntax.Block) bool {
			return b.Type == "asset"
		})
		infos := utils.Filter(site.Body.Blocks, func(b *hclsyntax.Block) bool {
			return b.Type == "info"
		})

		if len(assets) == 0 && len(infos) == 0 {
			return append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Insufficient \"site\" and \"info\" blocks",
				Detail:   "At least one asset or one info block must be defined inside a \"site\" block.",
				Subject:  &site.Body.SrcRange,
			})
		}

		for _, asset := range assets {
			// if "transform" blocks are present:
			//  - validate that the label is either "url" or "filename"
			//  - validate that there is not more than one "transform" block with the same label
			transforms := utils.Filter(asset.Body.Blocks, func(b *hclsyntax.Block) bool {
				return b.Type == "transform"
			})

			for _, transform := range transforms {
				label := transform.Labels[0]

				if label != "url" && label != "filename" {
					return append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid block label",
						Detail:   "\"transform\" block labels must be either \"url\" or \"filename\".",
						Subject:  &transform.LabelRanges[0],
					})
				}

				if len(utils.Filter(transforms, func(t *hclsyntax.Block) bool { return t.Labels[0] == label })) > 1 {
					return append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Duplicate block label",
						Detail:   fmt.Sprintf("No more than one \"transform\" block with the label \"%s\" is allowed.", label),
						Subject:  &transform.LabelRanges[0],
					})
				}
			}

		}

		// validate that, inside all "subdirectory" blocks, the "from" attribute is either "body" or "url"
		subdirectories := utils.Filter(site.Body.Blocks, func(b *hclsyntax.Block) bool {
			return b.Type == "subdirectory"
		})

		for _, subdirectory := range subdirectories {
			from := subdirectory.Body.Attributes["from"]

			val, moreDiags := from.Expr.Value(ctx)
			diags = append(diags, moreDiags...)

			if val != cty.StringVal("body") && val != cty.StringVal("url") {
				return append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid block attribute",
					Detail:   "The \"from\" attribute must be either \"body\" or \"url\".",
					Subject:  &from.EqualsRange,
				})
			}
		}
	}

	return nil
}

func EvaluateRegexPattern(attr *hclsyntax.Attribute, ctx *hcl.EvalContext) (string, *regexp.Regexp, hcl.Diagnostics) {
	val, diags := attr.Expr.Value(ctx)
	if diags.HasErrors() {
		return "", nil, diags
	}

	str := val.AsString()
	re, err := regexp.Compile(str)
	if err != nil {
		return "", nil, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid regex pattern",
			Detail:   err.Error(),
			Subject:  &attr.EqualsRange,
		})
	}

	return str, re, diags
}

// validate all regexp attributes
// - site*.test
// - site*.assets*.pattern
// - site*.assets*.transform*.pattern
// - site*.info*.pattern
// - site*.subdirectory.pattern
func BuildRegexCache(root *hclsyntax.Body, ctx *hcl.EvalContext) (RegexCacheMap, hcl.Diagnostics) {
	log.Trace().Msg("building regex cache")

	var regexCache = make(RegexCacheMap)

	sites := utils.Filter(root.Blocks, func(b *hclsyntax.Block) bool {
		return b.Type == "site"
	})

	for _, site := range sites {
		log.Trace().Str("name", site.Labels[0]).Msg("processing site")

		// test attribute is there
		if site.Body.Attributes["test"] != nil {
			str, re, diags := EvaluateRegexPattern(site.Body.Attributes["test"], ctx)
			if diags.HasErrors() {
				return nil, diags
			}

			log.Trace().Str("name", site.Labels[0]).Str("pattern", str).Msg("adding test regex")

			regexCache[str] = re
		}

		// gather all blocks that must contain the "pattern" attribute
		assets := utils.Filter(site.Body.Blocks, func(b *hclsyntax.Block) bool { return b.Type == "asset" })
		infos := utils.Filter(site.Body.Blocks, func(b *hclsyntax.Block) bool { return b.Type == "info" })
		subdirectories := utils.Filter(site.Body.Blocks, func(b *hclsyntax.Block) bool { return b.Type == "subdirectory" })

		transforms := make([]*hclsyntax.Block, 0)
		for _, asset := range assets {
			transforms = append(transforms, utils.Filter(asset.Body.Blocks, func(b *hclsyntax.Block) bool { return b.Type == "transform" })...)
		}

		patternBlocks := append(assets, infos...)
		patternBlocks = append(patternBlocks, subdirectories...)
		patternBlocks = append(patternBlocks, transforms...)

		for _, pb := range patternBlocks {
			if pb.Body.Attributes["pattern"] != nil {
				str, re, diags := EvaluateRegexPattern(pb.Body.Attributes["pattern"], ctx)
				if diags.HasErrors() {
					return nil, diags
				}

				log.Trace().Str("name", site.Labels[0]).Str("pattern", str).Msg("adding pattern regex")

				// FIXME: see if the regex has named captures or indexed captures and warn the user
				// if pb.Body.Attributes["capture"] != nil {
				// 	if utils.HasNamedCaptures(re) &&
				// }

				regexCache[str] = re
			}
		}
	}

	return regexCache, nil
}

func Parse(b []byte, filename string) (*Config, *hcl.EvalContext, RegexCacheMap, hcl.Diagnostics) {
	// parse
	p := hclparse.NewParser()
	file, diags := p.ParseHCL(b, filename)
	if diags.HasErrors() {
		return nil, nil, nil, diags
	}

	// create context
	ctx := context.BuildInitialContext()

	// validate against spec
	diags = ValidateSpec(&file.Body, ctx)
	if diags.HasErrors() {
		return nil, nil, nil, diags
	}

	root := file.Body.(*hclsyntax.Body)

	// validate what cannot be done in the spec
	diags = ValidateConfig(root, ctx)
	if diags.HasErrors() {
		return nil, nil, nil, diags
	}

	// validate regular expressions and build cache
	regexCache, diags := BuildRegexCache(root, ctx)
	if diags.HasErrors() {
		return nil, nil, nil, diags
	}

	// decode
	var config Config
	diags = gohcl.DecodeBody(file.Body, ctx, &config)
	if diags.HasErrors() {
		return nil, nil, nil, diags
	}

	return &config, ctx, regexCache, nil
}
