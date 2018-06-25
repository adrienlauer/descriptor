package model

import "errors"

type Provider struct {
	// The Name of the provider
	Name string
	// The environment referencing the provider
	root *Environment
	// The Repository/version of the provider
	Component
	// The provider attributes
	Parameters attributes
}

// Reference to a provider
type ProviderRef struct {
	// The referenced provider
	provider *Provider
	// The overwritten parameters of the provider
	Parameters attributes `yaml:",inline"`
}

// ProviderName returns the name of the referenced provider
func (p ProviderRef) ProviderName() string {
	return p.provider.Name
}

// ComponentId returns the id of the provider component
func (p ProviderRef) ComponentId() string {
	return p.provider.Component.Id
}

// Component returns the provider component
func (p ProviderRef) Component() Component {
	return p.provider.Component
}

// createProviders creates all the providers declared into the provided environment
func createProviders(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Provider {
	res := map[string]Provider{}
	if yamlEnv.Providers == nil || len(yamlEnv.Providers) == 0 {
		vErrs.AddError(errors.New("no provider specified"), "providers")
	} else {
		for name, yamlProvider := range yamlEnv.Providers {
			provider := Provider{
				root:       env,
				Parameters: createAttributes(yamlProvider.Params, nil),
				Name:       name,
			}

			provider.Component = createComponent(vErrs, env, "providers."+name, yamlProvider.Repository, yamlProvider.Version)

			res[name] = provider
		}
	}
	return res
}

// createProviderRef creates a reference to the provider declared into the yaml reference
func createProviderRef(vErrs *ValidationErrors, env *Environment, location string, yamlRef yamlRef) ProviderRef {
	if len(yamlRef.Name) == 0 {
		vErrs.AddError(errors.New("empty provider reference"), location)
	} else {
		if val, ok := env.Providers[yamlRef.Name]; ok {
			return ProviderRef{Parameters: createAttributes(yamlRef.Params, val.Parameters), provider: &val}
		} else {
			vErrs.AddError(errors.New("unknown provider reference: "+yamlRef.Name), location+".name")
		}
	}
	return ProviderRef{}
}
