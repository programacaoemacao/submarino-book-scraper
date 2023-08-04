package flags

import (
	gflag "github.com/jessevdk/go-flags"
)

type RunOptions struct {
	URLToCollect string `short:"u" long:"url" description:"submarino url to scrape" required:"true"`
	Output       string `short:"o" long:"output" description:"output file (supported: [json])" required:"true"`
}

func GetOptions(args ...string) (*RunOptions, error) {
	runOptions := new(RunOptions)

	_, err := gflag.ParseArgs(runOptions, args)
	if err != nil {
		return nil, err
	}

	return runOptions, nil
}
