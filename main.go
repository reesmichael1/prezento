package main

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/reesmichael1/prezento/tui"
)

type Slide struct {
	content string
}

type Slides []Slide

func main() {
	file, err := os.Open("./presentation.md")
	if err != nil {
		panic(err)
	}

	slides := Slides{}
	currentSlide := ""
	delimiter := "---"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == delimiter {
			slide := Slide{content: currentSlide}
			slides = append(slides, slide)
			currentSlide = ""
			continue
		}
		currentSlide += line + "\n"
	}
	if currentSlide != "" {
		slides = append(slides, Slide{content: currentSlide})
	}

	slidesContent := []string{}
	for _, slide := range slides {
		slidesContent = append(slidesContent, slide.content)
	}

	model := tui.NewPages(slidesContent)

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error running program: %v\n", err)
		os.Exit(1)
	}
}
