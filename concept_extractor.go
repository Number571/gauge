// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package main

import "strings"

const (
	CONCEPT_HEADING_TEMPLATE = "Concept Heading"
	SPEC_HEADING_TEMPLATE    = "# S\n"
)

type extractConceptResult struct {
	heading     string
	stepTexts   string
	conceptText string
	hasParam    bool
	isValid     bool
}

func getTextForConcept(stepsToExtract string) *extractConceptResult {
	specs, parseResult := new(specParser).parse(SPEC_HEADING_TEMPLATE+stepsToExtract, &conceptDictionary{})
	if !parseResult.ok {
		return &extractConceptResult{heading: "", stepTexts: "", conceptText: "", hasParam: false, isValid: false}
	}
	steps := make([]*step, 0)
	for _, item := range specs.items {
		if item.kind() == stepKind {
			steps = append(steps, item.(*step))
		} else if item.kind() != commentKind {
			return &extractConceptResult{heading: "", stepTexts: "", conceptText: "", hasParam: false, isValid: false}
		}
	}
	stepTexts, args, table := getSteps(steps)
	heading, text, hasParam := getHeadingAndText(args, table)
	return &extractConceptResult{heading: heading, stepTexts: stepTexts, conceptText: text, hasParam: hasParam, isValid: true}
}
func getSteps(steps []*step) (string, map[string]bool, table) {
	stepTexts := ""
	argsMap := getArgsMap(steps)
	args := make(map[string]bool)
	table := table{}
	for _, step := range steps {
		for _, arg := range step.args {
			value := arg.String()
			if arg.argType == tableArg {
				value = arg.table.String()
			}
			if argsMap[value] > 1 {
				value := arg.value
				if arg.argType == tableArg && (!table.isInitialized() || table.String() == arg.table.String()) {
					table = arg.table
					value = "table"
				} else if arg.argType == tableArg && table.String() != arg.table.String() {
					continue
				} else if arg.argType != tableArg {
					args[value] = true
				}
				arg.argType = dynamic
				arg.name = "<" + value + ">"
				arg.value = value
			}
		}
		stepTexts += formatItem(step)
	}
	return stepTexts, args, table
}

func getHeadingAndText(args map[string]bool, table table) (string, string, bool) {
	conceptHeading :=
		CONCEPT_HEADING_TEMPLATE
	conceptText := "* " + CONCEPT_HEADING_TEMPLATE
	hasParam := false
	for name, _ := range args {
		hasParam = true
		conceptHeading += " <" + name + ">"
		conceptText += " \"" + name + "\""
	}
	if table.isInitialized() {
		hasParam = true
		conceptHeading += " <table>"
		conceptText += "\n" + formatTable(&table)
	}
	return conceptHeading, conceptText, hasParam
}

func getArgsMap(steps []*step) map[string]int {
	argsMap := make(map[string]int)
	for _, step := range steps {
		for _, arg := range step.args {
			value := arg.String()
			if arg.argType == tableArg {
				value = arg.table.String()
			}
			if _, ok := argsMap[value]; !ok {
				argsMap[value] = 1
			} else {
				argsMap[value] += 1
			}
		}
	}
	return argsMap
}

func refactorConceptHeading(newConceptHeading string, oldConceptHeading string, oldConceptText string) string {
	removeIdentifier := func(text string, identifier string) string {
		if strings.HasPrefix(text, identifier) {
			text = strings.TrimPrefix(text, identifier)
		}
		return text
	}
	newConceptHeading = removeIdentifier(newConceptHeading, "#")
	oldConceptHeading = removeIdentifier(oldConceptHeading, "#")
	agent, _ := getRefactorAgent(oldConceptHeading, newConceptHeading)
	argsOrder := agent.createOrderOfArgs()
	spec, _ := new(specParser).parse(SPEC_HEADING_TEMPLATE+oldConceptText, &conceptDictionary{})
	oldConcept := spec.items[0].(*step)
	tokens, _ := new(specParser).generateTokens("*" + newConceptHeading)
	step, _ := (&specification{}).createStepUsingLookup(tokens[0], nil)
	oldConcept.value = step.value
	oldConcept.args = oldConcept.getArgsInOrder(*oldConcept, argsOrder)
	value := formatStep(oldConcept)
	return value
}
