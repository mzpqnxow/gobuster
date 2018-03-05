package libgobuster

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"github.com/sirupsen/logrus"
)

// PrepareSignalHandler ... Set up a SIGINT handler
func PrepareSignalHandler(s *State) {
	s.SignalChan = make(chan os.Signal, 1)
	signal.Notify(s.SignalChan, os.Interrupt)
	go func() {
		for range s.SignalChan {
			// caught CTRL+C
			if !s.Quiet {
				fmt.Println("[!] Keyboard interrupt detected, terminating.")
				s.Terminate = true
			}
		}
	}()
}

// Banner ... Print the Gobuster banner to the screen
func Banner(s *State) {
	RespectfulPrintf(s, "Gobuster v1.4.1              OJ Reeves (@TheColonial)\n")
}


func RespectfulPrintf(s *State, format string, args ...interface{}) {
    if !s.JSON && !s.Quiet {
    	msg := fmt.Sprintf(format, args...)
    	fmt.Print(msg)
    }
}

// ShowConfig ... Print the state to the screen
func ShowConfig(s *State) {
	//var logFields *logrus.Fields = {}
	var logFields = make(logrus.Fields)

	if s.Quiet {
		return
	}

	if s != nil {
		logFields["Mode"] = s.Mode
		logFields["URL/Domain"] = s.URL
		logFields["Threads"] = strconv.Itoa(s.Threads)
		
		wordlist := "stdin (pipe)"
		
		if !s.StdIn {
			wordlist = s.Wordlist
		}
		logFields["Wordlist"] = wordlist

		if s.OutputFileName != "" {
			logFields["Output file"] = s.OutputFileName
		}

		if s.Mode == "dir" {
			logFields["Status Codes"] = s.StatusCodes.Stringify()
			
			if s.ProxyURL != nil {
				logFields["Proxy"] = s.ProxyURL
			}

			if s.Cookies != "" {
				logFields["Cookies"] = s.Cookies
			}

			if s.UserAgent != "" {
				logFields["User Agent"] = s.UserAgent
			}

			if s.IncludeLength {
				logFields["Show length"] = "true"
			}

			if s.Username != "" {
				logFields["Auth User"] = s.Username
			}

			if len(s.Extensions) > 0 {
				logFields["Extensions"] = strings.Join(s.Extensions, ",")
			}

			if s.UseSlash {
				logFields["Add Slash"] = "true"
			}

			if s.FollowRedirect {
				logFields["Follow Redir"] = "true"
			}

			if s.Expanded {
				logFields["Expanded"] = "true"
			}

			if s.NoStatus {
				logFields["No status"] = "true"
			}

			if s.Verbose {
				logFields["Verbose"] = "true"
			}
		}
		logFields["Status"] = "starting"
		//s.Logger.WithFields(logFields).Info("Bust Begin")
		RespectfulPrintf(s, "--- Configuration")
		for key,value := range logFields {
			RespectfulPrintf(s, "\n  * %s=%s", key, value)
		}
		RespectfulPrintf(s, "\n\n")
	}
}

// Add ... Add an element to a set
func (set *StringSet) Add(s string) bool {
	_, found := set.Set[s]
	set.Set[s] = true
	return !found
}

// AddRange ... Add a list of elements to a set
func (set *StringSet) AddRange(ss []string) {
	for _, s := range ss {
		set.Set[s] = true
	}
}

// Contains ... Test if an element is in a set
func (set *StringSet) Contains(s string) bool {
	_, found := set.Set[s]
	return found
}

// ContainsAny ... Check if any of the elements exist
func (set *StringSet) ContainsAny(ss []string) bool {
	for _, s := range ss {
		if set.Set[s] {
			return true
		}
	}
	return false
}

// Stringify ... Stringify the set
func (set *StringSet) Stringify() string {
	values := []string{}
	for s := range set.Set {
		values = append(values, s)
	}
	return strings.Join(values, ",")
}

// Add ... Add an element to a set
func (set *IntSet) Add(i int) bool {
	_, found := set.Set[i]
	set.Set[i] = true
	return !found
}

// Contains ... Test if an element is in a set
func (set *IntSet) Contains(i int) bool {
	_, found := set.Set[i]
	return found
}

// Stringify ... Stringify the set
func (set *IntSet) Stringify() string {
	values := []string{}
	for s := range set.Set {
		values = append(values, strconv.Itoa(s))
	}
	return strings.Join(values, ",")
}
