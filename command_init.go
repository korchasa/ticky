package main

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type InitCommand struct {
	*Opts
}

func (c *InitCommand) Execute(_ []string) error {
	for _, s := range c.Statuses {
		sp := c.IssuesDir + "/" + s
		if _, err := os.Stat(sp); err == nil || os.IsExist(err) {
			continue
		}

		logrus.Infof("Create status directory `%s`", sp)
		if err := os.MkdirAll(sp, 0755); err != nil {
			return err
		}

		sf := sp + "/.gitkeep"
		logrus.Infof("Wrote git keep file `%s`", sf)
		if err := ioutil.WriteFile(sf, []byte{}, 0644); err != nil {
			return err
		}
	}
	return nil
}
