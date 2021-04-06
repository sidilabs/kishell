package options

import (
	"time"

	"github.com/alecthomas/kong"
	"github.com/sidilabs/kishell/pkg/config"
	"github.com/sidilabs/kishell/pkg/utils"
)

// Context configuration.
type Context struct {
	Debug         bool
	Configuration config.Configuration
}

// Option represents a common configuration for all options.
type Option struct {
	Context    *kong.Context
	ConfigFile config.Configuration
}

// ConfigureCmd represents CLI arguments for configure option.
type ConfigureCmd struct {
	Server bool `optional help:"Add a new server definition"`
	Role   bool `optional help:"Add a new role definition"`
	Reset  bool `optional help:"Reset the whole configuration"`
}

// UseCmd represents CLI arguments for use option.
type UseCmd struct {
	Server string `optional help:"Set what server to use when querying ES"`
	Role   string `optional help:"Set what role to use when querying ES"`
}

// ListCmd represents CLI arguments for list option.
type ListCmd struct {
}

// SearchCmd represents CLI arguments for search option.
type SearchCmd struct {
	Query      string           `optional help:"Text input to query data. Use the same format as you would use in Kibana"`
	Older      string           `optional default:"now" help:"Data older than. Defaults to current time when not provided (e.g. 30m, 1h, 1w, 1M, 1y)"`
	Newer      string           `optional default:"15m" help:"Data newer than. Defaults to 15m when not provided (e.g. 30m, 1h, 1w, 1M, 1y)"`
	Limit      int32            `optional default:"50" help:"Limit the number of messages fetched"`
	Server     string           `optional help:"Which server to query against. Used to override the current server config"`
	httpClient utils.HTTPClient `-`
}

// CLI represents possible CLI options.
var CLI struct {
	Debug     bool         `help:"Enable debug mode."`
	Configure ConfigureCmd `cmd help:"Init ES server configs"`
	List      ListCmd      `cmd help:"Show the current server configs"`
	Search    SearchCmd    `cmd help:"Search for data"`
	Use       UseCmd       `cmd help:"Update config options with ser/role preferences"`
}

// OlderAsTimestamp converts a ISO-8601 period as string in timestamp.
func (s *SearchCmd) OlderAsTimestamp() (int64, error) {
	return toTimestamp(s.Older)
}

// NewerAsTimestamp converts a ISO-8601 period as string in timestamp.
func (s *SearchCmd) NewerAsTimestamp() (int64, error) {
	return toTimestamp(s.Newer)
}

// AfterApply defines the http client instance to used once search option is identified to take execution.
func (s *SearchCmd) AfterApply(h *utils.DefaultHTTPClient) error {
	s.httpClient = h
	return nil
}

func toTimestamp(period string) (int64, error) {
	now := time.Now().Unix() * 1000
	if len(period) <= 0 || period == "now" {
		return now, nil
	}
	duration, err := time.ParseDuration(period)
	if err != nil {
		return -1, err
	}
	return now - duration.Milliseconds(), nil
}

// Run the option found through CLI arguments.
func (o *Option) Run() {
	err := o.Context.Run(&Context{
		Debug:         CLI.Debug,
		Configuration: o.ConfigFile,
	})
	o.Context.FatalIfErrorf(err)
}

// Parse the CLI arguments.
func Parse() Option {
	httpClient := &utils.DefaultHTTPClient{
		Timeout: 5 * time.Second,
	}
	context := kong.Parse(&CLI, kong.Bind(httpClient))
	opt := Option{
		Context:    context,
		ConfigFile: config.LoadDefaultConfig(),
	}
	return opt
}
