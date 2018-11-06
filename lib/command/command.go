package command

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning-web/lib/recommend"
	"github.com/syou6162/go-active-learning-web/lib/add"
	"github.com/syou6162/go-active-learning-web/lib/web"
)

var Commands = []cli.Command{
	add.CommandAdd,
	recommend.CommandRecommend,
	web.CommandServe,
}
