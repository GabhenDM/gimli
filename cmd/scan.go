package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gabhendm/gimli/service"
	"github.com/gabhendm/gimli/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ScanCmd)
	ScanCmd.Flags().StringP("domain-scan", "d", "example.com", "Set the TLD to be scanned")
	ScanCmd.MarkFlagRequired("domain-scan")
}

// ScanCmd is used to run a scan on a TLD especified via de -d flag
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Start a scan by specifying a TLD",
	Long:  `This command starts a Gimli scan with the specified TLD`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		domain, _ := cmd.Flags().GetString("domain-scan")
		fmt.Println(fmt.Sprintf("[!] Starting Subfinder scan for domain: %s...", domain))
		output, err := service.StartContainer("projectdiscovery/subfinder", []string{"-d", domain, "-oJ", "> output.json"}, debug)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(output)

		var subOutput types.SubfinderOutput

		for scanner.Scan() {
			var el types.SubfinderHost
			json.Unmarshal(scanner.Bytes(), &el)
			subOutput.Subdomains = append(subOutput.Subdomains, el)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		test, err := json.Marshal(subOutput)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(test))

	},
}
