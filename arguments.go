package main

import "flag"

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
	// connection flags
	ConsoleServiceName string
	WorkingDirectory string
	NetPipeServiceName string
	User string
	Password string

	// mailing flags
	ToEmailList stringSlice
	Body string
	From string
	Subject string
	SMTPServer string
	SMTPPort int
	SMTPLogin string
	SMTPPassword string

	// database flags
	DbDataSource string
	DatabaseName string
	BackupPath string
	UseCompression bool
}

// Init initializes argument flags
func (args *ArgumentOptions) Init()  {
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
	flag.StringVar(&args.Subject, "s", "", "Subject of email")	
	flag.StringVar(&args.SMTPServer, "smtp", "akforta.com", "address of SMTP server")	
	flag.IntVar(&args.SMTPPort, "port", 465, "Port of SMTP server")
	flag.StringVar(&args.SMTPLogin, "smtplogin", "", "Login to SMTP server")
	flag.StringVar(&args.SMTPPassword, "smtppass", "", "Password to SMTP server")

	// database flags
	flag.StringVar(&args.DbDataSource, "dbdatasource", "localhost", "Database data source name")
	flag.StringVar(&args.DatabaseName, "dbname", "", "Database Name")
	flag.StringVar(&args.BackupPath, "backuppath", "", "Path to store backups of database")
	flag.BoolVar(&args.UseCompression, "usecompr", false, "Use compression on backup database")

	flag.Parse()
	flag.PrintDefaults()
}