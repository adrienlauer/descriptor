package model

import (
	"errors"
)

//Platform the platform used to build an environment
type Platform struct {
	Base               Base
	Distribution       Distribution
	Components         map[string]Component
	componentResolvers []componentResolver
}

func createPlatform(yamlEnv *yamlEnvironment) (*Platform, error) {

	p := &Platform{}
	// Compute the component base for the environment
	base, e := CreateComponentBase(yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the base component : " + e.Error())
	}
	p.Base = base

	// Create the distribution component (mandatory)
	dist, e := CreateDistribution(base, yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the distribution : " + e.Error())
	}
	p.Distribution = dist

	// Create other components of the environment
	components := map[string]Component{}
	for name, yamlC := range yamlEnv.Ekara.Components {
		repo, e := CreateRepository(base, yamlC.Repository, yamlC.Ref, "")
		if e != nil {
			return p, errors.New("Error creating the repository: " + e.Error())
		}
		repo.setAuthentication(yamlC)
		component := CreateComponent(name, repo)

		components[name] = component
	}

	p.Components = components
	p.componentResolvers = make([]componentResolver, 0, 0)
	return p, nil
}

func (p *Platform) tagUsedComponent(cr componentResolver) {
	p.componentResolvers = append(p.componentResolvers, cr)
}

// UsedComponents returns an array of components effectively in used throughout the descriptor.
func (p *Platform) UsedComponents() ([]Component, error) {
	res := make([]Component, 0, 0)
	temp := make(map[string]Component)
	for _, cr := range p.componentResolvers {
		c, err := cr.ResolveComponent()
		if err != nil {
			return res, err
		}
		temp[c.Id] = c
	}
	for _, c := range temp {
		res = append(res, c)
	}
	return res, nil
}

func (p Platform) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	for _, c := range p.Components {
		vErrs.merge(ErrorOnInvalid(c))
	}
	return vErrs
}

func (p *Platform) merge(other Platform) error {
	for id, c := range other.Components {
		if id != "" {
			if _, ok := p.Components[id]; !ok {
				p.Components[id] = c
			}
		}
	}

	for _, c := range other.componentResolvers {
		p.tagUsedComponent(c)
	}

	if p.Distribution.Repository.Url == nil && other.Distribution.Repository.Url != nil {
		p.Distribution = other.Distribution
	}
	return nil
}
