package main

import (
	"cmp"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func serializeCards(path string, cards []card.Card) error {
	j, err := json.Marshal(cards)
	if err != nil {
		return fmt.Errorf("failed to serialize cards: %w", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o0644)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", path, err)
	}
	defer f.Close()

	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		panic(fmt.Errorf("compression level should always be valid: %w", err))
	}

	_, err = w.Write(j)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", path, err)
	}
	return w.Close()
}

func deserializeCards(path string) ([]card.Card, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open to %s: %w", path, err)
	}
	defer f.Close()

	w, err := gzip.NewReader(f)
	if err != nil {
		panic(fmt.Errorf("compression level should always be valid: %w", err))
	}
	defer w.Close()

	j, err := io.ReadAll(w)
	if err != nil {
		return nil, fmt.Errorf("failed to read to %s: %w", path, err)
	}

	var cards []card.Card
	err = json.Unmarshal(j, &cards)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize cards from %s: %w", path, err)
	}
	return cards, nil
}

var subcommands = map[string]struct {
	desc    string
	example string
	run     func(flags *flag.FlagSet, args []string)
}{
	"fetch": {
		"Fetch oracle cards.json",
		"cards.bin",
		func(flags *flag.FlagSet, args []string) {
			flags.Parse(args)
			const NARGS = 1
			if flags.NArg() != NARGS {
				fmt.Fprintf(flags.Output(), "please supply exactly %d arguments", NARGS)
				flags.Usage()
			}

			js, err := fetch.GetOracleCardsJSON()
			if err != nil {
				log.Fatal(err)
			}

			var cards []card.Card
			if err = json.Unmarshal(js, &cards); err != nil {
				log.Fatal(err)
			}

			if len(cards) == 0 {
				panic("serializing 0 cards")
			}

			if err := serializeCards(flags.Arg(0), cards); err != nil {
				log.Fatal(err)
			}
		},
	},
	"benchdecode": {
		"measure the time to decode cards.bin (requires -time)",
		"cards.bin",
		func(flags *flag.FlagSet, args []string) {
			flags.Parse(args)

			const NARGS = 1
			if flags.NArg() != NARGS {
				fmt.Fprintf(flags.Output(), "please supply exactly %d arguments", NARGS)
				flags.Usage()
			}

			cards, err := deserializeCards(flags.Arg(0))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("decoded %d cards\n", len(cards))
		},
	},
	"search": {
		"search serialized cards",
		"-max 20 cards.bin 'oracle:/:.*goblin/'",
		searchCmd,
	},
	"serve": {
		"serve serialized cards over a tcp socket",
		"cards.bin :8099",
		serveCmd,
	},
	"query": {
		"serve serialized cards over a tcp socket",
		"-max 20 :8099 'oracle:/:.*goblin/'",
		queryCmd,
	},
}

func searchCmd(flags *flag.FlagSet, args []string) {
	maxArg := flags.Uint("max", 5, "set the max amount of cards to print")
	short := flags.Bool("short", false, "show short output")
	flags.Parse(args)
	printMax := int(*maxArg)

	const NARGS = 2
	if flags.NArg() != NARGS {
		fmt.Fprintf(flags.Output(), "please supply exactly %d arguments", NARGS)
		flags.Usage()
	}

	queryString := flags.Arg(1)
	cardsPath := flags.Arg(0)

	q, err := query.Parse(queryString, true)
	if err != nil {
		log.Fatalf("failed to parse query '%s': %s", queryString, err)
	}
	slog.Info("Parsed Query", "query", q)

	cards, err := deserializeCards(cardsPath)
	if err != nil {
		log.Fatal(err)
	}

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

	printMax = min(printMax, len(matches))
	fmt.Printf("Showing %d/%d\n", printMax, len(matches))
	for i := range printMax {
		if *short {
			fmt.Printf("%d.\t%s\n", i, cards[matches[i]].Name)
		} else {
			fmt.Printf("\t%s -- %s\n%s\n\n", cards[matches[i]].Name, cards[matches[i]].OracleID, strings.Join(cards[matches[i]].GetOracleText(), "\n"))
		}
	}
}

func serveCmd(flags *flag.FlagSet, args []string) {
	flags.Parse(args)
	const NARGS = 2
	if flags.NArg() != NARGS {
		flags.Usage()
		os.Exit(1)
	}
	cardsPath := flags.Arg(0)
	serveAddr := flags.Arg(1)

	cards, err := deserializeCards(cardsPath)
	if err != nil {
		log.Fatal(err)
	}

	server := Server{cards}
	if err := rpc.DefaultServer.Register(&server); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", serveAddr)
	if err != nil {
		log.Fatal(err)
	}

	rpc.DefaultServer.Accept(listener)
}

type Server struct {
	cards []card.Card
}

type QueryResponse struct {
	Cards []card.Card
	Error error
}

func (s *Server) Query(req string, resp *QueryResponse) error {
	start := time.Now()
	slog.Info("Recieved Request", "query", req, "cards", len(s.cards))
	q, err := query.Parse(req, true)
	if err != nil {
		*resp = QueryResponse{Error: err}
		return nil
	}

	var matched []card.Card
	for _, c := range s.cards {
		if q.Matches(&c) {
			matched = append(matched, c)
		}
	}
	slices.SortFunc(matched, func(a, b card.Card) int { return cmp.Compare(a.Name, b.Name) })
	elapsed := time.Since(start)
	slog.Debug("handled request", "query", req, "matches", len(matched), "took", elapsed.String())

	*resp = QueryResponse{Cards: matched}
	return nil
}

func queryCmd(flags *flag.FlagSet, args []string) {
	maxArg := flags.Uint("max", 5, "set the max amount of cards to print")
	short := flags.Bool("short", false, "show short output")
	flags.Parse(args)
	printMax := int(*maxArg)

	const NARGS = 2
	if flags.NArg() != NARGS {
		fmt.Fprintf(flags.Output(), "please supply exactly %d arguments", NARGS)
		flags.Usage()
		os.Exit(1)
	}

	serverAddr := flags.Arg(0)

	conn, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var resp QueryResponse
	err = conn.Call("Server.Query", flags.Arg(1), &resp)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Error != nil {
		log.Fatal(err)
	}
	cards := resp.Cards

	printMax = min(printMax, len(cards))
	fmt.Printf("Showing %d/%d\n", printMax, len(resp.Cards))
	// slices.SortFunc(matches, func(a, b int) int { return cmp.Compare(cards[a].Name, cards[b].Name) })
	for i := range printMax {
		if *short {
			fmt.Printf("%d.\t%s\n", i, cards[i].Name)
		} else {
			fmt.Printf("\t%s -- %s\n%s\n\n", cards[i].Name, cards[i].OracleID, strings.Join(cards[i].GetOracleText(), "\n"))
		}
	}
}

func main() {
	log.SetFlags(log.Ltime)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	programName := filepath.Base(os.Args[0])

	flags := flag.NewFlagSet(programName, flag.ExitOnError)
	recordTime := flags.Bool("time", false, "log the time to execute the command")

	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s:\n", flags.Name())
		flags.PrintDefaults()
		fmt.Fprintf(flags.Output(), "\nsubcommands:\n")
		for name, subcommand := range subcommands {
			fmt.Fprintf(flags.Output(), "\t%s\t%s\n", name, subcommand.desc)
		}
		os.Exit(1)
	}

	flags.Parse(os.Args[1:])
	if flags.NArg() < 1 {
		fmt.Fprintf(flags.Output(), "please supply a subcommand\n\n")
		flags.Usage()
	}

	start := time.Now()

	name := strings.ToLower(flags.Arg(0))
	if cmd, exists := subcommands[name]; exists {
		cmdFlags := flag.NewFlagSet(name, flag.ExitOnError)
		cmdFlags.Usage = func() {
			fmt.Fprintf(flags.Output(), "%s %s:\n", programName, cmdFlags.Name())
			fmt.Fprintf(flags.Output(), "  ex. %s %s %s\n", programName, cmdFlags.Name(), cmd.example)
			cmdFlags.PrintDefaults()
			os.Exit(1)
		}
		cmd.run(cmdFlags, flags.Args()[1:])
	} else {
		fmt.Fprintf(flags.Output(), "unknown subcommand %s\n", name)
		flags.Usage()
	}

	end := time.Since(start)

	if *recordTime {
		log.Printf("finished in %s\n", end.String())
	}
}
