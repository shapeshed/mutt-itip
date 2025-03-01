package cmd

import (
	"fmt"
	"os"
	"path"

	ics "github.com/arran4/golang-ical"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates the invitation in question",
	Long:  "update currontly supports invitation cancellation only",
	Run: func(cmd *cobra.Command, args []string) {
		_, c, err := parseInput(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		err = updateLocalCalendar(c)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	updateCmd.Flags().BoolVar(&dry, "dry", false, "if set update does not change local calendar")
	rootCmd.AddCommand(updateCmd)
}

func updateLocalCalendar(c *ics.Calendar) error {
	e := oneEventOrDie(c)
	uid := e.GetProperty(ics.ComponentPropertyUniqueId).Value

	m := extractMethod(c)
	var err error
	switch m {
	case ics.MethodCancel:
		err = deleteICSFile(uid)
	default:
		err = fmt.Errorf("unsupported method: %s", m)
	}
	return err
}

func deleteICSFile(uid string) error {
	icsDir := viper.GetString("icsDir")
	if icsDir == "" {
		log.Fatalf("unable to update local calendar: you need to set icsDir in %s", viper.ConfigFileUsed())
	}
	icsDir = os.ExpandEnv(icsDir)

	f := path.Join(icsDir, uid) + ".ics"
	if dry {
		log.Infof("would remove file: %s", f)
	} else {
		log.Infof("removing file: %s", f)
		return os.Remove(f)
	}
	return nil
}
