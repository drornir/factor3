// GENERATED FILE DO NOT EDIT
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

var _ = strconv.ParseInt // just in case strconv is not used

func (self *Config) Factor3Load(argv []string) error {

	fmt.Printf("type name: Config\n")
	fmt.Printf("annotations: //factor3:generate --filename ./example/app/config.yaml --env-prefix EX,\n")
	fmt.Printf("fields:\n")

	fmt.Printf("\t-name=DBConnection, type=string\n")
	fmt.Printf("\t-annotations: \n")

	fmt.Printf("\t-name=Port, type=string\n")
	fmt.Printf("\t-annotations: //factor3:pflag port,\n")

	fmt.Printf("\t-name=SomeNumber, type=int\n")
	fmt.Printf("\t-annotations: //factor3:pflag some-number n,\n")

	fmt.Printf("\t-name=SomeFlag, type=bool\n")
	fmt.Printf("\t-annotations: //factor3:pflag some-flag,\n")

	conf := zz_factor3_Config{}
	conf.Filename = "./example/app/config.yaml"
	conf.EnvPrefix = "EX_"

	type jsonStruct struct {
		DBConnection string `json:"db_connection"`
		Port         string `json:"port"`
		SomeNumber   int    `json:"some_number"`
		SomeFlag     bool   `json:"some_flag"`
	}

	var jsoner json.Unmarshaler
	if x, ok := interface{}(self).(json.Unmarshaler); ok {
		jsoner = x
	} else {
		jsoner = &zz_factor3_JSONer[jsonStruct]{t: (*jsonStruct)(self)}
	}

	loadConfigFile := func(filename string) error {
		file, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		fmt.Printf("%s\n", file)

		fileExt := filename[strings.LastIndex(filename, ".")+1:]
		switch fileExt {
		case "yaml", "yml":
			intoMap := make(map[string]interface{})
			err = yaml.Unmarshal(file, intoMap)
			if err != nil {
				break
			}
			intoJSON, e := json.Marshal(intoMap)
			if e != nil {
				err = e
				break
			}
			err = json.Unmarshal(intoJSON, jsoner)
		case "json":
			err = json.Unmarshal(file, jsoner)
		default:
			return fmt.Errorf("unsupported file type %q", fileExt)
		}
		if err != nil {
			return fmt.Errorf("unmarshaling: %w", err)
		}

		return nil
	}

	loadEnv := func(prefix string) error {
		var s string
		s = os.Getenv(conf.EnvPrefix + "DB_CONNECTION")
		if s != "" {
			self.DBConnection = s
		}
		s = os.Getenv(conf.EnvPrefix + "PORT")
		if s != "" {
			self.Port = s
		}
		s = os.Getenv(conf.EnvPrefix + "SOME_NUMBER")
		if s != "" {
			if n, err := strconv.ParseInt(s, 10, 32); err != nil {
				return fmt.Errorf("parsing \"SomeNumber\" as \"int\": %w", err)
			} else {
				self.SomeNumber = int(n)
			}
		}
		s = os.Getenv(conf.EnvPrefix + "SOME_FLAG")
		if s != "" {
			if b, err := strconv.ParseBool(s); err != nil {
				return fmt.Errorf("parsing \"SomeFlag\" as \"bool\": %w", err)
			} else {
				self.SomeFlag = b
			}
		}

		fmt.Printf("here is was I got:\n\t%#v\n", *self)

		return nil
	}

	parseFlags := func(argv []string) error {
		if len(argv) == 0 {
			return nil
		}
		fset := pflag.NewFlagSet("Config", pflag.ContinueOnError)
		fset.StringVarP(&self.Port, "port", "", self.Port, "")
		fset.IntVarP(&self.SomeNumber, "some-number", "n", self.SomeNumber, "")
		fset.BoolVarP(&self.SomeFlag, "some-flag", "", self.SomeFlag, "")

		if err := fset.Parse(argv); err != nil {
			return fmt.Errorf("parsing flags: %w", err)
		}
		return nil
	}

	if err := loadConfigFile(conf.Filename); err != nil {
		return fmt.Errorf("loading config from file %q: %w", conf.Filename, err)
	}
	if err := loadEnv(conf.EnvPrefix); err != nil {
		return fmt.Errorf("loading config from env: %w", err)
	}
	if err := parseFlags(argv); err != nil {
		return fmt.Errorf("loading config from pflags: %w", err)
	}

	return nil
}
