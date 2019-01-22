package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/containers/buildah"
	ispecs "github.com/opencontainers/image-spec/specs-go"
	rspecs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/cobra"
)

//Overwritten at build time
var (
	GitCommit  string
	buildInfo  string
	cniVersion string
)

//Function to get and print info for version command
func versionCmd(c *cobra.Command, args []string) error {
	var err error
	buildTime := int64(0)
	if buildInfo != "" {
		//converting unix time from string to int64
		buildTime, err = strconv.ParseInt(buildInfo, 10, 64)
		if err != nil {
			return err
		}
	}

	fmt.Println("Version:        ", buildah.Version)
	fmt.Println("Go Version:     ", runtime.Version())
	fmt.Println("Image Spec:     ", ispecs.Version)
	fmt.Println("Runtime Spec:   ", rspecs.Version)
	fmt.Println("CNI Spec:       ", cniversion.Current())
	fmt.Println("libcni Version: ", cniVersion)
	fmt.Println("Git Commit:     ", GitCommit)

	//Prints out the build time in readable format
	fmt.Println("Built:          ", time.Unix(buildTime, 0).Format(time.ANSIC))
	fmt.Println("OS/Arch:        ", runtime.GOOS+"/"+runtime.GOARCH)

	return nil
}

//cli command to print out the version info of buildah
var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display the Buildah version information",
	Long:  "Displays Buildah version information.",
	RunE:  versionCmd,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(versionCommand)
}
