package add

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func doAdd(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "expand-url")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	err := db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	err = cache.Init()
	if err != nil {
		return err
	}
	defer cache.Close()

	examples, err := file.ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	cache.AttachMetadata(examples, true, false)

	for _, e := range examples {
		_, err = db.InsertOrUpdateExample(e)
		if err != nil {
			return err
		}

		url := e.Url
		if e.FinalUrl != "" {
			url = e.FinalUrl
		}
		fmt.Println(fmt.Sprintf("%s\t%d", url, e.Label))
	}

	return nil
}

var CommandAdd = cli.Command{
	Name:  "add",
	Usage: "add urls",
	Description: `
Add urls.
`,
	Action: doAdd,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
	},
}
