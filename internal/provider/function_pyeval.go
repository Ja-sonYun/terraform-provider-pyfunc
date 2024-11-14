// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var pyEvalReturnAttrTypes = map[string]attr.Type{
	"stdout": types.StringType,
}

var _ function.Function = &PyEvalFunction{}

type PyEvalFunction struct{}

func NewPyEvalFunction() function.Function {
	return &PyEvalFunction{}
}

func (f *PyEvalFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "pyeval"
}

func (f *PyEvalFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Run a Python script",
		Description: "Given a Python script, will run the script and return the stdout output.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "script",
				Description: "Python script to run",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: pyEvalReturnAttrTypes,
		},
	}
}

func (f *PyEvalFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var script string

	resp.Error = req.Arguments.Get(ctx, &script)
	if resp.Error != nil {
		return
	}

	// Run the Python script
	cmd := exec.Command("python", "-c", script)
	stdout, err := cmd.Output()
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed to run Python script, underlying error: %s", err.Error()))
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error running Python script: %q", script))
		return
	}

	pyEvalObj, diags := types.ObjectValue(
		pyEvalReturnAttrTypes,
		map[string]attr.Value{
			"stdout": types.StringValue(string(stdout)),
		},
	)

	resp.Error = function.FuncErrorFromDiags(ctx, diags)
	if resp.Error != nil {
		return
	}

	resp.Error = resp.Result.Set(ctx, &pyEvalObj)
}
