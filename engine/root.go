package engine

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:     "viction-datadir-hardlink",
		Short:   "Viction Blockchain data clone via hard link.",
		Version: version(),
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
