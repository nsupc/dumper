package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/codeclysm/extract/v3"
	"github.com/urfave/cli/v2"
)

const (
	regionDumpURL string = "https://www.nationstates.net/pages/regions.xml.gz"
	nationDumpURL string = "https://www.nationstates.net/pages/nations.xml.gz"
)

type Args struct {
	userAgent  string
	nations    bool
	nOutDir    string
	regions    bool
	rOutDir    string
	decompress bool
	dryRun     bool
}

func main() {
	var args Args

	app := &cli.App{
		Name:    "dumper",
		Usage:   "a simple utility for downloading NationStates nation and region dumps",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "user-agent",
				Aliases:     []string{"u"},
				Usage:       "NS nation or email address for API identification",
				Required:    true,
				Destination: &args.userAgent,
			},
			&cli.StringFlag{
				Name:        "out-dir-nations",
				Aliases:     []string{"N"},
				Usage:       "output directory for nation dump",
				Required:    false,
				Destination: &args.nOutDir,
			},
			&cli.StringFlag{
				Name:        "out-dir-regions",
				Aliases:     []string{"R"},
				Usage:       "output directory for region dump",
				Required:    false,
				Destination: &args.rOutDir,
			},
			&cli.BoolFlag{
				Name:        "nations",
				Aliases:     []string{"n"},
				Usage:       "download nation dump",
				Required:    false,
				Destination: &args.nations,
			},
			&cli.BoolFlag{
				Name:        "regions",
				Aliases:     []string{"r"},
				Usage:       "download region dump",
				Required:    false,
				Destination: &args.regions,
			},
			&cli.BoolFlag{
				Name:        "decompress",
				Aliases:     []string{"d"},
				Usage:       "decompress the gzip archives to xml files",
				Required:    false,
				Destination: &args.decompress,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"D"},
				Usage:       "Perform a test run without downloading anything; creates blank output files",
				Required:    false,
				Destination: &args.dryRun,
			},
		},
		Action: func(*cli.Context) error {
			if err := downloadDumps(args); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func downloadDumps(args Args) error {
	client := &http.Client{}

	if args.nations {
		log.Println("Downloading nation dump...")

		filename := filepath.Join(args.nOutDir, fmt.Sprintf("nations_%s.xml", time.Now().Format("2006_01_02")))

		if args.dryRun {
			if err := generateBlankOutputFile(filename); err != nil {
				return err
			}
		} else {
			if err := downloadDump(args, nationDumpURL, filename, client); err != nil {
				return err
			}

			if args.regions {
				time.Sleep(5 * time.Second)
			}
		}
	}

	if args.regions {
		log.Println("Downloading region dump...")

		filename := filepath.Join(args.rOutDir, fmt.Sprintf("regions_%s.xml", time.Now().Format("2006_01_02")))

		if args.dryRun {
			if err := generateBlankOutputFile(filename); err != nil {
				return err
			}
		} else {
			if err := downloadDump(args, regionDumpURL, filename, client); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadDump(args Args, url string, filename string, client *http.Client) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", args.userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if args.decompress {
		log.Println("Saving dump to ", filename)
		err = extract.Archive(
			context.Background(),
			resp.Body,
			filename,
			nil,
		)
		if err != nil {
			return err
		}
	} else {
		filename += ".gz"
		log.Println("Saving dump to", filename)

		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateBlankOutputFile(filename string) error {
	log.Println("Creating dry run output file at", filename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
