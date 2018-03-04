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

// Ruler ... Perform advanced screen I/O :>
func Ruler(s *State) {
	QuietPrintf(s, "=====================================================")
}

// Banner ... Print the Gobuster banner to the screen
func Banner(s *State) {
	QuietPrintf(s, "\n")
	QuietPrintf(s, "Gobuster v1.4.1              OJ Reeves (@TheColonial)")
	Ruler(s)
}


func QuietPrintf(s *State, format string, args ...interface{}) {
    if !s.JSON && !s.Quiet {
    	msg := fmt.Sprintf(format, args...)
    	fmt.Print(msg)
    }
}

// ShowConfig ... Print the state to the screen
func ShowConfig(s *State) {
	//var fields *logrus.Fields = {}
	var fields = make(logrus.Fields)

	if s.Quiet {
		return
	}

	if s != nil {
		QuietPrintf(s, "[+] Mode         : %s\n", s.Mode)
		fields["Mode"] = s.Mode
		QuietPrintf(s, "[+] URL/Domain   : %s\n", s.URL)
		fields["URL"] = s.URL
		QuietPrintf(s, "[+] Threads      : %d\n", s.Threads)
		fields["Threads"] = s.Threads

		wordlist := "stdin (pipe)"
		if !s.StdIn {
			wordlist = s.Wordlist
		}

		QuietPrintf(s, "[+] Wordlist     : %s\n", wordlist)
		fields["Output file"] = s.OutputFileName

		if s.OutputFileName != "" {
			QuietPrintf(s, "[+] Output file  : %s\n", s.OutputFileName)
			fields["Output file"] = s.OutputFileName
		}

		if s.Mode == "dir" {
			QuietPrintf(s, "[+] Status codes : %s\n", s.StatusCodes.Stringify())
			fields["Status codes"] = s.StatusCodes.Stringify()
			
			if s.ProxyURL != nil {
				QuietPrintf(s, "[+] Proxy        : %s\n", s.ProxyURL)
				fields["Proxy"] = s.ProxyURL
			}

			if s.Cookies != "" {
				QuietPrintf(s, "[+] Cookies      : %s\n", s.Cookies)
				fields["Cookies"] = s.Cookies
			}

			if s.UserAgent != "" {
				QuietPrintf(s, "[+] User Agent   : %s\n", s.UserAgent)
				fields["User Agent"] = s.UserAgent
			}

			if s.IncludeLength {
				QuietPrintf(s, "[+] Show length  : true\n")
				fields["Show length"] = "true"
			}

			if s.Username != "" {
				QuietPrintf(s, "[+] Auth User    : %s\n", s.Username)
				fields["Auth User"] = s.Username
			}

			if len(s.Extensions) > 0 {
				QuietPrintf(s, "[+] Extensions   : %s\n", strings.Join(s.Extensions, ","))
				fields["Extensions"] = strings.Join(s.Extensions, ",")
			}

			if s.UseSlash {
				QuietPrintf(s, "[+] Add Slash    : true\n")
				fields["Add Slash"] = "true"
			}

			if s.FollowRedirect {
				QuietPrintf(s, "[+] Follow Redir : true\n")
				fields["Follow Redir"] = "true"
			}

			if s.Expanded {
				QuietPrintf(s, "[+] Expanded     : true\n")
				fields["Expanded"] = "true"
			}

			if s.NoStatus {
				QuietPrintf(s, "[+] No status    : true\n")
				fields["No status"] = "true"
			}

			if s.Verbose {
				QuietPrintf(s, "[+] Verbose      : true\n")
				fields["Verbose"] = "true"
			}
		}
		fields["Status"] = "starting"

		if s.JSON {
			s.Logger.WithFields(fields).Info("Scan initializing")
		}
		Ruler(s)
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
