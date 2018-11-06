package command

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning-web/lib/apply"
	"github.com/syou6162/go-active-learning-web/lib/add"
	"github.com/syou6162/go-active-learning-web/lib/web"
)

var Commands = []cli.Command{
	apply.CommandApply,
	add.CommandAdd,
	web.CommandServe,
}
