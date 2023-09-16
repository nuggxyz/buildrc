package schema

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/json"
	"github.com/k0kubun/pp/v3"
	"github.com/zclconf/go-cty/cty"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// load json or yaml schema file

func LoadJsonSchemaFile(ctx context.Context, document_uri string) (*jsonschema.Schema, error) {

	loader := jsonschema.NewCompiler()
	loader.Draft = jsonschema.Draft7

	resp, err := http.Get(document_uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status code %d", document_uri, resp.StatusCode)
	}

	// // add schema to loader
	// err = loader.AddResource(document_uri, dat)
	// if err != nil {
	// 	return nil, err
	// }

	// compile schema
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	schema, err := jsonschema.CompileString(document_uri, string(str))
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func JsonScemaToHCLSpec(ctx context.Context, schema *jsonschema.Schema) (hcldec.Spec, error) {

	// convert json schema to hcl schema
	// compile schema

	blk := hcldec.BlockObjectSpec{
		TypeName:   "export",
		LabelNames: []string{"root"},
		Nested:     nil,
	}

	spec, err := convertJSONSchemaToHCLDec(ctx, "test", false, schema)
	if err != nil {
		return nil, err
	}

	blk.Nested = spec

	return &blk, nil
}

// github actions schema https://raw.githubusercontent.com/SchemaStore/schemastore/master/src/schemas/json/github-workflow.json

// convert to hcl schema file
func LoadJSONSchemaFileAsHCL(ctx context.Context, document_uri string) (hcldec.Spec, error) {

	resp, err := http.Get(document_uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status code %d", document_uri, resp.StatusCode)
	} // convert json schema to hcl schema
	// compile schema
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	hcl, errd := json.Parse(str, "test.json")
	if errd != nil && errd.HasErrors() {
		return nil, errd
	}

	pp.Println(hcl)

	ctn, errd := loadSpecFile(hcl)
	if errd != nil && errd.HasErrors() {
		return nil, errd
	}

	return ctn.RootSpec, nil

}

func convertJSONSchemaToHCLDec(ctx context.Context, name string, required bool, prop *jsonschema.Schema) (hcldec.Spec, error) {
	if name == "jobs" {
		pp.Println(prop)
	}
	if prop.Ref != nil {
		return convertJSONSchemaToHCLDec(ctx, name, required, prop.Ref)
	}
	if prop.OneOf != nil {
		return convertTouple(ctx, name, required, prop.OneOf)
	}
	if prop.AnyOf != nil {
		return convertTouple(ctx, name, required, prop.AnyOf)
	}
	if len(prop.Types) == 0 {
		return convertObject(ctx, name, required, prop)
	}
	switch prop.Types[0] {
	case "object":
		return convertObject(ctx, name, required, prop)
	case "array":
		return convertArray(ctx, name, required, prop)
	case "string":
		return convertString(ctx, name, required, prop)
	case "number", "integer":
		return convertNumber(ctx, name, required, prop)
	case "boolean":
		return convertBoolean(ctx, name, required, prop)
	case "null":
		return &hcldec.AttrSpec{
			Name:     name,
			Type:     cty.NilType,
			Required: required,
		}, nil
	default:
		return nil, errors.New("Unsupported type definition: " + prop.Types[0])
	}
}

func convertTouple(ctx context.Context, name string, required bool, prop []*jsonschema.Schema) (hcldec.Spec, error) {
	objSpec := hcldec.ObjectSpec{}

	for _, v := range prop {
		// Convert each JSON schema to its corresponding HCL dec spec
		hclSpec, err := convertJSONSchemaToHCLDec(ctx, v.Title, false /* tuples usually don't have required elements */, v)
		if err != nil {
			return nil, err
		}

		// Append it to the tuple spec
		objSpec[v.Title] = hclSpec
	}

	return &hcldec.ValidateSpec{
		Wrapped: objSpec,
		Func: func(value cty.Value) hcl.Diagnostics {
			if value.IsNull() {
				return nil
			}

			// Validate the tuple
			return nil
		},
	}, nil
}

func convertObject(ctx context.Context, name string, required bool, prop *jsonschema.Schema) (hcldec.Spec, error) {
	objSpec := hcldec.ObjectSpec{}

	for propName, childProp := range prop.Properties {
		hclDecSpec, err := convertJSONSchemaToHCLDec(ctx, propName, slices.Contains(childProp.Required, propName), childProp)
		if err != nil {
			return nil, err
		}
		objSpec[propName] = hclDecSpec
	}

	return objSpec, nil
}

func convertArray(ctx context.Context, name string, required bool, prop *jsonschema.Schema) (hcldec.Spec, error) {
	if prop.Items == nil {
		return nil, errors.New("array items not defined")
	}

	var spec hcldec.Spec

	if itm, ok := prop.Items.(*jsonschema.Schema); ok {
		hclDecSpec, err := convertJSONSchemaToHCLDec(ctx, itm.Title, required, itm)
		if err != nil {
			return nil, err
		}
		spec = hclDecSpec
	} else if itms, ok := prop.Items.([]*jsonschema.Schema); ok {
		hcldecspec, err := convertTouple(ctx, name, required, itms)
		if err != nil {
			return nil, err
		}
		spec = hcldecspec
	} else {
		return nil, errors.New("array items not defined")
	}

	return &hcldec.BlockListSpec{
		TypeName: name,
		Nested:   spec,
		MinItems: prop.MinItems,
		MaxItems: prop.MaxItems,
	}, nil
}

func convertString(ctx context.Context, name string, required bool, prop *jsonschema.Schema) (hcldec.Spec, error) {

	return &hcldec.ValidateSpec{
		Wrapped: &hcldec.AttrSpec{
			Name:     name,
			Type:     cty.String,
			Required: required,
		},
		Func: func(value cty.Value) hcl.Diagnostics {
			if prop.Pattern != nil && !value.IsNull() {
				if !prop.Pattern.MatchString(value.AsString()) {
					return hcl.Diagnostics{
						{
							Severity: hcl.DiagError,
							Summary:  "Invalid value",
							Detail:   fmt.Sprintf("The value %q is not valid for %q.", value.AsString(), name),
						},
					}
				}
			}
			return nil
		},
	}, nil
}

func convertNumber(ctx context.Context, name string, required bool, prop *jsonschema.Schema) (hcldec.Spec, error) {
	return &hcldec.AttrSpec{
		Name:     name,
		Type:     cty.Number,
		Required: required,
	}, nil
}

func convertBoolean(ctx context.Context, name string, required bool, prop *jsonschema.Schema) (hcldec.Spec, error) {
	return &hcldec.AttrSpec{
		Name:     name,
		Type:     cty.Bool,
		Required: required,
	}, nil
}
