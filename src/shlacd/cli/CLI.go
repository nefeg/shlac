package cli


import (
	"time"
	"github.com/urfave/cli"
)

func New() *cli.App{

	Cli := cli.NewApp()
	Cli.Version             = "0.4-rc1"
	Cli.Name                = "ShLAC(client)"
	Cli.Usage               = "[SH]lac [L]ike [A]s [C]ron"
	Cli.Author              = "Evgeny Nefedkin"
	Cli.Compiled            = time.Now()
	Cli.Email               = "evgeny.nefedkin@umbrella-web.com"
	Cli.EnableBashCompletion= true
	Cli.Description         = "Distributed and concurrency Job manager\n" +

		"\t\tSupported extended syntax:\n" +
		"\t\t------------------------------------------------------------------------\n" +
		"\t\tField name     Mandatory?   Allowed values    Allowed special characters\n" +
		"\t\t----------     ----------   --------------    --------------------------\n" +
		"\t\tSeconds        No           0-59              * / , -\n" +
		"\t\tMinutes        Yes          0-59              * / , -\n" +
		"\t\tHours          Yes          0-23              * / , -\n" +
		"\t\tDay of month   Yes          1-31              * / , - L W\n" +
		"\t\tMonth          Yes          1-12 or JAN-DEC   * / , -\n" +
		"\t\tDay of week    Yes          0-6 or SUN-SAT    * / , - L #\n" +
		"\t\tYear           No           1970–2099         * / , -\n" +

		"\n\n" +

		"\t\tand aliases:\n" +
		"\t\t-------------------------------------------------------------------------------------------------\n" +
		"\t\tEntry       Description                                                             Equivalent to\n" +
		"\t\t-------------------------------------------------------------------------------------------------\n" +
		"\t\t@annually   Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *\n" +
		"\t\t@yearly     Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *\n" +
		"\t\t@monthly    Run once a month at midnight in the morning of the first of the month   0 0 0 1 * * *\n" +
		"\t\t@weekly     Run once a week at midnight in the morning of Sunday                    0 0 0 * * 0 *\n" +
		"\t\t@daily      Run once a day at midnight                                              0 0 0 * * * *\n" +
		"\t\t@hourly     Run once an hour at the beginning of the hour                           0 0 * * * * *\n" +
		"\t\t@reboot     Not supported"

	return Cli
}