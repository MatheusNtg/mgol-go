package parser

import (
	"encoding/json"
	"io/ioutil"
)

type Rule struct {
	Number int      `json:"rule_number"`
	Left   string   `json:"left"`
	Right  []string `json:"right"`
}

type RulesMap map[int]Rule

var rulesMapInstance *RulesMap

func loadGrammarRules(path string) []Rule {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	rules := []Rule{}
	err = json.Unmarshal(file, &rules)
	if err != nil {
		panic(err)
	}
	return rules
}

func createMapFromSlice(rules []Rule) *RulesMap {
	rulesMap := make(RulesMap)
	for _, rule := range rules {
		rulesMap[rule.Number] = rule
	}
	return &rulesMap
}

func GetRulesMap(path string) *RulesMap {
	if rulesMapInstance == nil {
		rules := loadGrammarRules(path)
		rulesMapInstance = createMapFromSlice(rules)
		return rulesMapInstance
	} else {
		return rulesMapInstance
	}
}
