package getoptions

import (
	"regexp"
	"strings"
)

// 1: leading dashes
// 2: option
// 3: =arg
var isOptionRegex = regexp.MustCompile(`^(--?)([^=]+)(.*?)$`)

type optionPair struct {
	Option string
	// We allow multiple args in case of splitting on comma.
	Args []string
}

/*
isOptionV2 - Enhanced version of isOption, this one returns pairs of options and arguments
At this level we don't agregate results in case we have -- and then other options, basically we can parse one option at a time.
This makes the caller have to agregate multiple calls to the same option.
TODO: Here is where we should handle windows /option types.
*/
func isOption(s string, mode Mode) ([]optionPair, bool) {
	// Handle especial cases
	if s == "--" {
		return []optionPair{{Option: "--"}}, false
	} else if s == "-" {
		return []optionPair{{Option: "-"}}, true
	}

	match := isOptionRegex.FindStringSubmatch(s)
	if len(match) > 0 {
		// check long option
		if match[1] == "--" {
			opt := optionPair{}
			opt.Option = match[2]
			args := strings.TrimPrefix(match[3], "=")
			if args != "" {
				// TODO: Here is where we could split on comma
				opt.Args = []string{args}
			}
			return []optionPair{opt}, true
		}
		// check short option
		switch mode {
		case Bundling:
			opts := []optionPair{}
			for _, option := range strings.Split(match[2], "") {
				opt := optionPair{}
				opt.Option = option
				opts = append(opts, opt)
			}
			if len(opts) > 0 {
				args := strings.TrimPrefix(match[3], "=")
				if args != "" {
					opts[len(opts)-1].Args = []string{args}
				}
			}
			return opts, true
		case SingleDash:
			opts := []optionPair{}
			for _, option := range []string{strings.Split(match[2], "")[0]} {
				opt := optionPair{}
				opt.Option = option
				opts = append(opts, opt)
			}
			if len(opts) > 0 {
				args := strings.Join(strings.Split(match[2], "")[1:], "") + match[3]
				opts[len(opts)-1].Args = []string{args}
			}
			return opts, true
		default:
			opt := optionPair{}
			opt.Option = match[2]
			args := strings.TrimPrefix(match[3], "=")
			if args != "" {
				opt.Args = []string{args}
			}
			return []optionPair{opt}, true
		}
	}
	return []optionPair{}, false
}

// stringSliceIndex - indicates if an element is found in the slice and what its index is
func stringSliceIndex(ss []string, e string) (int, bool) {
	for i, s := range ss {
		if s == e {
			return i, true
		}
	}
	return -1, false
}
