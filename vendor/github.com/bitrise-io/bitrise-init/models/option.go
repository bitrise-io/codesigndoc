package models

import (
	"encoding/json"
	"fmt"
)

// Type is to select the user interaction type that is required to fill an option
type Type string

// OptionTypes we currently support, list, user-input, optional-user-input
const (
	// Originally:
	// - if there was only one key then this question was not asked
	// - if there was more than one keys then it was asked in a list to select from
	// Now, if this type is selected:
	// - if there is only one key then this question shall not not be asked
	// - if there are more than one keys then the selection must be asked in a list to select from
	TypeSelector Type = "selector"
	// Originally:
	// - if there was only one key then this question was not asked
	// - if there was more than one keys then it was asked in a list to select from
	// Now, if this type is selected:
	// - if there is only one key then the view should let the user select of the only item the list has or manual input
	// - if there are more than one keys then the selection must be asked in a list to select from or have a manual button to let the user to input anything else
	TypeOptionalSelector Type = "selector_optional"
	// Originally:
	// - if there was only one key and it's name was `_` then we shown an input field to the user to type his value in and it was a requirement to have an input value
	// Now, if this type is selected:
	// - we must show an input field to the user and it is required to fill, any name for the key will be the placeholder value for the input field
	TypeUserInput Type = "user_input"
	// Originally:
	// - if there was only one key and it's name was `_` then we shown an input field to the user to type his value in and it was a requirement to have an input value
	// Now, if this type is selected:
	// - we must show an input field to the user and it is NOT required to be filled, can be empty, and any name for the key will be the placeholder value for the input field
	TypeOptionalUserInput Type = "user_input_optional"
)

// OptionNode ...
type OptionNode struct {
	Title          string                 `json:"title,omitempty" yaml:"title,omitempty"`
	Summary        string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	EnvKey         string                 `json:"env_key,omitempty" yaml:"env_key,omitempty"`
	Type           Type                   `json:"type,omitempty" yaml:"type,omitempty"`
	ChildOptionMap map[string]*OptionNode `json:"value_map,omitempty" yaml:"value_map,omitempty"`
	// Leafs only
	Config string   `json:"config,omitempty" yaml:"config,omitempty"`
	Icons  []string `json:"icons,omitempty" yaml:"icons,omitempty"`

	Components []string    `json:"-" yaml:"-"`
	Head       *OptionNode `json:"-" yaml:"-"`
}

// NewOption ...
func NewOption(title, summary, envKey string, optionType Type) *OptionNode {
	return &OptionNode{
		Title:          title,
		Summary:        summary,
		EnvKey:         envKey,
		ChildOptionMap: map[string]*OptionNode{},
		Components:     []string{},
		Type:           optionType,
	}
}

// NewConfigOption ...
func NewConfigOption(name string, icons []string) *OptionNode {
	return &OptionNode{
		ChildOptionMap: map[string]*OptionNode{},
		Config:         name,
		Icons:          icons,
		Components:     []string{},
	}
}

func (option *OptionNode) String() string {
	bytes, err := json.MarshalIndent(option, "", "\t")
	if err != nil {
		return fmt.Sprintf("failed to marshal, error: %s", err)
	}
	return string(bytes)
}

// IsConfigOption ...
func (option *OptionNode) IsConfigOption() bool {
	return option.Config != ""
}

// IsValueOption ...
func (option *OptionNode) IsValueOption() bool {
	return option.Title != ""
}

// IsEmpty ...
func (option *OptionNode) IsEmpty() bool {
	return !option.IsValueOption() && !option.IsConfigOption()
}

// AddOption ...
func (option *OptionNode) AddOption(forValue string, newOption *OptionNode) {
	option.ChildOptionMap[forValue] = newOption

	if newOption != nil {
		newOption.Components = append(option.Components, forValue)

		if option.Head == nil {
			// first option's head is nil
			newOption.Head = option
		} else {
			newOption.Head = option.Head
		}
	}
}

// AddConfig ...
func (option *OptionNode) AddConfig(forValue string, newConfigOption *OptionNode) {
	option.ChildOptionMap[forValue] = newConfigOption

	if newConfigOption != nil {
		newConfigOption.Components = append(option.Components, forValue)

		if option.Head == nil {
			// first option's head is nil
			newConfigOption.Head = option
		} else {
			newConfigOption.Head = option.Head
		}
	}
}

// Parent ...
func (option *OptionNode) Parent() (*OptionNode, string, bool) {
	if option.Head == nil {
		return nil, "", false
	}

	parentComponents := option.Components[:len(option.Components)-1]
	parentOption, ok := option.Head.Child(parentComponents...)
	if !ok {
		return nil, "", false
	}
	underKey := option.Components[len(option.Components)-1:][0]
	return parentOption, underKey, true
}

// Child ...
func (option *OptionNode) Child(components ...string) (*OptionNode, bool) {
	currentOption := option
	for _, component := range components {
		childOption := currentOption.ChildOptionMap[component]
		if childOption == nil {
			return nil, false
		}
		currentOption = childOption
	}
	return currentOption, true
}

// LastChilds ...
func (option *OptionNode) LastChilds() []*OptionNode {
	lastOptions := []*OptionNode{}

	var walk func(*OptionNode)
	walk = func(opt *OptionNode) {
		if len(opt.ChildOptionMap) == 0 {
			lastOptions = append(lastOptions, opt)
			return
		}

		for _, childOption := range opt.ChildOptionMap {
			if childOption == nil {
				lastOptions = append(lastOptions, opt)
				return
			}

			if childOption.IsConfigOption() {
				lastOptions = append(lastOptions, opt)
				return
			}

			if childOption.IsEmpty() {
				lastOptions = append(lastOptions, opt)
				return
			}

			walk(childOption)
		}
	}

	walk(option)

	return lastOptions
}

// RemoveConfigs ...
func (option *OptionNode) RemoveConfigs() {
	lastChilds := option.LastChilds()
	for _, child := range lastChilds {
		for _, child := range child.ChildOptionMap {
			child.Config = ""
		}
	}
}

// AttachToLastChilds ...
func (option *OptionNode) AttachToLastChilds(opt *OptionNode) {
	childs := option.LastChilds()
	for _, child := range childs {
		values := child.GetValues()
		for _, value := range values {
			child.AddOption(value, opt)
		}
	}
}

// Copy ...
func (option *OptionNode) Copy() *OptionNode {
	bytes, err := json.Marshal(*option)
	if err != nil {
		return nil
	}

	var optionCopy OptionNode
	if err := json.Unmarshal(bytes, &optionCopy); err != nil {
		return nil
	}

	return &optionCopy
}

// GetValues ...
func (option *OptionNode) GetValues() []string {
	if option.Config != "" {
		return []string{option.Config}
	}

	values := []string{}
	for value := range option.ChildOptionMap {
		values = append(values, value)
	}
	return values
}
