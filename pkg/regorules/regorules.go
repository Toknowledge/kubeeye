package regorules

import (
	"context"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed rules
var defaultRegoRules embed.FS

// GetRegoRules get rego rules , put it into the channel RegoRulesListChan.
// func GetRegoRules(additionalRegoRulePath string) {
// 	var regoRulesList kube.RegoRulesList
// 	if additionalRegoRulePath != "" {
// 		GetRegoRulesfiles(additionalRegoRulePath, &regoRulesList)
// 	}
// 	GetDefaultRegofile("rules", &regoRulesList)

// 	kube.RegoRulesListChan <- regoRulesList
// }

// GetRegoRulesfiles get rego rules , put it into pointer of RegoRulesList
func GetRegoRulesfiles(path string) []string{
	var regoRules []string
	if path == "" {
		return regoRules
	}
	pathabs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(pathabs, "/") == false {
		pathabs += "/"
	}
	files, err := ioutil.ReadDir(pathabs)
	if err != nil {
		fmt.Println("Failed to read the dir of rego rule files.")
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".rego") == false {
			continue
		}

		getregoRule, _ := ioutil.ReadFile(pathabs + file.Name())
		regoRule := string(getregoRule)
		regoRules = append(regoRules, regoRule)
	}
	return regoRules
}

func GetDefaultRegofile(path string) []string{
	var regoRules []string
	files, err := defaultRegoRules.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		rule, _ := defaultRegoRules.ReadFile(path + "/" + file.Name())
		regoRule := string(rule)
		regoRules = append(regoRules, regoRule)
	}
	return regoRules
}

func MergeRegoRules(ctx context.Context, channels ...[]string) <-chan string {
	res := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	mergeRegoRuls := func(ctx context.Context, ch []string) {
		defer wg.Done()
		for _, c := range ch {
			res <- c
		}
	}

	for _, c := range channels {
		go mergeRegoRuls(ctx, c)
	}

	go func() {
		wg.Wait()
		defer close(res)
	}()
	return res
}