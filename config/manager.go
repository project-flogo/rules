package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/project-flogo/core/app/resource"
)

const (
	uriSchemeRes        = "res://"
	RESTYPE_RULESESSION = "rulesession"
)

type ResourceManager struct {
	configs map[string]*RuleActionDescriptor
}

func NewResourceManager() *ResourceManager {
	manager := &ResourceManager{}
	manager.configs = make(map[string]*RuleActionDescriptor)

	return manager
}

func (m *ResourceManager) LoadResource(resConfig *resource.Config) (*resource.Resource, error) {
	var rsConfig *RuleActionDescriptor
	err := json.Unmarshal(resConfig.Data, &rsConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling rulesession resource with id '%s', %s", resConfig.ID, err.Error())
	}

	m.configs[resConfig.ID] = rsConfig
	return resource.New("rulesession", m.configs), nil
}

func (m *ResourceManager) GetResource(id string) interface{} {
	return m.configs[id]
}

func (m *ResourceManager) GetRuleSessionDescriptor(uri string) (*RuleSessionDescriptor, error) {

	if strings.HasPrefix(uri, uriSchemeRes) {
		return &RuleSessionDescriptor{m.configs[uri[len(uriSchemeRes):]].Rules}, nil
	}

	return nil, errors.New("cannot find RuleSession: " + uri)
}

func (m *ResourceManager) GetRuleActionDescriptor(uri string) (*RuleActionDescriptor, error) {

	if strings.HasPrefix(uri, uriSchemeRes) {
		return m.configs[uri[len(uriSchemeRes):]], nil
	}

	return nil, errors.New("cannot find RuleSession: " + uri)
}

//ioMetadata support
/*
type ActionResource struct {
	IOMetadata *metadata.IOMetadata `json:"metadata"`
}

type ResManager struct {
	IOMetadata *metadata.IOMetadata
}

func (m *ResManager) LoadResource(resConfig *resource.Config) (*resource.Resource, error) {

	var res *ActionResource
	err := json.Unmarshal(resConfig.Data, &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling metadata resource with id '%s', %s", resConfig.ID, err.Error())
	}

	m.IOMetadata = res.IOMetadata
	return resource.New("ruleaction", m.IOMetadata), nil
}
*/
