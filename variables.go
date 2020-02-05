package main

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

const BadFormat = "Error - Environment variable '%s' should be of type '%s' (%s)"

func GetVariable(name string, defaultValue interface{} ) interface{} {
	env := os.Getenv(name)

	if env == "" {
		return defaultValue
	}

	returnType := reflect.TypeOf(defaultValue).String()
	var out, err interface{}

	switch returnType {
		case "int":
			out, err = strconv.Atoi(env)
		case "bool":
			out, err = strconv.ParseBool(env)
		case "float64":
			out, err = strconv.ParseFloat(env,64)
		case "string":
			return env
		default:
			log.Fatalf("Variable %s has a default value which type cannot be converted by GetVariable (type: %s)", name, returnType)
	}
	if err != nil {
		log.Fatalf(BadFormat, name, returnType, err)
	}
	return out
}