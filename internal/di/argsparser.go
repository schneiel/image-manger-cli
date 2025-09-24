package di

// ArgumentParser extracts configuration values from command line arguments.
type ArgumentParser struct{}

// NewArgumentParser creates a new argument parser.
func NewArgumentParser() *ArgumentParser {
	return &ArgumentParser{}
}

// ExtractLanguage extracts the language flag from command line arguments.
func (p *ArgumentParser) ExtractLanguage(args []string) string {
	for i, arg := range args {
		if (arg == "--lang" || arg == "-l" || arg == "--language") && i+1 < len(args) {
			return args[i+1]
		}
	}

	return "en"
}

// ExtractConfigPath extracts the config path flag from command line arguments.
func (p *ArgumentParser) ExtractConfigPath(args []string) string {
	for i, arg := range args {
		if (arg == "--config" || arg == "-c") && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}
