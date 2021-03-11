// При каждой записи в канал проверяется не прервана ли горутина,
// которая этот канал слушает.
// Тем самым, если этот канал никто не слушает,
// то нельзя будет записать что-то в канал и заблокировать выполнение.
package main

import (
	"flag"
	"log"

	"github.com/pryg/safe_channel/parser"
)

func main() {
	var (
		Name string
		ChannelType string
		ImportType string
	)
	flag.StringVar(&Name, "name", "", "name of generated channel")
	flag.StringVar(&ChannelType, "type", "", "type for generate channel")
	flag.StringVar(&ImportType, "import", "", "add package to import")
	flag.Parse()

	prs := parser.New(Name, ChannelType, ImportType)
	if err := prs.Execute(); err != nil {
		log.Fatalf("execute error: %s", err)
	}
}
