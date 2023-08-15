/*
Copyright © 2023 Jose Cueto

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	portForwardCmdArgs struct {
		port           string
		ocSvcNamespace string
		svcSource      string
		svcName        string
	}
)

// portForwardCmd represents the portForward command
var portForwardCmd = &cobra.Command{
	Use:   "portForward",
	Short: "Forward service port to container port",
	Run: func(cmd *cobra.Command, args []string) {
		portArg := strings.TrimSpace(portForwardCmdArgs.port)
		ports := strings.Split(portArg, ":")

		if len(ports) != 2 {
			log.Fatalf("Port map should be in the form of \"containerport:serviceport\"")
		}

		cmdArgs := []string{
			"port-forward",
			portForwardCmdArgs.svcSource,
			"--address",
			"0.0.0.0",
			fmt.Sprintf("%s:%s", ports[0], ports[1]),
			"-n",
			portForwardCmdArgs.ocSvcNamespace,
		}
		log.Printf("Forwarding %s/%s port to %s", portForwardCmdArgs.svcName, portForwardCmdArgs.ocSvcNamespace, ports[0])

		log.Printf("Running command: %s %s\n", "oc", cmdArgs)

		execCmd := exec.Command("oc", cmdArgs...)
		err := execCmd.Start()

		if err != nil {
			log.Fatal(err)
		}

		go func() {
			err = execCmd.Wait()
			log.Printf("Command finished with error: %v", err)
		}()
	},
}

func init() {
	rootCmd.AddCommand(portForwardCmd)

	flags := portForwardCmd.Flags()
	flags.StringVarP(
		&portForwardCmdArgs.port,
		"ports",
		"p",
		"",
		"Container to service port map or \"container:port\"",
	)

	flags.StringVarP(
		&portForwardCmdArgs.ocSvcNamespace,
		"ocSvcNamespace",
		"n",
		"",
		"The namespace where the service resides.",
	)

	flags.StringVarP(
		&portForwardCmdArgs.svcSource,
		"svcSource",
		"s",
		"",
		"The source name (e.g. pod name) of the port forward from the service side.",
	)

	flags.StringVarP(
		&portForwardCmdArgs.svcName,
		"svcName",
		"m",
		"",
		"An arbitrary name of the service.",
	)

	portForwardCmd.MarkFlagRequired("svcName")
	portForwardCmd.MarkFlagRequired("svcSource")
	portForwardCmd.MarkFlagRequired("ocSvcNamespace")
	portForwardCmd.MarkFlagRequired("ports")

}
