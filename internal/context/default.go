package context

import (
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func GetEnvironmentMap() map[string]cty.Value {
	result := make(map[string]cty.Value)

	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		result[parts[0]] = cty.StringVal(parts[1])
	}

	return result
}

func BuildInitialContext() *hcl.EvalContext {
	result := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: map[string]function.Function{},
	}

	result.Variables["env"] = cty.ObjectVal(GetEnvironmentMap())

	result.Variables["os"] = cty.ObjectVal(map[string]cty.Value{
		"name": cty.StringVal(runtime.GOOS),
		"arch": cty.StringVal(runtime.GOARCH),
	})

	result.Variables["body"] = cty.StringVal("body")
	result.Variables["url"] = cty.StringVal("url")

	return result
}
