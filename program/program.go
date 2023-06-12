package program

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"runtime"
	"sync"
)

// Options is the structure of program options
type Options struct {
	Version bool `help:"Show program version"`
	// VersionCmd VersionCmd `name:"version" cmd:"" help:"show program version"`

	KubeConfig      string   `group:"Input" short:"k" help:"Kubeconfig file" type:"path" default:"~/.kube/config"`
	CredentialsFile string   `group:"Input" short:"c" help:"AWS Credentials File" type:"existingfile" default:"~/.aws/credentials"`
	Regions         []string `group:"Input" help:"List of regions to check" env:"AWS_REGIONS" default:"us-east-1,us-east-2,us-west-1,us-west-2,ap-south-1,ap-northeast-3,ap-northeast-2,ap-southeast-1,ap-southeast-2,ap-northeast-1,ca-central-1,eu-central-1,eu-west-1,eu-west-2,eu-west-3,eu-north-1,sa-east-1"`
	Profiles        []string `group:"Input" help:"List of AWS profiles to use.  Will discover profiles if not specified" env:"AWS_PROFILES"`

	Debug        bool   `group:"Info" help:"Show debugging information"`
	OutputFormat string `group:"Info" enum:"auto,jsonl,terminal" default:"auto" help:"How to show program output (auto|terminal|jsonl)"`
	Quiet        bool   `group:"Info" help:"Be less verbose than usual"`
}

// Parse calls the CLI parsing routines
func (program *Options) Parse(args []string) (*kong.Context, error) {
	parser, err := kong.New(program,
		kong.ShortUsageOnError(),
		kong.Description("Download kubeconfigs in bulk by examining clusters across multiple profiles and regions"),
	) 

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return parser.Parse(args)
}

// Run runs the program
func (program *Options) Run(options *Options) error {
	config, err := program.ReadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to read kubeconfig file")
		return err
	}

	clusters := make(chan ClusterInfo)

	wg := sync.WaitGroup{}
	for sess := range program.getUniqueSessions() {
		wg.Add(1)
		go func(sess *sessionInfo) {
			defer wg.Done()
			program.getClustersFrom(sess, clusters)
		}(sess)
	}

	go func() {
		wg.Wait()
		close(clusters)
	}()

	for c := range clusters {
		if err := captureConfig(c, config); err != nil {
			stats.Errors.Add(1)
			log.Error().Err(err).Msg("Error capturing cluster configuration")
		}
	}

	if err := program.WriteConfig(config); err != nil {
		stats.Errors.Add(1)
		log.Error().
			Err(err).
			Str("file", program.KubeConfig).
			Msg("Error saving kubeconfig")
	}

	stats.Log()

	if stats.Errors.Load() > 0 {
		return errors.New("Errors encountered during run")
	}
	return nil
}

// AfterApply runs after the options are parsed but before anything runs
func (program *Options) AfterApply() error {
	program.initLogging()

	if len(program.Regions) < 1 {
		return errors.New("Must specify at least one region")
	}
	return nil
}

func (program *Options) initLogging() {
	if program.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	switch {
	case program.Debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case program.Quiet:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var out io.Writer = os.Stdout

	if os.Getenv("TERM") == "" && runtime.GOOS == "windows" {
		out = colorable.NewColorableStdout()
	}

	if program.OutputFormat == "terminal" ||
		(program.OutputFormat == "auto" && isTerminal(os.Stdout)) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	} else {
		log.Logger = log.Output(out)
	}

	log.Logger.Debug().
		Str("version", Version).
		Str("program", os.Args[0]).
		Msg("Starting")
}

// isTerminal returns true if the file given points to a character device (i.e. a terminal)
func isTerminal(file *os.File) bool {
	if fileInfo, err := file.Stat(); err != nil {
		log.Err(err).Msg("Error running stat")
		return false
	} else {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
}
