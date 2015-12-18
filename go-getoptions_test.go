// This file is part of go-getoptions.
//
// Copyright (C) 2015  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package getoptions

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestIsOption(t *testing.T) {
	Debug.SetOutput(os.Stderr)
	Debug.SetOutput(ioutil.Discard)

	cases := []struct {
		in       string
		mode     string
		options  []string
		argument string
	}{
		{"opt", "bundling", []string{}, ""},
		{"--opt", "bundling", []string{"opt"}, ""},
		{"--opt=arg", "bundling", []string{"opt"}, "arg"},
		{"-opt", "bundling", []string{"o", "p", "t"}, ""},
		{"-opt=arg", "bundling", []string{"o", "p", "t"}, "arg"},
		{"-", "bundling", []string{"-"}, ""},
		{"--", "bundling", []string{"--"}, ""},

		{"opt", "singleDash", []string{}, ""},
		{"--opt", "singleDash", []string{"opt"}, ""},
		{"--opt=arg", "singleDash", []string{"opt"}, "arg"},
		{"-opt", "singleDash", []string{"o"}, "pt"},
		{"-", "singleDash", []string{"-"}, ""},
		{"--", "singleDash", []string{"--"}, ""},

		{"opt", "normal", []string{}, ""},
		{"--opt", "normal", []string{"opt"}, ""},
		{"--opt=arg", "normal", []string{"opt"}, "arg"},
		{"-opt", "normal", []string{"opt"}, ""},
		{"-", "normal", []string{"-"}, ""},
		{"--", "normal", []string{"--"}, ""},
	}
	for _, c := range cases {
		options, argument := isOption(c.in, c.mode)
		if !reflect.DeepEqual(options, c.options) || argument != c.argument {
			t.Errorf("isOption(%q, %q) == (%q, %q), want (%q, %q)",
				c.in, c.mode, options, argument, c.options, c.argument)
		}
	}
}

/*
// Verifies that a panic is reached when the same option is defined twice.
func TestDuplicateDefinition(t *testing.T) {
	opt := GetOptions()
	opt.Flag("flag", []string{"f"})
	opt.Flag("flag", []string{"f"})
}
*/

func TestWarningOrErrorOnUnknown(t *testing.T) {
	opt := GetOptions()
	_, err := opt.Parse([]string{"--flags"})
	if err == nil {
		t.Errorf("Unknown option 'flags' didn't raise error")
	}
	if err != nil && err.Error() != "Unknown option 'flags'" {
		t.Errorf("Error string didn't match expected value")
	}
}

func TestMissingArgument(t *testing.T) {
	// Missing argument without default
	opt := GetOptions()
	opt.String("string")
	_, err := opt.Parse([]string{"--string"})
	if err == nil {
		t.Errorf("Missing argument for option 'string' didn't raise error")
	}
	if err != nil && err.Error() != "Missing argument for option 'string'!" {
		t.Errorf("Error string didn't match expected value")
	}

	// Missing argument with default
	opt = GetOptions()
	opt.StringOptional("string", "default")
	_, err = opt.Parse([]string{"--string"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if opt.Option["string"] != "default" {
		t.Errorf("Default value not set for 'string'")
	}
}

func TestGetOptFlag(t *testing.T) {
	setup := func() *GetOpt {
		opt := GetOptions()
		opt.Flag("flag")
		opt.NFlag("nflag")
		return opt
	}

	cases := []struct {
		opt    *GetOpt
		option string
		input  []string
		value  bool
	}{
		{setup(),
			"flag",
			[]string{"--flag"},
			true,
		},
		{setup(),
			"nflag",
			[]string{"--nflag"},
			true,
		},
		{setup(),
			"nflag",
			[]string{"--no-nflag"},
			false,
		},
	}
	for _, c := range cases {
		_, err := c.opt.Parse(c.input)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if c.opt.Option[c.option] != c.value {
			t.Errorf("Wrong value: %v != %v", c.opt.Option[c.option], c.value)
		}
	}
}

func TestGetOptAliases(t *testing.T) {
	setup := func() *GetOpt {
		opt := GetOptions()
		opt.Flag("flag", "f", "h")
		return opt
	}

	cases := []struct {
		opt    *GetOpt
		option string
		input  []string
		value  bool
	}{
		{setup(),
			"flag",
			[]string{"--flag"},
			true,
		},
		{setup(),
			"flag",
			[]string{"-f"},
			true,
		},
		{setup(),
			"flag",
			[]string{"-h"},
			true,
		},
		// TODO: Add flag to allow for this.
		{setup(),
			"flag",
			[]string{"--fl"},
			true,
		},
	}
	for _, c := range cases {
		_, err := c.opt.Parse(c.input)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if c.opt.Option[c.option] != c.value {
			t.Errorf("Wrong value: %v != %v", c.opt.Option[c.option], c.value)
		}
	}

	opt := GetOptions()
	opt.Flag("flag")
	opt.Flag("fleg")
	_, err := opt.Parse([]string{"--fl"})
	if err == nil {
		t.Errorf("Ambiguous argument 'fl' didn't raise unknown option error")
	}
	if err != nil && err.Error() != "Unknown option 'fl'" {
		t.Errorf("Error string didn't match expected value")
	}
}

func TestGetOptString(t *testing.T) {
	setup := func() *GetOpt {
		opt := GetOptions()
		opt.String("string")
		return opt
	}

	cases := []struct {
		opt    *GetOpt
		option string
		input  []string
		value  string
	}{
		{setup(),
			"string",
			[]string{"--string=hello"},
			"hello",
		},
		{setup(),
			"string",
			[]string{"--string=hello", "world"},
			"hello",
		},
		{setup(),
			"string",
			[]string{"--string", "hello"},
			"hello",
		},
		{setup(),
			"string",
			[]string{"--string", "hello", "world"},
			"hello",
		},
		// TODO: Set a flag to decide wheter or not to allow this
		{setup(),
			"string",
			[]string{"--string", "--hello", "world"},
			"--hello",
		},
		// TODO: Set up a flag to decide wheter or not to err on this
		{setup(),
			"string",
			[]string{"--string", "hello", "--string", "world"},
			"world",
		},
	}
	for _, c := range cases {
		_, err := c.opt.Parse(c.input)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if c.opt.Option[c.option] != c.value {
			t.Errorf("Wrong value: %v != %v", c.opt.Option[c.option], c.value)
		}
	}
}

func TestGetOptInt(t *testing.T) {
	setup := func() *GetOpt {
		opt := GetOptions()
		opt.Int("int")
		return opt
	}

	cases := []struct {
		opt    *GetOpt
		option string
		input  []string
		value  int
	}{
		{setup(),
			"int",
			[]string{"--int=123"},
			123,
		},
		{setup(),
			"int",
			[]string{"--int=123", "world"},
			123,
		},
		{setup(),
			"int",
			[]string{"--int", "123"},
			123,
		},
		{setup(),
			"int",
			[]string{"--int", "123", "world"},
			123,
		},
	}
	for _, c := range cases {
		_, err := c.opt.Parse(c.input)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if c.opt.Option[c.option] != c.value {
			t.Errorf("Wrong value: %v != %v", c.opt.Option[c.option], c.value)
		}
	}
}

func TestGetOptStringRepeat(t *testing.T) {
	setup := func() *GetOpt {
		opt := GetOptions()
		opt.StringRepeat("string")
		return opt
	}

	cases := []struct {
		opt    *GetOpt
		option string
		input  []string
		value  []string
	}{
		{setup(),
			"string",
			[]string{"--string=hello"},
			[]string{"hello"},
		},
		{setup(),
			"string",
			[]string{"--string=hello", "world"},
			[]string{"hello"},
		},
		{setup(),
			"string",
			[]string{"--string", "hello"},
			[]string{"hello"},
		},
		{setup(),
			"string",
			[]string{"--string", "hello", "world"},
			[]string{"hello"},
		},
		// TODO: Set a flag to decide wheter or not to allow this
		{setup(),
			"string",
			[]string{"--string", "--hello", "world"},
			[]string{"--hello"},
		},
		{setup(),
			"string",
			[]string{"--string", "hello", "--string", "happy", "--string", "world"},
			[]string{"hello", "happy", "world"},
		},
	}
	for _, c := range cases {
		_, err := c.opt.Parse(c.input)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(c.opt.Option[c.option], c.value) {
			t.Errorf("Wrong value: %v != %v", c.opt.Option[c.option], c.value)
		}
	}
}

// TODO: Allow passig : as the map divider
func TestGetOptStringMap(t *testing.T) {
	setup := func() *GetOpt {
		opt := GetOptions()
		opt.StringMap("string")
		return opt
	}

	// TODO: Check error when there is no equal sign.

	cases := []struct {
		opt    *GetOpt
		option string
		input  []string
		value  map[string]string
	}{
		{setup(),
			"string",
			[]string{"--string=hello=world"},
			map[string]string{"hello": "world"},
		},
		{setup(),
			"string",
			[]string{"--string=hello=happy", "world"},
			map[string]string{"hello": "happy"},
		},
		{setup(),
			"string",
			[]string{"--string", "hello=world"},
			map[string]string{"hello": "world"},
		},
		{setup(),
			"string",
			[]string{"--string", "hello=happy", "world"},
			map[string]string{"hello": "happy"},
		},
		// TODO: Set a flag to decide wheter or not to allow this
		{setup(),
			"string",
			[]string{"--string", "--hello=happy", "world"},
			map[string]string{"--hello": "happy"},
		},
		{setup(),
			"string",
			[]string{"--string", "hello=world", "--string", "key=value", "--string", "key2=value2"},
			map[string]string{"hello": "world", "key": "value", "key2": "value2"},
		},
	}
	for _, c := range cases {
		_, err := c.opt.Parse(c.input)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(c.opt.Option[c.option], c.value) {
			t.Errorf("Wrong value: %v != %v", c.opt.Option[c.option], c.value)
		}
	}
}

func TestVars(t *testing.T) {
	opt := GetOptions()
	var flag, nflag, nflag2 bool
	var str string
	var integer int
	opt.FlagVar(&flag, "flag")
	opt.NFlagVar(&nflag, "nflag")
	opt.NFlagVar(&nflag2, "n2")
	opt.StringVar(&str, "stringVar")
	opt.IntVar(&integer, "intVar")

	_, err := opt.Parse([]string{
		"-f",
		"-nf",
		"--no-n2",
		"--stringVar", "hello",
		"--intVar", "123",
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if flag != true {
		t.Errorf("flag didn't have expected value: %v != %v", flag, true)
	}
	if nflag != true {
		t.Errorf("nflag didn't have expected value: %v != %v", nflag, true)
	}
	if nflag2 != false {
		t.Errorf("nflag2 didn't have expected value: %v != %v", nflag2, false)
	}
	if str != "hello" {
		t.Errorf("str didn't have expected value: %v != %v", str, "hello")
	}
	if integer != 123 {
		t.Errorf("integer didn't have expected value: %v != %v", integer, 123)
	}
}

func TestDefaultValues(t *testing.T) {
	var flag, nflag bool
	var str string
	var integer int

	opt := GetOptions()
	opt.Flag("flag")
	opt.FlagVar(&flag, "varflag")
	opt.NFlag("nflag")
	opt.NFlagVar(&nflag, "varnflag")
	opt.String("string")
	opt.StringVar(&str, "stringVar")
	opt.Int("int")
	opt.IntVar(&integer, "intVar")
	opt.StringRepeat("string-repeat")
	opt.StringMap("string-map")

	_, err := opt.Parse([]string{})

	if err != nil {
		log.Println(err)
	}

	expected := map[string]interface{}{
		"flag":          nil,
		"nflag":         nil,
		"string":        nil,
		"int":           nil,
		"string-repeat": nil,
		"string-map":    nil,
	}

	for k := range expected {
		if !reflect.DeepEqual(opt.Option[k], expected[k]) {
			t.Errorf("Wrong value: %v != %v", opt.Option, expected)
		}
	}

	if flag != false {
		t.Errorf("flag didn't have expected value: %v != %v", flag, true)
	}
	if nflag != false {
		t.Errorf("nflag didn't have expected value: %v != %v", nflag, true)
	}
	if str != "" {
		t.Errorf("str didn't have expected value: %v != %v", str, "hello")
	}
	if integer != 0 {
		t.Errorf("integer didn't have expected value: %v != %v", integer, 123)
	}

	// Tested above, but it gives me a feel for how it would be used

	if opt.Option["flag"] != nil && !opt.Option["flag"].(bool) {
		t.Errorf("flag didn't have expected value: %v != %v", opt.Option["flag"], nil)
	}
	if opt.Option["non-used-flag"] != nil && opt.Option["non-used-flag"].(bool) {
		t.Errorf("non-used-flag didn't have expected value: %v != %v", opt.Option["non-used-flag"], nil)
	}
	if opt.Option["flag"] != nil && opt.Option["nflag"].(bool) {
		t.Errorf("nflag didn't have expected value: %v != %v", opt.Option["nflag"], nil)
	}
	if opt.Option["string"] != nil {
		t.Errorf("str didn't have expected value: %v != %v", opt.Option["string"], nil)
	}
	if opt.Option["int"] != nil {
		t.Errorf("int didn't have expected value: %v != %v", opt.Option["int"], nil)
	}
}

func TestAll(t *testing.T) {
	var flag, nflag, nflag2 bool
	var str string
	var integer int
	opt := GetOptions()
	opt.Flag("flag")
	opt.FlagVar(&flag, "varflag")
	opt.Flag("non-used-flag")
	opt.NFlag("nflag")
	opt.NFlag("nftrue")
	opt.NFlag("nfnil")
	opt.NFlagVar(&nflag, "varnflag")
	opt.NFlagVar(&nflag2, "varnflag2")
	opt.String("string")
	opt.StringVar(&str, "stringVar")
	opt.Int("int")
	opt.IntVar(&integer, "intVar")
	opt.StringRepeat("string-repeat")
	opt.StringMap("string-map")

	// log.Println(opt)

	remaining, err := opt.Parse([]string{
		"hello",
		"--flag",
		"--varflag",
		"--no-nflag",
		"--nft",
		"happy",
		"--varnflag",
		"--no-varnflag2",
		"--string", "hello",
		"--stringVar", "hello",
		"--int", "123",
		"--intVar", "123",
		"--string-repeat", "hello", "--string-repeat", "world",
		"--string-map", "hello=world", "--string-map", "server=name",
		"world",
	})

	if err != nil {
		log.Println(err)
	}

	if !reflect.DeepEqual(remaining, []string{"hello", "happy", "world"}) {
		t.Errorf("remaining didn't have expected value: %v != %v", remaining, []string{"hello", "happy", "world"})
	}

	expected := map[string]interface{}{
		"flag":          true,
		"nflag":         false,
		"nftrue":        true,
		"string":        "hello",
		"int":           123,
		"string-repeat": []string{"hello", "world"},
		"string-map":    map[string]string{"hello": "world", "server": "name"},
	}

	for k := range expected {
		if !reflect.DeepEqual(opt.Option[k], expected[k]) {
			t.Errorf("Wrong value: %v != %v", opt.Option, expected)
		}
	}

	if flag != true {
		t.Errorf("flag didn't have expected value: %v != %v", flag, true)
	}
	if nflag != true {
		t.Errorf("nflag didn't have expected value: %v != %v", nflag, true)
	}
	if nflag2 != false {
		t.Errorf("nflag2 didn't have expected value: %v != %v", nflag2, false)
	}
	if str != "hello" {
		t.Errorf("str didn't have expected value: %v != %v", str, "hello")
	}
	if integer != 123 {
		t.Errorf("int didn't have expected value: %v != %v", integer, 123)
	}

	// Tested above, but it gives me a feel for how it would be used

	if opt.Option["flag"] != nil && !opt.Option["flag"].(bool) {
		t.Errorf("flag didn't have expected value: %v != %v", opt.Option["flag"], true)
	}
	if opt.Option["non-used-flag"] != nil && opt.Option["non-used-flag"].(bool) {
		t.Errorf("non-used-flag didn't have expected value: %v != %v", opt.Option["non-used-flag"], false)
	}
	if opt.Option["flag"] != nil && opt.Option["nflag"].(bool) {
		t.Errorf("nflag didn't have expected value: %v != %v", opt.Option["nflag"], true)
	}
	if opt.Option["string"] != "hello" {
		t.Errorf("str didn't have expected value: %v != %v", opt.Option["string"], "hello")
	}
	if opt.Option["int"] != 123 {
		t.Errorf("int didn't have expected value: %v != %v", opt.Option["int"], 123)
	}
}