package main

//----------------------------------------------------
// Gobuster -- by OJ Reeves
//
// A crap attempt at building something that resembles
// dirbuster or dirb using Go. The goal was to build
// a tool that would help learn Go and to actually do
// something useful. The idea of having this compile
// to native code is also appealing.
//
// Run: gobuster -h
//
// Please see THANKS file for contributors.
// Please see LICENSE file for license details.
//
//----------------------------------------------------

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	log "github.com/sirupsen/logrus"
	"./libgobuster"
	//"github.com/mzpqnxow/gobuster/libgobuster"
)

// ParseCmdLine ... Parse all the command line options into a settings
// instance for future use.
func ParseCmdLine() *libgobuster.State {
	var extensions string
	var codes string
	var proxy string

	s := libgobuster.InitState()

	// Set up the variables we're interested in parsing.
	flag.IntVar(&s.Threads, "t", 10, "Number of concurrent threads")
	flag.StringVar(&s.Mode, "m", "dir", "Directory/File mode (dir) or DNS mode (dns)")
	flag.StringVar(&s.Wordlist, "w", "", "Path to the wordlist")
	flag.StringVar(&codes, "s", "200,204,301,302,307", "Positive status codes (dir mode only)")
	flag.StringVar(&s.OutputFileName, "o", "", "Output file to write results to (defaults to stdout)")
	flag.StringVar(&s.URL, "u", "", "The target URL or Domain")
	flag.StringVar(&s.Cookies, "c", "", "Cookies to use for the requests (dir mode only)")
	flag.StringVar(&s.Username, "U", "", "Username for Basic Auth (dir mode only)")
	flag.StringVar(&s.Password, "P", "", "Password for Basic Auth (dir mode only)")
	flag.StringVar(&extensions, "x", "", "File extension(s) to search for (dir mode only)")
	flag.StringVar(&s.UserAgent, "a", "", "Set the User-Agent string (dir mode only)")
	flag.StringVar(&proxy, "p", "", "Proxy to use for requests [http(s)://host:port] (dir mode only)")
	flag.StringVar(&s.ContentType, "ct", "", "Default Content-Type for POST requests")
	flag.StringVar(&s.Verb, "X", "GET", "Verb to use instead of GET (GET, POST, PUT) are valid)")
	flag.StringVar(&s.Body, "b", "", "Content of POST body, i.e. '{}' for Application/JSON")
	flag.StringVar(&s.Headers, "H", "", "List of arbitrary headers to supply, separated by '|' characters")
	flag.BoolVar(&s.JSON, "J", false, "Use JSON formatting for output to stdout and output file")
	flag.BoolVar(&s.Verbose, "v", false, "Verbose output (errors)")
	flag.BoolVar(&s.ShowIPs, "i", false, "Show IP addresses (dns mode only)")
	flag.BoolVar(&s.ShowCNAME, "cn", false, "Show CNAME records (dns mode only, cannot be used with '-i' option)")
	flag.BoolVar(&s.FollowRedirect, "r", false, "Follow redirects")
	flag.BoolVar(&s.Quiet, "q", false, "Don't print the banner and other noise")
	flag.BoolVar(&s.Expanded, "e", false, "Expanded mode, print full URLs")
	flag.BoolVar(&s.NoStatus, "n", false, "Don't print status codes")
	flag.BoolVar(&s.IncludeLength, "l", false, "Include the length of the body in the output (dir mode only)")
	flag.BoolVar(&s.FullLogging, "F", false, "Print timestamp and loglevel on every non-JSON stdout output line")
	flag.BoolVar(&s.UseSlash, "f", false, "Append a forward-slash to each directory request (dir mode only)")
	flag.BoolVar(&s.WildcardForced, "fw", false, "Force continued operation when wildcard found")
	flag.BoolVar(&s.InsecureSSL, "k", false, "Skip SSL certificate verification")

	flag.Parse()

	libgobuster.Banner(&s)

	if err := libgobuster.ValidateState(&s, extensions, codes, proxy); err.ErrorOrNil() != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}

	return &s
}

type (
	Formatter struct {
		Name      string
		Formatter log.Formatter
	}

	HumanFormatter struct{
		KVSeparator string
		MessageWrapperString [2]string
	}
)



/*func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	f.Do(func() { f.init(entry) })

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		}
		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, "message", entry.Message)
		}
		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

*/

func (f *HumanFormatter) Format(entry *log.Entry) ([]byte, error) {
	keys := make([]string, 0, len(entry.Data))
	b := &bytes.Buffer{}
	kvSeparator := " "
	messageWrapperString := [2]string{"(", ")"}
	if (messageWrapperString != f.MessageWrapperString) {
		messageWrapperString = f.MessageWrapperString
	}
	// if (kvSeparator != f.KVSeparator) {
	//	kvSeparator = f.KVSeparator
	// }
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//fmt.Fprintf(b, "[%s]", entry.Time.String())
	//fmt.Fprintf(b, " [%s]", strings.ToUpper(entry.Level.String()))
	//fmt.Fprintf(b, " %s", entry.Message)
	fmt.Fprintf(b, "%s%s%s", messageWrapperString[0], entry.Message, messageWrapperString[1])

	for _, key := range keys {
		if key != "time" && key != "level" && key != "message" {
			fmt.Fprintf(b, "%s%s=%s", kvSeparator, key, entry.Data[key])
		}
	}
	fmt.Fprintf(b, "\n")
	
	return b.Bytes(), nil
}

func setStdoutLogging(s *libgobuster.State) {
	s.Logger = log.New()
	s.Logger.Out = os.Stdout
	if s.JSON {
		s.Logger.Formatter = (&log.JSONFormatter{ TimestampFormat: "2006-01-02 15:04:05"})
	} else if !s.FullLogging {
		s.Logger.Formatter = &HumanFormatter {}
	} else {
		s.StdoutLogger.Formatter = (&log.TextFormatter{})
	}
  	s.Logger.SetLevel(log.InfoLevel)
}

func setLogging(s *libgobuster.State) {
	setStdoutLogging(s)
}

func main() {
	state := ParseCmdLine()
	if state != nil {
		setLogging(state)
		libgobuster.Process(state)
	}
}
