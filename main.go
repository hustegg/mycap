package main

import (
    "errors"
    "fmt"
    "flag"
    "os"
    "strings"

    "github.com/sirupsen/logrus"
)

const (
    Version = "mycap-1.0"
    MyLike = "2006-01-02 15:04:05.999999"
)

var log = logrus.New()
var display = logrus.New()

type keywords []string

func (kws *keywords) String() string {
    return fmt.Sprint(*kws)
}

func (kws *keywords) Set(value string) error {
    if len(*kws) > 0 {
        return errors.New("Keywords already set")
    }

    for _, kw := range strings.Split(value, ",") {
        *kws = append(*kws, kw)
    }
    return nil
}

//var arg_kws keywords
var arg_whiteips keywords
var arg_blackips keywords

var arg_direction *string
var arg_capfile *string
var arg_dev *string
var arg_snaplen *int
var arg_promisc *bool
var arg_myport *int
var arg_pcnt *int64
var arg_detail *bool
var arg_verbose *bool
var arg_jsonfmt *bool
var arg_version *bool

var bpf_port string
var bpf_whiteips string
var bpf_blackips string
var BPF string


func init() {
    flag.CommandLine.SetOutput(os.Stdout)

    flag.Var(&arg_whiteips, "w", "Packets white ip list separated by comma")
    flag.Var(&arg_blackips, "b", "Packets white ip list separated by comma")

    arg_direction = flag.String("d", "client", "Capture MySQL Packet direction [client|server|both]")
    arg_dev = flag.String("i", "eth0", "Network interface name")
    arg_capfile = flag.String("f", "", "Captured packets filename")
    arg_snaplen = flag.Int("s", 65535, "Snap length for pcap packet capture")
    arg_promisc = flag.Bool("m", false, "Capture with promisc mode")
    arg_myport = flag.Int("p", 3306, "MySQL server port capture")
    arg_pcnt = flag.Int64("c", 1024, "Packets number captured before exit")
    arg_detail = flag.Bool("v", false, "Logging in detail")
    arg_verbose = flag.Bool("vv", false, "Logging in verbose")
    arg_jsonfmt = flag.Bool("j", false, "Logging with JSON formatter")
    arg_version = flag.Bool("V", false, "Show version and exit")

    flag.Parse()
    if *arg_version {
        fmt.Println(Version)
        os.Exit(0)
    }

    if *arg_direction == "client" {
        bpf_port = fmt.Sprintf("(dst port %d)", *arg_myport)
    } else if *arg_direction == "server" {
        bpf_port = fmt.Sprintf("(src port %d)", *arg_myport)
    } else {
        bpf_port = fmt.Sprintf("(port %d)", *arg_myport)
    }

    if len(arg_whiteips) > 0 {
        bpf_whiteips = fmt.Sprintf("and (host %s)", strings.Join(arg_whiteips, " or "))
    }
    if len(arg_blackips) > 0 {
        bpf_blackips = fmt.Sprintf("and not (host %s)", strings.Join(arg_blackips, " or "))
    }

    BPF = fmt.Sprintf("tcp and %s %s %s", bpf_port, bpf_whiteips, bpf_blackips)

    if *arg_verbose {
        log.SetLevel(logrus.TraceLevel)
        //log.SetReportCaller(true)
    } else if *arg_detail {
        log.SetLevel(logrus.DebugLevel)
    } else {
        log.SetLevel(logrus.InfoLevel)
    }

    display.SetLevel(logrus.InfoLevel)

    if *arg_jsonfmt {
        log.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: MyLike,
        })
        display.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: MyLike,
        })
    } else {
        log.SetFormatter(&logrus.TextFormatter{
            TimestampFormat: MyLike,
            FullTimestamp: true,
        })
        display.SetFormatter(&logrus.TextFormatter{
            TimestampFormat: MyLike,
            FullTimestamp: true,
        })
    }

    display.SetOutput(os.Stdout)

}

func main() {

    go DisplayPayload(DisplayClient)
    go DisplayPayload(DisplayServer)

    if err := Capture(); err != nil {
        log.Fatal(err)
    }
}

