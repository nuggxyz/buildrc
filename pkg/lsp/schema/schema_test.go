package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/k0kubun/pp/v3"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func TestIntegrationJSONLoad(t *testing.T) {
	ctx := context.Background()
	pp.SetDefaultMaxDepth(5)

	// load schema file
	schema, err := LoadJsonSchemaFile(ctx, "https://raw.githubusercontent.com/SchemaStore/schemastore/master/src/schemas/json/github-workflow.json")
	if err != nil {
		if errd, ok := err.(hcl.Diagnostics); ok {
			for _, d := range errd {
				t.Log(d)
			}
		}
		t.Fatal(err)
	}

	spec, err := JsonScemaToHCLSpec(ctx, schema)
	if err != nil {
		t.Fatal(err)
	}

	pp.Println(spec)

}

const validHCL = `

schedule "main" {
	  cron = "0 0 * * *" # every day at midnight
}
export = {
  name = "test"

  schedule = schedule.main
}
`

func TestValidHCLDecoding(t *testing.T) {
	ctx := context.Background()
	// pp.SetDefaultMaxDepth(5)

	// load schema file
	s, err := LoadJsonSchemaFile(ctx, "https://raw.githubusercontent.com/SchemaStore/schemastore/master/src/schemas/json/github-workflow.json")
	if err != nil {
		if errd, ok := err.(hcl.Diagnostics); ok {
			for _, d := range errd {
				t.Log(d)
			}
		}
		t.Fatal(err)
	}

	// spec, err := JsonScemaToHCLSpec(ctx, schema)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	file, diags := hclparse.NewParser().ParseHCL([]byte(validHCL), "test.hcl")
	if diags.HasErrors() {
		for _, d := range diags {
			t.Log(d)
		}
		t.Fatal(diags)
	}
	// pp.Println(spec)

	scheme := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "export",
				Required: true,
			},
			{
				Name:     "schedule",
				Required: false,
			},
		},
	}
	ctn, diags := file.Body.Content(scheme)
	if diags.HasErrors() {
		for _, d := range diags {
			t.Log(d)
		}
		t.Fatal(diags)
	}
	var attr hcl.Attribute
	for _, a := range ctn.Attributes {
		if a.Name == "export" {
			attr = *a
		}
	}

	pp.Println(attr)

	Transform(ctx, s, attr)

	// val, diags := hcldec.Decode(file.Body, spec, nil)
	// if diags.HasErrors() {
	// 	for _, d := range diags {
	// 		t.Log(d)
	// 	}
	// 	t.Fatal(diags)
	// }
	// assert.True(t, val.GetAttr("enabled").True())
	// More assertions
}

func Transform(ctx context.Context, schema *jsonschema.Schema, attr hcl.Attribute) {
	// Create an evaluation context. This may include variables, functions, etc.
	// For the sake of this example, I'm using an empty context.
	evalContext := &hcl.EvalContext{}

	// Evaluate the attribute's expression to get a cty.Value
	val, diag := attr.Expr.Value(evalContext)
	if diag.HasErrors() {
		// Handle errors
		fmt.Println("Error evaluating HCL:", diag)
		return
	}

	// Convert the cty.Value to a JSON-friendly form
	jsonFriendlyVal, err := stdlib.JSONEncode(val)
	if err != nil {
		// Handle errors
		fmt.Println("Error converting to JSON-friendly value:", err)
		return
	}

	jsonData := map[string]interface{}{}
	// Marshal the JSON-friendly value to a JSON string
	err = json.Unmarshal([]byte(jsonFriendlyVal.AsString()), &jsonData)
	if err != nil {
		// Handle errors
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	pp.Println(jsonData)

	// Validate the JSON against the schema
	err = schema.Validate(jsonData)
	if err != nil {
		// Handle errors
		fmt.Println("Error validating JSON:", err)
		return
	}

	// The JSON is valid!
	fmt.Println("JSON is valid!")
}
