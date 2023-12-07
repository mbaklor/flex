package main

import "github.com/urfave/cli/v2"

func flexPack(ctx *cli.Context) error {
	dev, err := CheckForDevice(ctx)
	if err != nil {
		return ShowHelpAndError(ctx, err)
	}
	println(dev.Address.String())
	return nil
}
