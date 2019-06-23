package main

import "github.com/sahilm/fuzzy"

// OptsCompleter creates completer function for ishell from list of opts
func OptsCompleter(opts []string) func(prefix string, args []string) []string {
	return func(prefix string, args []string) []string {
		if prefix == "" {
			return opts
		}

		var completion []string

		matches := fuzzy.Find(prefix, opts)
		for _, match := range matches {
			completion = append(completion, match.Str)
		}

		return completion

	}
}
