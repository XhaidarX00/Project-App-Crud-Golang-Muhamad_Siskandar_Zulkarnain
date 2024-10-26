package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

type Utils struct {
	key   interface{}
	value interface{}
}

func ErrorMessage(text string) {
	red := color.New(color.FgRed).SprintFunc()
	err := errors.New(text)
	fmt.Printf("⚠️   %s%s\n", red("Error : "), err)
}

func SuccesMessage(text string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("✔️   %s%s\n", green("Succes : "), text)
}

func ColorMessage(color_ string, text string) string {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	switch color_ {
	case "red":
		return red(text)
	case "green":
		return green(text)
	case "blue":
		return blue(text)
	case "yellow":
		return yellow(text)
	}

	return "0"
}

func PrintColorMsg(color string, text string) {
	fmt.Println(ColorMessage(color, text))
}

func Capitalize(s string) string {
	s = strings.ToLower(s)

	capitalizeNext := true
	result := []rune(s)

	for i, r := range result {
		if capitalizeNext && unicode.IsLetter(r) {
			result[i] = unicode.ToUpper(r)
			capitalizeNext = false
		} else if unicode.IsSpace(r) {
			capitalizeNext = true
		}
	}

	return string(result)
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		module := exec.Command("cmd", "/c", "cls")
		module.Stdout = os.Stdout
		module.Run()
	} else {
		module := exec.Command("clear")
		module.Stdout = os.Stdout
		module.Run()
	}
}

func IsLenVar(input any, len_ int) bool {
	lenVal := reflect.ValueOf(input)

	switch lenVal.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if lenVal.Len() >= len_ {
			return true
		}
	default:
		ErrorMessage(fmt.Sprintf("tipe data %s tidak mendukung operasi Len()\n", lenVal.Kind()))
	}

	return false
}

func ConvertSliceToMap(slice interface{}) []map[string]interface{} {
	var data []map[string]interface{}

	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return nil
	}

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()
		itemMap := make(map[string]interface{})
		itemVal := reflect.ValueOf(item)

		for j := 0; j < itemVal.NumField(); j++ {
			field := itemVal.Type().Field(j)
			itemMap[field.Name] = itemVal.Field(j).Interface()
		}
		data = append(data, itemMap)
	}

	return data
}

func GetStructKeys(p interface{}, dk []string) []string {
	var keys []string
	val := reflect.ValueOf(p)

	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		excluded := false
		for _, v := range dk {
			if fieldName == v {
				excluded = true
				break
			}
		}
		if !excluded {
			keys = append(keys, fieldName)
		}
	}

	return keys
}

// func ConvertMapToStruct(data map[string]interface{}, outputType reflect.Type) (interface{}, error) {
// 	if outputType.Kind() != reflect.Struct {
// 		return nil, fmt.Errorf("outputType must be a struct type")
// 	}

// 	output := reflect.New(outputType).Elem()

// 	for key, value := range data {
// 		fieldVal := output.FieldByName(key)
// 		if !fieldVal.IsValid() {
// 			continue
// 		}
// 		if !fieldVal.CanSet() {
// 			continue
// 		}

// 		valType := reflect.TypeOf(value)
// 		if fieldVal.Type().AssignableTo(valType) {
// 			fieldVal.Set(reflect.ValueOf(value))
// 		} else {
// 			return nil, fmt.Errorf("type mismatch for field %s: expected %s but got %s", key, fieldVal.Type(), valType)
// 		}
// 	}

// 	return output.Interface(), nil
// }
