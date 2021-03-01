package main

import (
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Release",
		Short: "Manage OCP Releases locally",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newReleaseRetrieveCmd())
	//cmd.AddCommand(newDestroyClusterCmd())
	return cmd
}

func newReleaseRetrieveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "Retrieve",
		Short: "Retrieve OCP installation media for offline use",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			//	cleanup := setupFileHook(rootOpts.dir)
			//	defer cleanup()
			//
			//	err := runDestroyCmd(rootOpts.dir)
			//	if err != nil {
			//		logrus.Fatal(err)
			//	}
		},
	}
}
