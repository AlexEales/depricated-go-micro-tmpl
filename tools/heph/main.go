package main

import (
	"github.com/AlexEales/go-micro-tmpl/tools/heph/cmd"
	"github.com/sirupsen/logrus"
)

type Command func() error

var (
	// TODO: not a fan of the logrus formatting should make a common lib with a common format
	log = logrus.New()
)

func main() {
	cmd.Execute()
}
