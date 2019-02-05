package corehttp

import (
	"crypto/subtle"
	"errors"

	"github.com/ipsn/go-ipfs/core/commands"
	"github.com/ipsn/go-ipfs/core/commands/cmdenv"
	cmds "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-cmds"
)

func NewAuthCommand(cmd *cmds.Command) *cmds.Command {
	res := *cmd

	res.Subcommands = make(map[string]*cmds.Command)
	for k, v := range cmd.Subcommands {
		res.Subcommands[k] = NewAuthCommand(v)
	}

	res.Run = func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		if cfg.API.Auth != "" {
			auth, err := getAuth(req)
			if err != nil {
				return err
			}

			if subtle.ConstantTimeCompare([]byte(auth), []byte(cfg.API.Auth)) != 1 {
				return errors.New("invalid authentication token")
			}
		}

		return cmd.Run(req, res, env)
	}

	return &res
}

func getAuth(req *cmds.Request) (string, error) {
	auth, ok := req.Options[commands.AuthOption].(string)
	if ok {
		delete(req.Options, commands.AuthOption)
		return auth, nil
	}

	auth, ok = req.Options[commands.ConfigOption].(string)
	if ok {
		delete(req.Options, commands.ConfigOption)
		return auth, nil
	}

	return "", errors.New("no authentication token")
}
