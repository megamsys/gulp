/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package carton

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

var (
	//the state actions available are.
	STATE   = "state"
	BOOT    = "boot"
	STATEUP = "stateup"
	DONE    = "done"

	//the control actions available are.
	CONTROL = "control"
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"

	//the operation actions is just one called upgrade
	OPERATIONS    = "operations"
	UPGRADE       = "upgrade"
	RESETPASSWORD = "resetpassword"
)

type ReqParser struct {
	name string
}

// NewParser returns a new instance of Parser.
func NewReqParser(n string) *ReqParser {
	return &ReqParser{name: n}
}

// ParseRequest parses a request string and returns its MegdProcess representation.
// eg: (state, create) => CreateProcess{}
// After figuring out the process, we operate on it.
func ParseRequest(n, s, a string) (MegdProcessor, error) {
	return NewReqParser(n).ParseRequest(s, a)
}

func (p *ReqParser) ParseRequest(category string, action string) (MegdProcessor, error) {
	switch category {
	case STATE:
		return p.parseState(action)
	case CONTROL:
		return p.parseControl(action)
	case OPERATIONS:
		return p.parseOperations(action)
	default:
		return nil, newParseError([]string{category, action}, []string{STATE, CONTROL, OPERATIONS})
	}
}

func (p *ReqParser) parseState(action string) (MegdProcessor, error) {
	switch action {
	case BOOT:
		return BootProcess{
			Name: p.name,
		}, nil
	case STATEUP:
		return StateupProcess{
			Name: p.name,
		}, nil
	default:
		return nil, newParseError([]string{STATE, action}, []string{BOOT, STATEUP})
	}
}

func (p *ReqParser) parseControl(action string) (MegdProcessor, error) {
	switch action {
	case START:
		return StartProcess{
			Name: p.name,
		}, nil
	case STOP:
		return StopProcess{
			Name: p.name,
		}, nil
	case RESTART:
		return RestartProcess{
			Name: p.name,
		}, nil
	default:
		return nil, newParseError([]string{CONTROL, action}, []string{START, STOP, RESTART})
	}
}

func (p *ReqParser) parseOperations(action string) (MegdProcessor, error) {
	switch action {
	case UPGRADE:
		return UpgradeProcess{
			Name: p.name,
		}, nil
	case RESETPASSWORD:
		return ResetPasswordProcess{
			Name: p.name,
		}, nil
	default:
		return nil, newParseError([]string{OPERATIONS, action}, []string{UPGRADE})
	}
}

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Found    string
	Expected []string
}

// newParseError returns a new instance of ParseError.
func newParseError(found []string, expected []string) *ParseError {
	return &ParseError{Found: strings.Join(found, ","), Expected: expected}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("found %s, expected %s", e.Found, strings.Join(e.Expected, ", "))
}

type Requests struct {
	Id        string `json:"id"`     //assembly id
	CatId     string `json:"cat_id"` // assemblies_id
	Action    string `json:"action"` // start, stop ...
	AccountId string `json:"email"`
	Category  string `json:"category"` // state, control, policy
	CreatedAt string `json:"created_at"`
}

type ApiRequests struct {
	JsonClaz string     `json:"json_claz" cql:"json_claz"`
	Results  []Requests `json:"results" cql:"results"`
}

func (r *Requests) String() string {
	if d, err := yaml.Marshal(r); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}
