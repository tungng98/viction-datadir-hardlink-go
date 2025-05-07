package engine

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"viction-datadir-clone-go/config"
	"viction-datadir-clone-go/filesystem"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type CloneModule struct {
	config *config.RootConfig
	logger zerolog.Logger
}

func NewCloneModule(c *Controller, cmdName string) *CloneModule {
	return &CloneModule{
		config: c.Root,
		logger: c.CommandLogger("clone", cmdName),
	}
}

func (m *CloneModule) Main(from, to string) error {
	if from == "" || to == "" {
		return errors.New("source and target directory must be specified")
	}
	err := m.clone(filepath.Join(from, "tomo/chaindata"), filepath.Join(to, "tomo/chaindata"))
	if err != nil {
		return err
	}
	err = m.link(filepath.Join(from, "tomo/rewards"), filepath.Join(to, "tomo/rewards"))
	if err != nil {
		return err
	}
	err = m.clone(filepath.Join(from, "tomox"), filepath.Join(to, "tomox"))
	if err != nil {
		return err
	}

	return nil
}

func (m *CloneModule) logError(err error) {
	if err != nil {
		m.logger.Err(err).Msg("Unexpected error has occurred. Program will exit.")
	}
}

func CloneCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "clone",
		Short: "Create another datadir for static node using hard link.",
		Run: func(cmd *cobra.Command, args []string) {
			c := InitApp()
			defer c.Close()
			flags := ParseHardlinkFlags(cmd)
			m := NewCloneModule(c, "clone")
			m.logError(m.Main(flags.From, flags.To))
		},
	}
	rootCmd.Flags().StringP("from", "f", "", "Source directory to hard link from")
	rootCmd.Flags().StringP("to", "t", "", "Target directory to hard link from")

	return rootCmd
}

type HardlinkFlags struct {
	From string
	To   string
}

func ParseHardlinkFlags(cmd *cobra.Command) *HardlinkFlags {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	return &HardlinkFlags{
		From: from,
		To:   to,
	}
}

var skipList = []string{
	"LOCK",
	"LOG",
	"nodekey",
	"transactions.rlp",
}

func (m *CloneModule) clone(fromDir, toDir string) error {
	if !filesystem.IsExist(fromDir) {
		return nil
	}
	contents, err := os.ReadDir(fromDir)

	if err != nil {
		return err
	}

	for _, item := range contents {
		itemName := item.Name()
		if contains(skipList, itemName) {
			continue
		}

		fromItem := filepath.Join(fromDir, item.Name())
		toItem := filepath.Join(toDir, item.Name())
		if item.IsDir() {
			err := m.clone(fromItem, toItem)
			if err != nil {
				return err
			}
			continue
		}

		if filesystem.IsExist(toItem) {
			m.logger.Info().Msgf("File existed: %s", toItem)
			continue
		}

		err := os.MkdirAll(toDir, os.ModeDir)
		if err != nil {
			return err
		}

		if strings.HasSuffix(itemName, ".ldb") || strings.HasSuffix(itemName, ".json") {
			err = os.Link(fromItem, toItem)
			if err != nil {
				return err
			}
			m.logger.Info().Msgf("File linked: %s -> %s", fromItem, toItem)
		} else {
			err = filesystem.CopyFile(fromItem, toItem)
			if err != nil {
				return err
			}
			m.logger.Info().Msgf("File copied: %s -> %s", fromItem, toItem)
		}
	}

	return nil
}

func (m *CloneModule) link(fromDir, toDir string) error {
	contents, err := os.ReadDir(fromDir)

	if err != nil {
		return err
	}

	for _, item := range contents {
		itemName := item.Name()
		if contains(skipList, itemName) {
			continue
		}

		fromItem := filepath.Join(fromDir, item.Name())
		toItem := filepath.Join(toDir, item.Name())
		if item.IsDir() {
			err := m.link(fromItem, toItem)
			if err != nil {
				return err
			}
			continue
		}

		if filesystem.IsExist(toItem) {
			m.logger.Info().Msgf("File existed: %s", toItem)
			continue
		}

		err := os.MkdirAll(toDir, os.ModeDir)
		if err != nil {
			return err
		}

		err = os.Link(fromItem, toItem)
		if err != nil {
			return err
		}
		m.logger.Info().Msgf("File linked: %s -> %s", fromItem, toItem)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
