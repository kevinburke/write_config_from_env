package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const Version = "0.1"

type whitelistArg []string

func (w *whitelistArg) String() string {
	return fmt.Sprint(*w)
}

func (w *whitelistArg) Set(value string) error {
	*w = append(*w, value)
	return nil
}

func init() {
	flag.Usage = func() {
		os.Stderr.WriteString(`write_config_from_env

Read configuration from environment variables and write it to a yml file. By
default this script prints the config to stdout. Pass --config=<file> to write
to a file instead.

Usage of write_config_from_env:
`)
		flag.PrintDefaults()
	}
}

func checkErr(err error, activity string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s: %s\n", activity, err.Error())
		os.Exit(2)
	}
}

// http://stackoverflow.com/a/22235064/329700
var specialChars = ":{}[],&*#?|-<>=!%@\\"

func escape(arg string) string {
	if strings.ContainsAny(arg, specialChars) {
		// this is probably incorrect
		return "'" + strings.Replace(arg, "'", "\\'", -1) + "'"
	}
	return arg
}

func writeConfig(b *bytes.Buffer, environ []string, whitelist map[string]bool) {
	for i := range environ {
		// IndexByte would be faster but performance or allocations are not
		// really a problem here
		parts := strings.SplitN(environ[i], "=", 2)
		if len(parts) < 2 {
			continue
		}
		if len(whitelist) > 0 && whitelist[parts[0]] == false {
			continue
		}
		b.WriteString(strings.ToLower(parts[0]))
		b.WriteByte(':')
		if strings.IndexByte(parts[1], ',') >= 0 {
			b.WriteByte('\n')
			args := strings.Split(parts[1], ",")
			for j := range args {
				if args[j] == "" {
					continue
				}
				b.WriteString("  - ")
				b.WriteString(escape(args[j]))
				b.WriteByte('\n')
			}
		} else {
			b.WriteByte(' ')
			b.WriteString(escape(parts[1]))
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
}

func main() {
	whitelistFlag := new(whitelistArg)
	flag.Var(whitelistFlag, "whitelist", "Environment variables to whitelist. If unspecified, all environment variables will be written to the config")
	cfg := flag.String("config", "", "Path to a config file")
	flag.Parse()
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "write_config_from_env: Too many arguments\n")
		os.Exit(2)
	}
	if flag.NArg() == 1 {
		switch flag.Arg(0) {
		case "version":
			fmt.Fprintf(os.Stderr, "write_config_from_env version %s\n", Version)
			os.Exit(2)
		case "help":
			flag.Usage()
			os.Exit(2)
		default:
			fmt.Fprintf(os.Stderr, "Unknown argument: %s\n", flag.Arg(0))
			os.Exit(2)
		}
	}
	b := new(bytes.Buffer)
	whitelistMap := make(map[string]bool, len(*whitelistFlag))
	for i := range *whitelistFlag {
		whitelistMap[(*whitelistFlag)[i]] = true
	}
	writeConfig(b, os.Environ(), whitelistMap)
	var w io.Writer
	if *cfg == "" {
		w = os.Stdout
	} else {
		f, err := os.OpenFile(*cfg, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0400)
		checkErr(err, "creating config file")
		defer f.Close()
		w = f
	}
	_, err := io.Copy(w, b)
	checkErr(err, "writing config file")
}
