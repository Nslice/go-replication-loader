package argsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/sergeyzalunin/go-replication-loader/logger"
	"golang.org/x/crypto/ssh/terminal"
)

// StartInteractiveMode gives a user rights to set up arguments settings in the console mode
// It returns passed (UseInteractive == false) or newly created arguments.
func StartInteractiveMode(inArgs *ArgumentOptions, log *logger.Log) *ArgumentOptions {
	if !inArgs.UseInteractive {
		return inArgs
	}

	args := deepCopy(inArgs, log)
	switchOfUnecessaryAttributes(args)

	startMode(args, log)
	Serialize(args, log)

	return args
}

// SaveArguments in the file. Expected that arguments passed via flag package
func SaveArguments(inArgs *ArgumentOptions, log *logger.Log) {
	if !inArgs.SaveArgs {
		return 
	}

	args := deepCopy(inArgs, log)
	switchOfUnecessaryAttributes(args)

	Serialize(args, log)
}

// ReadArguments from the file
func ReadArguments(log *logger.Log) *ArgumentOptions {
	args := Deserialize(log)
	return args
}

func deepCopy(args *ArgumentOptions, log *logger.Log) *ArgumentOptions {
	arr, err := json.Marshal(args)
	if err != nil {
		log.Error(err, "An error occured on creating deep copy of arguments ")
		return args
	}

	var newArgs ArgumentOptions
	err = json.Unmarshal(arr, &newArgs)
	if err != nil {
		log.Error(err, "An error occured on creating deep copy of arguments")
		return args
	}

	return &newArgs
}

func switchOfUnecessaryAttributes(args *ArgumentOptions) {
	args.UseInteractive = false
	args.SaveArgs = false
	args.ReadSavedArgs = false
}

func startMode(args *ArgumentOptions, log *logger.Log) {
	setProjectName(args, log)

	// connection flags
	setConnectionArgs(args, log)

	// mailing flags
	setEMailSettings(args, log)

	// database flags
	setDatabaseSettings(args, log)

	switchOfUnecessaryAttributes(args)
}

func yes(log *logger.Log) bool {
	str := readStringLine(log, "yes")
	result := strings.TrimSpace(strings.ToLower(str)) == "yes" || strings.TrimSpace(str) == ""
	return result
}

func setConnectionArgs(args *ArgumentOptions, log *logger.Log) {
	fmt.Println("\nDo you want to enter eleed settings (default - yes)?")
	if yes(log) {
		setConsoleServiceName(args, log)
		setConsoleWorkingDirectory(args, log)
		setNetPipeServiceName(args, log)
		setUser(args, log)
		setPassword(args, log)
	}
}

func setEMailSettings(args *ArgumentOptions, log *logger.Log) {
	fmt.Println("\nDo you want to enter email settings (default - yes)?")
	if yes(log) {
		setEmailFrom(args, log)
		setToEmailList(args, log)
		setEmailBody(args, log)
		setEmailSMTPServer(args, log)
		setEmailSMTPPort(args, log)
		setEmailSMTPLogin(args, log)
		setEmailSMTPPassword(args, log)
	}
}

func setDatabaseSettings(args *ArgumentOptions, log *logger.Log) {
	fmt.Println("\nDo you want to enter MSSQL database settings (default - yes)?")
	if yes(log) {
		setDbDataSource(args, log)
		setDatabaseName(args, log)
		setDatabaseUserID(args, log)
		setDatabasePassword(args, log)
		setBackupPath(args, log)
		setUseCompression(args, log)
		setTrustedConnection(args, log)
		setConnectionTimeout(args, log)
	}
}

func printStringDefaults(message string, defaultValue string) {
	if defaultValue == "" {
		fmt.Printf("%s: ", message)
	} else {
		fmt.Printf("%s (previous - %s): ", message, defaultValue)
	}
}

func readStringLine(log *logger.Log, defaultValue interface{}) string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	result := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Error(err)
	}

	result = strings.TrimSuffix(result, "\r\n")
	if result == "" {
		result = fmt.Sprint(defaultValue)
	}

	return result
}

func readPassword(log *logger.Log, defaultValue interface{}) string {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		log.Error(err)
	}

	result := strings.TrimSuffix(string(bytePassword), "\r\n")
	if result == "" {
		result = fmt.Sprint(defaultValue)
	}

	return result
}

func setProjectName(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter the Project Name", args.ProjectName)
	args.ProjectName = readStringLine(log, args.ProjectName)
}

func setConsoleServiceName(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter the Console Monolithic Service Name", args.ConsoleServiceName)
	args.ConsoleServiceName = readStringLine(log, args.ConsoleServiceName)
}

func setConsoleWorkingDirectory(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter the Working Directory", args.WorkingDirectory)
	args.WorkingDirectory = readStringLine(log, args.WorkingDirectory)
}

func setNetPipeServiceName(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter the NetPipe Service Name", args.NetPipeServiceName)
	args.NetPipeServiceName = readStringLine(log, args.NetPipeServiceName)
}

func setUser(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter the eLeed User Name", args.User)
	args.User = readStringLine(log, args.User)
}

func setPassword(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter the eLeed Password: ")
	args.Password = readPassword(log, args.Password)
}

// Email settings

func setEmailFrom(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults(`Enter "From" Email address`, args.From)
	args.From = readStringLine(log, args.From)
}

func setToEmailList(args *ArgumentOptions, log *logger.Log) {
	defaultValue := ""
	if len(args.ToEmailList) > 0 {
		defaultValue = strings.Join(args.ToEmailList, " ")
	}
	printStringDefaults(`Enter "To" Email address(es) divided by a space`, defaultValue)
	line := readStringLine(log, defaultValue)
	args.ToEmailList = strings.Split(line, " ")
}

func setEmailBody(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter Email Body address", args.Body)
	args.Body = readStringLine(log, args.Body)
}

func setEmailSMTPServer(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter SMTP Server", args.SMTPServer)
	args.SMTPServer = readStringLine(log, args.SMTPServer)
}

func setEmailSMTPPort(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter SMTP Port (previous - %d): ", args.SMTPPort)
	line := readStringLine(log, args.SMTPPort)

	if line == "" {
		args.SMTPPort = 0
		return
	}

	port, err := strconv.Atoi(line)
	if err != nil {
		log.Error(err)
	}

	args.SMTPPort = port
}

func setEmailSMTPLogin(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter SMTP Login", args.SMTPLogin)
	args.SMTPLogin = readStringLine(log, args.SMTPLogin)
}

func setEmailSMTPPassword(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter SMTP Password: ")
	args.SMTPPassword = readPassword(log, args.SMTPPassword)
}

// database flags

func setDbDataSource(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter Data Source", args.DbDataSource)
	args.DbDataSource = readStringLine(log, args.DbDataSource)
}

func setDatabaseName(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter Database Name", args.DatabaseName)
	args.DatabaseName = readStringLine(log, args.DatabaseName)
}

func setDatabaseUserID(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter User ID", args.DatabaseUserID)
	args.DatabaseUserID = readStringLine(log, args.DatabaseUserID)
}

func setDatabasePassword(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter Password: ")
	args.DatabasePassword = readPassword(log, args.DatabasePassword)
}

func setBackupPath(args *ArgumentOptions, log *logger.Log) {
	printStringDefaults("Enter Backup Path", args.BackupPath)
	args.BackupPath = readStringLine(log, args.BackupPath)
}

func setUseCompression(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter Use Compression (previous - %t): ", args.UseCompression)
	line := readStringLine(log, args.UseCompression)
	args.UseCompression = strings.TrimSpace(strings.ToLower(line)) == "true"
}

func setTrustedConnection(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter Trusted Connection (previous - %t): ", args.TrustedConnection)
	line := readStringLine(log, args.TrustedConnection)
	args.TrustedConnection = strings.TrimSpace(strings.ToLower(line)) == "true"
}

func setConnectionTimeout(args *ArgumentOptions, log *logger.Log) {
	fmt.Printf("Enter Connection Timeout (previous - %d): ", args.ConnectionTimeout)

	line := readStringLine(log, args.ConnectionTimeout)
	if line == "" {
		args.ConnectionTimeout = 0
		return
	}

	timeout, err := strconv.Atoi(line)
	if err != nil {
		log.Error(err)
	}
	args.ConnectionTimeout = timeout
}