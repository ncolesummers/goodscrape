/*
Copyright © 2021 N Cole Summers <nsummers72@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var (
	searchURL     string         = "https://www.goodreads.com/quotes/search?q=%s&commit=Search"
	contentRegexp *regexp.Regexp = regexp.MustCompile("“(.+?)”")
	quotesCmd                    = &cobra.Command{
		Use:   "quotes",
		Short: "Search for quotes",
		Long: `Query the quotes database of goodreads
		Example: goodscrape quotes <query>`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var quotes []Quote

			amount, err := cmd.Flags().GetInt("amount")
			if err != nil {
				log.Fatalln("Error parsing flag", err)
			}

			filename, err := cmd.Flags().GetString("filename")
			if err != nil {
				log.Fatalln("Error parsing flag", err)
			}

			query := strings.Join(args, "+")
			c := colly.NewCollector(
				colly.AllowedDomains("www.goodreads.com"),
			)

			// scrape all quotes on a page
			c.OnHTML(".quoteDetails", func(e *colly.HTMLElement) {
				results := contentRegexp.FindAllStringSubmatch(e.ChildText("div.quoteText"), -1)

				if len(results) < 1 {
					return
				}

				if len(results[0]) < 1 {
					return
				}

				quotes = append(quotes, Quote{
					Content: results[0][0],
					Author:  e.ChildText(".authorOrTitle"),
				})

				fmt.Print(".")
			})
			// continue to next page if we need more results
			c.OnHTML(".next_page", func(e *colly.HTMLElement) {
				if len(quotes) < amount {
					e.Request.Visit(e.Attr("href"))
				}
			})

			fmt.Println("Scrape in progress...")

			c.Visit(fmt.Sprintf(searchURL, query))

			fmt.Printf("Scraped %d quotes.\n\n", len(quotes))

			toWrite, err := json.MarshalIndent(quotes, "", "  ")
			if err != nil {
				log.Fatalln("Marshall error:", err.Error())
			}

			err = ioutil.WriteFile(filename, toWrite, 0644)
			if err != nil {
				log.Fatalln("File write error: ", err.Error())

			}
		},
	}
)

type Quote struct {
	Author  string `json: "author"`
	Content string `json: "content"`
}

func (q *Quote) String() string {
	return fmt.Sprintf("%s ― %s", q.Content, q.Author)
}

func init() {
	rootCmd.AddCommand(quotesCmd)

	// Here you will define your flags and configuration settings.

	// quotesCmd.PersistentFlags().StringP("tag", "t", "", "Display the top quotes by tag")
}
