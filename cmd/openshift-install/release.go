package main

import (
	"github.com/openshift/installer/pkg/release"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Manage OCP Releases locally",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newReleaseCreateCmd())
	cmd.AddCommand(newReleasePushCmd())

	return cmd
}

func newReleaseCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create OCP release objects",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newReleaseCreateBundleCmd())
	return cmd
}

func newReleaseCreateBundleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bundle",
		Short: "Create local OCP release bundle",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			err := release.CreateBundle(rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newReleasePushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Manage uploads of openshift content",
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
