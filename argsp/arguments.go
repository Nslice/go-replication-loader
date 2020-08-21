package argsp

import (
	"flag"
	"reflect"
)

type stringSlice []string

func (str *stringSlice) String() string {
	return "my string representation"
}

func (str *stringSlice) Set(value string) error {
	*str = append(*str, value)
	return nil
}

// ArgumentOptions provides argument parameters
type ArgumentOptions struct {
	ProjectName string

	// connection flags
	ConsoleServiceName string
	WorkingDirectory   string
	NetPipeServiceName string
	User               string
	Password           string

	// mailing flags
	ToEmailList  stringSlice
	Body         string
	From         string
	SMTPServer   string
	SMTPPort     int
	SMTPLogin    string
	SMTPPassword string

	// database flags
	DbDataSource      string
	DatabaseName      string
	DatabaseUserID    string
	DatabasePassword  string
	BackupPath        string
	UseCompression    bool
	TrustedConnection bool
	ConnectionTimeout int
	// SkipBackup added to skip backup of the second installation.
	// E.g the first installation was failed.
	// So it is not necessary to make backup from brocken database
	SkipBackup bool

	// interactive mode
	UseInteractive bool
	SaveArgs       bool
	ReadSavedArgs  bool
}

// IsEmpty checks current arguments has default state without any data
func (args ArgumentOptions) IsEmpty() bool {
	empty := ArgumentOptions{}
	return reflect.DeepEqual(args, empty)
}

// Init initializes argument flags
func (args *ArgumentOptions) Init() {
	flag.StringVar(&args.ProjectName, "prjName", "", "Name of the project")

	// connection flags
	flag.StringVar(&args.ConsoleServiceName, "c", "", "Name of console monolitic service")
	flag.StringVar(&args.WorkingDirectory, "d", "", "Path to the working directory")
	flag.StringVar(&args.NetPipeServiceName, "n", "", "Name of netpipe service")
	flag.StringVar(&args.User, "u", "", "User name with admin permisions")
	flag.StringVar(&args.Password, "p", "", "Password for user")

	// mailing flags
	flag.Var(&args.ToEmailList, "t", "List of emails which will send message")
	flag.StringVar(&args.Body, "b", "See the attached log file for details", "Message body")
	flag.StringVar(&args.From, "f", "sergey.zalunin@akforta.com", "From email")
	flag.StringVar(&args.SMTPServer, "smtp", "mail.akforta.com", "address of SMTP server")
	flag.IntVar(&args.SMTPPort, "port", 465, "Port of SMTP server")
	flag.StringVar(&args.SMTPLogin, "smtplogin", "", "Login to SMTP server")
	flag.StringVar(&args.SMTPPassword, "smtppass", "", "Password to SMTP server")

	// database flags
	flag.StringVar(&args.DbDataSource, "dbdatasource", "localhost", "Database data source name")
	flag.StringVar(&args.DatabaseName, "dbname", "", "Database Name")
	flag.StringVar(&args.DatabaseUserID, "dbuserid", "", "Database Login")
	flag.StringVar(&args.DatabasePassword, "dbpassword", "", "Database Password")
	flag.StringVar(&args.BackupPath, "backuppath", "", "Path to store backups of database")
	flag.BoolVar(&args.UseCompression, "usecompr", false, "Use compression on backup database")
	flag.BoolVar(&args.TrustedConnection, "dbtrust", false, "Database allows trusted connection")
	flag.IntVar(&args.ConnectionTimeout, "dbtimeout", 7200, "Connection timeout to mssql")
	flag.BoolVar(&args.SkipBackup, "skipbackup", false, "SkipBackup added to skip backup of the second installation. "+
		"E.g the first installation was failed. "+
		"So it is not necessary to make backup from brocken database")

	// interactive mode
	flag.BoolVar(&args.UseInteractive, "interactive", false,
		"Call the process to set up arguments settings in the console. "+
			"All entered values will store in the file which reads on starting this program")
	flag.BoolVar(&args.SaveArgs, "saveargs", false,
		"Saving entered arguments in the file which reads on starting this programm")
	flag.BoolVar(&args.ReadSavedArgs, "rsd", false,
		"Reading saving arguments from the data.dat file")

	flag.Parse()
}
