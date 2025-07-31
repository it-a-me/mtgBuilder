package main

import (
	"cmp"
	"compress/flate"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"mtgBuilder/card"
	"mtgBuilder/fetch"
	"mtgBuilder/query"
)

var name string

func init() {
	name = filepath.Base(os.Args[0])
}

func fetchCmd(args []string) {
	flags := flag.NewFlagSet("fetch", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: Fetch oracle cards.json\n", flags.Name())
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if flags.NArg() != 0 {
		fmt.Fprintf(flags.Output(), "expected no arguments and got %d\n", flags.NArg())
		flags.Usage()
		os.Exit(1)
	}

	bytes, err := fetch.GetOracleCardsJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(bytes))
}

func serializeCmd(args []string) {
	flags := flag.NewFlagSet("serialize", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: Fetch oracle cards.json\n", flags.Name())
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if flags.NArg() != 2 {
		fmt.Fprintf(flags.Output(), "serialize expects arguments in form `%s serialize src.json dst.bin`", name)
		flags.Usage()
		os.Exit(1)
	}

	var cards []card.Card
	content, err := os.ReadFile(flags.Arg(0))
	if err != nil {
		log.Fatalf("failed to read %s: %s", flags.Arg(0), err)
	}

	err = json.Unmarshal(content, &cards)
	if err != nil {
		log.Fatalf("failed to parse json from %s: %s", flags.Arg(0), err)
	}

	out, err := os.OpenFile(flags.Arg(1), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o600)
	if err != nil {
		log.Fatalf("failed to open %s: %s", flags.Arg(1), err)
	}
	defer out.Close()

	w, err := flate.NewWriter(out, flate.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	encoder := gob.NewEncoder(w)
	if err = encoder.Encode(cards); err != nil {
		log.Fatalf("failed to encode cards: %s", err)
	}
}

func benchCmd(args []string) {
	flags := flag.NewFlagSet("benchDecode", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: Time the decoding of all serialized cards\n", flags.Name())
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if flags.NArg() != 1 {
		fmt.Fprintf(flags.Output(), "benchDecode expects arguments in form `%s benchDecode dst.bin`", name)
		flags.Usage()
		os.Exit(1)
	}

	cards, err := decode(flags.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decoded %d cards\n", len(cards))
}

func decode(path string) ([]card.Card, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %s", path, err)
	}
	r := flate.NewReader(f)
	defer r.Close()

	var cards []card.Card
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&cards); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %s", path, err)
	}
	return cards, nil
}

func oracleCmd(args []string) {
	flags := flag.NewFlagSet("search", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: search serialized cards\n", flags.Name())
		flags.PrintDefaults()
	}
	maxArg := flags.Uint("max", 5, "set the max amount of cards to print")
	short := flags.Bool("short", false, "show short output")
	flags.Parse(args)
	printMax := int(*maxArg)
	if flags.NArg() != 2 {
		fmt.Fprintf(flags.Output(), "benchDecode expects arguments in form `%s benchDecode dst.bin 'o:goblin'`", name)
		flags.Usage()
		os.Exit(1)
	}

	q, err := query.Parse(flags.Arg(1), true)
	if err != nil {
		log.Fatalf("invalid regular expression /%s/: %s", flags.Arg(1), err)
	}
	cards, err := decode(flags.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Parsed Queries", "queries", q)
	start := time.Now()
	var matches []int
	for i, c := range cards {
		if q.Matches(&c) {
			matches = append(matches, i)
		}
	}
	slices.SortFunc(matches, func(a, b int) int { return cmp.Compare(cards[a].Name, cards[b].Name) })
	elapsed := time.Since(start)
	log.Printf("searched %d cards in %s", len(cards), elapsed.String())

	fmt.Printf("Showing %d/%d\n", min(printMax, len(matches)), len(matches))
	for i := range min(len(matches), printMax) {
		if *short {
			fmt.Printf("%d.\t%s\n", i, cards[matches[i]].Name)
		} else {
			fmt.Printf("\t%s -- %s\n%s\n\n", cards[matches[i]].Name, cards[matches[i]].OracleID, strings.Join(cards[matches[i]].GetOracleText(), "\n"))
		}
	}
}

func dumpCmd(args []string) {
	flags := flag.NewFlagSet("dump", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: dump card json based on oracle id\n", flags.Name())
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if flags.NArg() != 2 {
		fmt.Fprintf(flags.Output(), "dump expects arguments in form `%s dump cards.bin '75336d52-5278-4e98-9875-b3ac8b91bfe2'`\n", name)
		flags.Usage()
		os.Exit(1)
	}
	cards, err := decode(flags.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	q := query.OracleID{ID: strings.ToLower(flags.Arg(1))}
	for _, c := range cards {
		if q.Matches(&c) {
			j, err := json.MarshalIndent(c, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(j))
			return
		}
	}
}

func serveCmd(args []string) {
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: serve serialized cards over a tcp socket\n", flags.Name())
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if flags.NArg() != 2 {
		fmt.Fprintf(flags.Output(), "serve expects arguments in form `%s serve dst.bin 'localhost:8022'\n", name)
		flags.Usage()
		os.Exit(1)
	}

	cards, err := decode(flags.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	server := Server{cards}
	if err := rpc.DefaultServer.Register(&server); err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", flags.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	rpc.DefaultServer.Accept(listener)
}

type Server struct {
	cards []card.Card
}

func (s *Server) Query(req string, ret *[]card.Card) error {
	// start := time.Now()
	slog.Info("Recieved Request", "query", req, "cards", len(s.cards))

	panic("Unimplemented")
	// re, err := query.NewRegexMatcher(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // q := query.OracleText{Re: re}
	// *ret = (*ret)[:0]
	// for _, c := range s.cards {
	// 	if slices.ContainsFunc(c.GetOracleText(), func(s string) bool { return re.MatchString(strings.ToLower(s)) }) {
	// 		// if q.Matches(&c) {
	// 		*ret = append((*ret), c)
	// 	}
	// }
	// elapsed := time.Since(start)
	// slog.Debug("handled request", "query", req, "matches", len(*ret), "took", elapsed.String())
	// return nil
}

func queryCmd(args []string) {
	flags := flag.NewFlagSet("query", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s: query cards over a tcp socket\n", flags.Name())
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if flags.NArg() != 2 {
		fmt.Fprintf(flags.Output(), "query expects arguments in form `%s query 'localhost:8022' 'goblin'\n", name)
		flags.Usage()
		os.Exit(1)
	}

	conn, err := rpc.Dial("tcp", flags.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var cards []card.Card
	err = conn.Call("Server.Query", flags.Arg(1), &cards)
	if err != nil {
		log.Fatal(err)
	}

	const MAX = 7
	if len(cards) > MAX {
		fmt.Printf("Showing %d/%d\n", MAX, len(cards))
	}
	for i := range min(len(cards), MAX) {
		fmt.Printf("\t%s\n%s\n\n", cards[i].Name, strings.Join(cards[i].GetOracleText(), "\n"))
		// fmt.Printf("%d.\t%s\n", i, cards[matches[i]].Name)
	}
}

var subcommands = []string{"serialize", "benchDecode", "search", "fetch", "serve", "query", "dump"}

func main() {
	log.SetFlags(log.Ltime)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	flags := flag.NewFlagSet(name, flag.ExitOnError)
	recordTime := flags.Bool("time", false, "log the time to execute the command")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s:\n", flags.Name())
		flags.PrintDefaults()
		fmt.Fprintf(flags.Output(), "\nsubcommands:\n  %s\n", strings.Join(subcommands, "\n  "))
		os.Exit(1)
	}
	flags.Parse(os.Args[1:])
	if flags.NArg() < 1 {
		fmt.Fprintf(flags.Output(), "please supply a subcommand, [%s]\n\n", strings.Join(subcommands, ", "))
		flags.Usage()
		os.Exit(1)
	}

	start := time.Now()
	switch strings.ToLower(flags.Arg(0)) {
	case "serialize":
		serializeCmd(flags.Args()[1:])
	case "benchdecode":
		benchCmd(flags.Args()[1:])
	case "search":
		oracleCmd(flags.Args()[1:])
	case "fetch":
		fetchCmd(flags.Args()[1:])
	case "serve":
		serveCmd(flags.Args()[1:])
	case "query":
		queryCmd(flags.Args()[1:])
	case "dump":
		dumpCmd(flags.Args()[1:])
	default:
		fmt.Fprintf(flags.Output(), "unknown command '%s, subcommands: [%s]\n\n", flags.Arg(0), strings.Join(subcommands, ", "))
		flags.Usage()
		os.Exit(1)
	}
	end := time.Since(start)

	if *recordTime {
		log.Printf("finished in %s\n", end.String())
	}
}
