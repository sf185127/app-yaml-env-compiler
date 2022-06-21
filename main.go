package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	fmt.Println("Ready to compile ...")
	fileNameString := os.Getenv("INPUT_APPYAMLPATH")
	fmt.Println("Using " + fileNameString)
	filename, _ := filepath.Abs(fileNameString)
	
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var mapResult map[interface{}]interface{}
	err = yaml.Unmarshal(yamlFile, &mapResult)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(fmt.Sprintf("Env variables will be replaced: %v",mapResult["env_variables"]))

	for k, any := range mapResult {
		if k == "env_variables" || k == "build_env_variables" {
			err := checkIsPointer(&any)
			if err != nil {
				panic(err)
			}
			valueOf := reflect.ValueOf(any)
			val := reflect.Indirect(valueOf)
			switch val.Type().Kind() {
			case reflect.Map:
				envMap := any.(map[interface{}]interface{})
				for in, iv := range envMap {

					envName := in.(string)
					envVal := iv.(string)

					env := strings.Replace(strings.TrimSpace(envVal), "$", "", -1)
					osVal := os.Getenv(env)
					
					if len(osVal) > 0 {
						envMap[envName] = osVal
					}
				}
			default:
				panic(fmt.Sprintf("This is not supposed to happen, but if it does, good luck"))
			}
		}
	}
	
	fmt.Println(fmt.Sprintf("Compiled env variables: %v", mapResult["env_variables"]))
	fmt.Println(fmt.Sprintf("Compiled build env variables: %v", mapResult["build_env_variables"]))

	out, err := yaml.Marshal(mapResult)
	// write the whole body at once
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		panic(err)
	}
}

func checkIsPointer(any interface{}) error {
	if reflect.ValueOf(any).Kind() != reflect.Ptr {
		return fmt.Errorf("You passed something that was not a pointer: %s", any)
	}
	return nil
}
