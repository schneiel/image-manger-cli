// Package main is the entry point for the ImageManager program.
// It parses command-line arguments, initializes internationalization (i18n)
// and the logger, and dispatches execution to the appropriate command module.
package main

import (
	"ImageManager/cmd"
	"ImageManager/i18n"
	"ImageManager/log"
	"flag"
	"fmt"
	"os"
)

// Global flag for language.
var lang = flag.String("lang", "en", "Language for the application output (en/de).")

// main is the primary function that starts the program.
func main() {
	// First, parse the flags to get the language.
	flag.Parse()

	// Initialize i18n.
	if err := i18n.Init(*lang, "i18n/locales"); err != nil {
		fmt.Printf("Fatal: Could not initialize localization: %v\n", err)
		os.Exit(1)
	}

	if err := log.InitFileLogger(); err != nil {
		// Use i18n for the warning.
		fmt.Printf(i18n.T("FileLoggerWarning", map[string]interface{}{"Error": err}))
	}
	defer log.CloseFileLogger()

	// Get arguments after parsing flags.
	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		return
	}

	commands := map[string]cmd.Command{
		"sort":  cmd.NewSortCommand(),
		"dedup": cmd.NewDedupCommand(),
	}

	commandName := args[0]
	command, exists := commands[commandName]

	if !exists {
		log.LogError(i18n.T("UnknownCommand", map[string]interface{}{"CommandName": commandName}))
		printUsage()
		return
	}

	// Pass the remaining arguments to the command.
	command.Init(args[1:])
	if err := command.Run(); err != nil {
		log.LogError(i18n.T("CommandError", map[string]interface{}{"CommandName": command.Name(), "Error": err}))
		os.Exit(1)
	}
}

// printUsage displays usage information for the program.
func printUsage() {
	// Replace all hardcoded strings.
	fmt.Println("\n" + i18n.T("AppName") + " - " + i18n.T("AppDescription"))
	fmt.Printf("\n%s: go run . [--lang=en|de] <%s> [options]\n", i18n.T("Usage"), i18n.T("Commands"))
	fmt.Println("\n" + i18n.T("Commands") + ":")
	fmt.Printf("  %s\n", i18n.T("SortCommandName"))
	fmt.Printf("    %s\n", i18n.T("SortCommandDesc"))
	fmt.Printf("\n  %s\n", i18n.T("DedupCommandName"))
	fmt.Printf("    %s\n", i18n.T("DedupCommandDesc"))
}
