package cli

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brandenc40/hackernews"
	tm "github.com/buger/goterm"
	"github.com/logrusorgru/aurora"
	cli "github.com/urfave/cli/v2"
)

const (
	top      = "t\n"
	best     = "b\n"
	new      = "n\n"
	quit     = "q\n"
	next     = "n\n"
	previous = "p\n"
)

// HNCliApp -
func HNCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "HackerNews CLI"
	app.Usage = "Get your news on the terminal"
	app.Flags = []cli.Flag{
		&cli.Int64Flag{
			Name:  "limit",
			Value: 10,
			Usage: "Number of stories to be fetched",
		},
		&cli.StringFlag{
			Name:        "type",
			DefaultText: "",
			Usage:       "The type of stories to be fetched",
		},
	}
	app.Action = func(c *cli.Context) error {
		printLogo()
		limit := c.Int64("limit")
		inputType := c.String("type")
		storyType, err := promptForStoryType(inputType)
		if err != nil {
			return err
		}
		pageNum := 1
		for {
			err = displayStories(storyType, int(limit), pageNum)
			if err != nil {
				return err
			}
			actionStr, err := promptForAction()
			if err != nil {
				return err
			}
			val, err := strconv.Atoi(strings.TrimSuffix(actionStr, "\n"))
			if err == nil {
				pageNum = val
				continue
			}
			switch actionStr {
			case quit:
				return nil
			case next:
				pageNum++
				continue
			case previous:
				if pageNum > 1 {
					pageNum--
				}
				continue
			}
		}
	}
	return app
}

func printLogo() {
	tm.Clear()
	tm.Flush()
	tm.MoveCursor(0, 0)
	tm.Print("\n")
	tm.Println(`
	 _    _            _             _   _                   
	| |  | |          | |           | \ | |                  
	| |__| | __ _  ___| | _____ _ __|  \| | _____      _____ 
	|  __  |/ _  |/ __| |/ / _ \ '__| . ' |/ _ \ \ /\ / / __|
	| |  | | (_| | (__|   <  __/ |  | |\  |  __/\ V  V /\__ \
	|_|  |_|\__,_|\___|_|\_\___|_|  |_| \_|\___| \_/\_/ |___/													
  `)
	tm.Print("\n")
	tm.Flush()
}

func printStoryTypePrompt() {
	tm.Print(tm.Color("Select an option:", tm.CYAN))
	tm.Print("\n\n")
	tm.Println("[t] Top Stories")
	tm.Println("[b] Best Stories")
	tm.Println("[n] New Stories")
	tm.Print("\n")
	tm.Flush()
}

func promptForStoryType(inputType string) (hackernews.StoryType, error) {
	if inputType == "" {
		printStoryTypePrompt()
	}
	for {
		var char string
		var err error
		if inputType == "" {
			reader := bufio.NewReader(os.Stdin)
			char, err = reader.ReadString('\n')
			if err != nil {
				return 0, err
			}
		} else {
			char = inputType + "\n"
		}
		tm.Clear()
		tm.Flush()
		tm.MoveCursor(0, 0)
		switch char {
		case top:
			tm.Println(tm.Color("\nTop Stories", tm.WHITE))
			tm.Flush()
			return hackernews.StoriesTop, nil
		case best:
			tm.Println(tm.Color("\nBest Stories", tm.WHITE))
			tm.Flush()
			return hackernews.StoriesBest, nil
		case new:
			tm.Println(tm.Color("\nNew Stories", tm.WHITE))
			tm.Flush()
			return hackernews.StoriesNew, nil
		default:
			tm.Println("Invalid type selection, try again..")
			tm.Print("\n\n")
			tm.Flush()
			inputType = ""
			printStoryTypePrompt()
		}

	}
}

func displayStories(storyType hackernews.StoryType, limit int, pageNumber int) error {
	tm.Clear()
	tm.Flush()
	tm.MoveCursor(0, 0)
	tm.Flush()
	paginatedStories, err := hackernews.GetPaginatedStories(storyType, int(limit), pageNumber)
	if err != nil {
		return err
	}
	switch storyType {
	case hackernews.StoriesTop:
		tm.Println(tm.Color("\nTop Stories", tm.WHITE))
		tm.Flush()
	case hackernews.StoriesNew:
		tm.Println(tm.Color("\nNew Stories", tm.WHITE))
		tm.Flush()
	case hackernews.StoriesBest:
		tm.Println(tm.Color("\nBest Stories", tm.WHITE))
		tm.Flush()
	}
	for _, item := range paginatedStories.Stories {
		unixTime := time.Unix(item.Time, 0)
		tm.Println("---------------------------------------------------------------------------")
		tm.Println(tm.Color(strconv.Itoa(item.Score), tm.GREEN), aurora.Faint(unixTime.Format(time.RFC1123)))
		tm.Println(tm.Bold(item.Title))
		tm.Println(aurora.Faint(tm.Color(item.URL, tm.CYAN)))
		tm.Print("\n")
		tm.Flush()
	}
	tm.Print("\n")
	tm.Printf("Page %d of %d", paginatedStories.PageNumber, paginatedStories.TotalResults/paginatedStories.Limit)
	tm.Print("\n")
	tm.Flush()
	return nil
}

func promptForAction() (string, error) {
	var char string
	var err error
	tm.Println("")
	tm.Println(tm.Color("[n] Next Page   [p] Previous Page   [#] Page Number   [q] Quit", tm.YELLOW))
	tm.Flush()
	reader := bufio.NewReader(os.Stdin)
	char, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return char, nil
}
