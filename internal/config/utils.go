package config

import (
	"reflect"

	"github.com/spf13/pflag"
)

func isFlagPassed(name string) bool {
	found := false
	pflag.Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func ParseFlags(cfg *Config) {
	t := reflect.TypeOf(*cfg)

	// Custom Types
	var dsn string
	pflag.StringVarP(
		&dsn,
		"dsn",
		"d",
		"",
		"database dsn",
	)

	// General Types
	var str string
	strValues := map[string]*string{}

	var number int64
	numberValues := map[string]*int64{}

	var boolean bool
	booleanValues := map[string]*bool{}

	fieldNames := map[string]string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		flagName := field.Tag.Get("flag")
		flagShort := field.Tag.Get("flagShort")
		flagDescription := field.Tag.Get("flagDescription")

		fieldNames[field.Name] = flagName

		switch {
		case field.Type == reflect.TypeOf(str):
			strVal := ""
			strValues[field.Name] = &strVal
			pflag.StringVarP(
				strValues[field.Name],
				flagName,
				flagShort,
				strVal,
				flagDescription,
			)
		case field.Type == reflect.TypeOf(number):
			numberVal := int64(0)
			numberValues[field.Name] = &numberVal
			pflag.Int64VarP(
				numberValues[field.Name],
				flagName,
				flagShort,
				numberVal,
				flagDescription,
			)
		case field.Type == reflect.TypeOf(boolean):
			booleanVal := false
			booleanValues[field.Name] = &booleanVal
			pflag.BoolVarP(
				booleanValues[field.Name],
				flagName,
				flagShort,
				booleanVal,
				flagDescription,
			)
		}
	}

	pflag.Parse()

	// Custom Types
	if dsn != "" {
		cfg.Postgres.DSN = dsn
	}

	// General Types
	for k, v := range strValues {
		flagKey := fieldNames[k]
		if isFlagPassed(flagKey) {
			reflect.
				ValueOf(cfg).
				Elem().
				FieldByName(k).
				SetString(*v)
		}
	}

	for k, v := range numberValues {
		flagKey := fieldNames[k]
		if isFlagPassed(flagKey) {
			reflect.
				ValueOf(cfg).
				Elem().
				FieldByName(k).
				SetInt(*v)
		}
	}

	for k, v := range booleanValues {
		flagKey := fieldNames[k]
		if isFlagPassed(flagKey) {
			reflect.
				ValueOf(cfg).
				Elem().
				FieldByName(k).
				SetBool(*v)
		}
	}
}
