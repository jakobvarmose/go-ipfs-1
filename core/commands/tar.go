package commands

import (
	"fmt"
	"io"

	core "github.com/ipsn/go-ipfs/core"
	cmdenv "github.com/ipsn/go-ipfs/core/commands/cmdenv"
	coreiface "github.com/ipsn/go-ipfs/core/coreapi/interface"
	tar "github.com/ipsn/go-ipfs/tar"

	"github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-path"
	cmds "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-cmds"
	dag "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-merkledag"
	cmdkit "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-cmdkit"
)

var TarCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Utility functions for tar files in ipfs.",
	},

	Subcommands: map[string]*cmds.Command{
		"add": tarAddCmd,
		"cat": tarCatCmd,
	},
}

var tarAddCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Import a tar file into ipfs.",
		ShortDescription: `
'ipfs tar add' will parse a tar file and create a merkledag structure to
represent it.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.FileArg("file", true, false, "Tar file to add.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		fi, err := req.Files.NextFile()
		if err != nil {
			return err
		}

		node, err := tar.ImportTar(req.Context, fi, nd.DAG)
		if err != nil {
			return err
		}

		c := node.Cid()

		fi.FileName()
		return res.Emit(&coreiface.AddEvent{
			Name: fi.FileName(),
			Hash: c.String(),
		})
	},
	Type: coreiface.AddEvent{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *coreiface.AddEvent) error {
			fmt.Fprintln(w, out.Hash)
			return nil
		}),
	},
}

var tarCatCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Export a tar file from IPFS.",
		ShortDescription: `
'ipfs tar cat' will export a tar file from a previously imported one in IPFS.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("path", true, false, "ipfs path of archive to export.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		p, err := path.ParsePath(req.Arguments[0])
		if err != nil {
			return err
		}

		root, err := core.Resolve(req.Context, nd.Namesys, nd.Resolver, p)
		if err != nil {
			return err
		}

		rootpb, ok := root.(*dag.ProtoNode)
		if !ok {
			return dag.ErrNotProtobuf
		}

		r, err := tar.ExportTar(req.Context, rootpb, nd.DAG)
		if err != nil {
			return err
		}

		return res.Emit(r)
	},
}
