package portald

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func (s *Service) handlers() cmd.Handlers {
	return cmd.Handlers{
		cli.Handler,
		cli.EncodingHandler,
		{
			Func: s.Api,
			Name: "api",
			Desc: "Print API schema.",
		},
		{
			Func: s.Ping,
			Name: "ping",
		},
		{
			Func: Join,
			Name: "join",
		},
		{
			Func: s.Open().Start,
			Name: "open o",
			Desc: "Open portal app.",
			Params: cmd.Params{
				{Type: "string", Desc: "Absolute path to app bundle or directory."},
				{Type: "...string", Desc: "Optional arguments."},
			},
		},
	}.Plus(s.publicHandlers()...)
}

func (s *Service) publicHandlers() cmd.Handlers {
	return cmd.Handlers{
		{
			Func: s.Setup,
			Name: "setup",
			Desc: "Setup portal environment.",
			Params: cmd.Params{
				{Name: "user u", Type: "string", Desc: "Optional user alias. When specified, the installed node will be assigned to a new user identity associated with the name. Otherwise, the installed node will be ready to claim by existing user."},
				{Name: "apps a", Type: "string", Desc: "Path to directory containing application to install."},
			},
		},
		{
			Name: "user u",
			Desc: "Manage user.",
			Sub: cmd.Handlers{
				{
					Func: s.CreateUser,
					Name: "create",
					Desc: "Create user.",
					Params: cmd.Params{
						{Type: "string", Desc: "alias"},
					},
				},
				{
					Func: s.Claim,
					Name: "claim",
					Desc: "Claim user.",
					Params: cmd.Params{
						{Type: "string", Desc: "alias"},
					},
				},
				{
					Func: s.UserInfo,
					Name: "info i",
					Desc: "Print user info.",
				},
				{
					Func: s.HasUser,
					Name: "check",
					Desc: "Check if user exists.",
				},
			},
		},
		{
			Func: s.Connect,
			Name: "connect c",
			Desc: "Open portal app and redirect standard IO to client.",
			Params: cmd.Params{
				{Type: "string", Desc: "Absolute path to app bundle or directory."},
				{Type: "...string", Desc: "Optional arguments."},
			},
		},
		{
			Name: "token t",
			Desc: "Manage tokens.",
			Sub: cmd.Handlers{
				{
					Func: s.Tokens().Resolve,
					Name: "create c",
					Desc: "Create a new token or return existing one.",
					Params: cmd.Params{
						{Type: "string", Desc: "Token name."},
					},
				},
				{
					Func: s.Tokens().List,
					Name: "list l",
					Desc: "List apphost tokens.",
					Params: cmd.Params{
						{Name: "format f", Type: "string", Desc: "Format [json, bin]"},
					},
				},
			},
		},
		{
			Name: "app a",
			Desc: "Manage applications.",
			Sub: cmd.Handlers{
				{
					Func: s.InstallApp,
					Name: "install i",
					Desc: "Install app.",
					Params: cmd.Params{
						{Type: "string", Desc: "App name, package name, or path to the directory containing app bundle."},
					},
				},
				{
					Func: s.ClaimPackage,
					Name: "claim c",
					Desc: "Claim the app by package name.",
					Params: cmd.Params{
						{Type: "string", Desc: "App package name"},
					},
				},
				{
					Func: s.Installer().Uninstall,
					Name: "uninstall d",
					Desc: "Uninstall app.",
					Params: cmd.Params{
						{Type: "string", Desc: "Application name or package name"},
					},
				},
				{
					Func: s.InstalledApps,
					Name: "list l",
					Desc: "List installed apps.",
					Params: cmd.Params{
						{Name: "hidden h", Type: "boolean", Desc: "Include hidden apps."},
					},
					Sub: cmd.Handlers{{
						Func: s.ObserveApps,
						Name: "observe o",
						Desc: "Observe installed list apps.",
						Params: cmd.Params{
							{Name: "hidden h", Type: "boolean", Desc: "Include hidden apps."},
						},
					}},
				},
				{
					Func: s.ObserveApps,
					Name: "observe o",
					Desc: "Observe installed list apps.",
					Params: cmd.Params{
						{Name: "hidden h", Type: "boolean", Desc: "Include hidden apps."},
					},
				},
				{
					Func: s.AvailableApps,
					Name: "available a",
					Desc: "List available apps.",
					Params: cmd.Params{
						{Type: "boolean", Desc: "subscribe."},
					},
				},
				{
					Func: s.PublishApps,
					Name: "publish p",
					Desc: "Publish app bundles.",
					Params: cmd.Params{
						{Type: "string", Desc: "Path to app bundle or directory."},
					},
				},
			},
		},
		{
			Name: "node n",
			Desc: "Manage nodes.",
			Sub: cmd.Handlers{
				{
					Func: s.nodeInfo,
					Name: "info i",
					Desc: "Print node info.",
				},
			},
		},
		{
			Name: "fs",
			Desc: "Manage fs module.",
			Sub: cmd.Handlers{
				{
					Name: "config c",
					Func: s.FsConfigRead,
					Desc: "Read fs config.",
				},
				{
					Name: "watch w",
					Desc: "Read only directories.",
					Sub: cmd.Handlers{
						{
							Func: s.FsConfigWatchAdd,
							Name: "add a",
							Desc: "Observe given directory.",
							Params: cmd.Params{
								{Type: "string", Desc: "Alias."},
								{Type: "string", Desc: "Path to dir."},
							},
						},
					},
				},
			},
		},
		{
			Func: s.Stop,
			Name: "close",
			Desc: "Shutdown portal environment and close all running apps.",
		},
	}
}
