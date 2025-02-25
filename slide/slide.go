package slide

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/creack/pty"
)

type SlideType int

const (
	ContentSlide SlideType = iota
	CodeSlide
)

type Slide struct {
	Content  string
	Kind     SlideType
	Commands []*exec.Cmd
	Cursor   int
}

type Slides []Slide

func New(content string) Slide {
	trimmed := strings.TrimSpace(content)

	if strings.HasPrefix(trimmed, "===") {
		trimmed := strings.Trim(trimmed, "===")

		var commands []*exec.Cmd
		for _, line := range strings.Split(trimmed, "\n") {
			chunks := strings.Split(line, " ")
			commands = append(commands, exec.Command(chunks[0], chunks...))
		}

		return Slide{
			Content:  fmt.Sprintf("```\n%s\n```", trimmed),
			Kind:     CodeSlide,
			Commands: commands,
			Cursor:   0,
		}
	}

	return Slide{
		Content: trimmed,
		Kind:    ContentSlide,
	}
}

func (slide *Slide) Execute() tea.Msg {
	// c := exec.Command(comm, args...)
	if slide.Cursor < len(slide.Commands) {
		return nil
	}
	c := slide.Commands[slide.Cursor]
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	allArgs := strings.Join(c.Args, " ")
	f.Write([]byte("> " + c.Path + " " + allArgs + "\n"))

	// os.Read(f)

	bytes := []byte{}
	if _, err = f.Read(bytes); err != nil {
		panic(err)
	}

	slide.Content = slide.Content + string(bytes)

	// io.Copy(os.Stdout, f)
	slide.Cursor += 1

	return nil
}
