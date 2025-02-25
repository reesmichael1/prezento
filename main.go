package main

import (
	// "bufio"
	// "fmt"
	// "os"

	// tea "github.com/charmbracelet/bubbletea"
	// "github.com/reesmichael1/prezento/tui"
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/reesmichael1/prezento/slide"
	"github.com/reesmichael1/prezento/tui"
)

func main() {
	// runCommand("echo", []string{"'Hi'"})
	// time.Sleep(time.Second * 5)

	file, err := os.Open("./presentation.md")
	if err != nil {
		panic(err)
	}

	slides := slide.Slides{}
	currentSlide := ""
	delimiter := "---"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == delimiter {
			// slide := slide.Slide{Content: currentSlide}
			slide := slide.New(currentSlide)
			slides = append(slides, slide)
			currentSlide = ""
			continue
		}
		currentSlide += line + "\n"
	}
	if currentSlide != "" {
		slides = append(slides, slide.Slide{Content: currentSlide, Kind: slide.ContentSlide})
	}

	// slidesContent := []string{}
	// for _, slide := range slides {
	// 	slidesContent = append(slidesContent, slide.Content)
	// }
	//
	model := tui.NewPages(slides)

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error running program: %v\n", err)
		os.Exit(1)
	}
}
