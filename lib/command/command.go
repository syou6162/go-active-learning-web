package command

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning-web/lib/apply"
	"github.com/syou6162/go-active-learning-web/lib/expand_url"
	"github.com/syou6162/go-active-learning-web/lib/web"
)

var Commands = []cli.Command{
	apply.CommandApply,
	expand_url.CommandExpandURL,
	web.CommandServe,
}
