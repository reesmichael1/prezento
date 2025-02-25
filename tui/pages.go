package tui

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/creack/pty"
)

type Pages struct {
	slides []string
	paginator paginator.Model
	ready bool
	viewport viewport.Model
	codeCursor int
	codeBlocks []CodeBlock
	codeBlockToRender string
	commandOutput string
	altScreen bool
}

type CodeBlock struct {
	commands []string
	language string
	code     string
	inline   bool
}

func (p Pages) ExecuteCommand(command string, args []string) (string, error) {
	var output bytes.Buffer

	cmd := exec.Command(command, args...)

	f, err := pty.Start(cmd)

	if err != nil {
		cmd.Cancel()
		return "", err
	}

	// read the output from the pty
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		output.Write(append(scanner.Bytes(), '\n'))
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading command output: %w", err)
	}

	f.Close()

	p.Update(tea.Msg(output.String()))

	return output.String(), nil
}

func NewPages(slides []string) Pages {
	pages := paginator.New()
	pages.Type = paginator.Dots
	pages.PerPage = 1
	pages.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	pages.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	pages.SetTotalPages(len(slides))

	return Pages{
		slides: slides,
		paginator: pages,
		ready: false,
		viewport: viewport.New(0, 0),
		codeCursor: 0,
		codeBlocks: []CodeBlock{},
		codeBlockToRender: "",
		commandOutput: "",
	}
}

func (p Pages) Init() tea.Cmd {
	return nil
}

func (p Pages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		cmds []tea.Cmd
	)

	oldPage := p.paginator.Page

	p.paginator, cmd = p.paginator.Update(msg)
	cmds = append(cmds, cmd)

	if oldPage != p.paginator.Page {
		p.codeBlocks = []CodeBlock{}
		p.codeBlockToRender = ""
		p.codeCursor = 0
		p.commandOutput = ""
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		
		case "enter":
			if len(p.codeBlocks) > 0 && p.codeCursor < len(p.codeBlocks) {
				commandToExecute := p.codeBlocks[p.codeCursor].commands[0]
				commandArgs := p.codeBlocks[p.codeCursor].commands[1:]

				p.altScreen = !p.altScreen

				if p.altScreen {
					cmd = tea.ExitAltScreen
				} else {
					cmd = tea.EnterAltScreen
				}

				output, err := p.ExecuteCommand(commandToExecute, commandArgs)

				if err != nil {
					panic(err)
				}

				p.codeCursor++

				p.commandOutput = output

				return p, cmd
			}
		}

	case tea.WindowSizeMsg:
		// headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(p.paginator.View())

		// verticalMarginHeight := headerHeight + footerHeight
		if !p.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			p.viewport = viewport.New(msg.Width, msg.Height-footerHeight)
			// m.viewport.YPosition = 0
			p.ready = true
		} else {
			p.viewport.Width = msg.Width
			p.viewport.Height = msg.Height - footerHeight
		}
	}

	// find any code blocks in the slide that use the run:language tag syntax
	// so we can pull the code out and execute it using Bubbletea's `cmd.ExecProcess()`
	var command []string = []string{}
	var language string = ""
	var code []string = []string{}

	lines := strings.Split(p.slides[p.paginator.Page], "\n")

	inCodeBlock := false
	runPrefix := "$>"

	for _, line := range lines {
		if !inCodeBlock && (strings.HasPrefix(line, "```"+runPrefix) || strings.Contains(line, "`"+runPrefix)) {
			inCodeBlock = true
			code = []string{}
			command = []string{}

			// if we're in a code block, set the language and remove the $> from the line
			// example:
			// ```$>shell
			// echo "Pinging 8.8.8.8..."
			// ping 8.8.8.8
			// ``` -> {
			// 	language = "shell"
			// 	code = ["```shell\necho \"Pinging 8.8.8.8...\"\n```"] "what we'll display"
			// 	command = ["ping", "8.8.8.8"] "what we'll execute"
			// }
			if strings.HasPrefix(line, "```$>") {
				// FENCED CODE BLOCK
				language = strings.Replace(line, "```"+runPrefix, "", 1)

				code = append(code, strings.Replace(line, "```"+runPrefix, "```", 1))

				command = []string{}
			} else {
				// INLINE CODE BLOCK
				// find the first instance of the runPrefix in the line
				runPrefixIndex := strings.Index(line, runPrefix)
				// get the next instance of the backtick in the line after the runPrefixIndex
				closingBackTickIndex := strings.Index(line[runPrefixIndex:], "`") + runPrefixIndex

				if closingBackTickIndex == -1 {
					closingBackTickIndex = len(line)
				}

				// get the code from runPrefixIndes to the closing "`"
				inlineCode := line[runPrefixIndex:closingBackTickIndex]

				// get the actual command from the inline code block, example: `$>ping 8.8.8.8` -> "ping 8.8.8.8"
				// remove the runPrefix from the inline code
				inlineCode = strings.Replace(inlineCode, runPrefix, "", 1)
				// remove the backticks from the inline code
				inlineCode = strings.Replace(inlineCode, "`", "", 1)

				if inlineCode != "" {

					command = append(command, inlineCode)

					p.codeBlocks = append(p.codeBlocks, CodeBlock{
						commands: strings.Split(inlineCode, " "),
						language: "shell",
						code: "`"+inlineCode+"`",
						inline: true,
					})
				}

				// we're done with the inline code block, so set inCodeBlock to false
				inCodeBlock = false
			}
		} else if inCodeBlock {
			// if we're at the end of a code block
			if strings.HasSuffix(line, "```") {
				code = append(code, line)

				// add the command to the codeBlocks slice
				p.codeBlocks = append(p.codeBlocks, CodeBlock{
					commands: command,
					language: language,
					code:     strings.Join(code, "\n"),
					inline:   false,
				})

				// reset the command and language variables
				command = []string{}

				// reset the inCodeBlock flag
				inCodeBlock = false
			} else {
				command = append(command, strings.Split(line, " ")...)
				code = append(code, line)
			}
		}
	}

	var newSlide string = p.slides[p.paginator.Page]

	for _, codeBlock := range p.codeBlocks {
		var commandString string

		if codeBlock.inline {
			commandString = strings.TrimSpace("`" + runPrefix + strings.Join(codeBlock.commands, " ")) + "`"
		} else {
			commandString = strings.TrimSpace("```" + runPrefix + codeBlock.language + "\n" + strings.Join(codeBlock.commands, "\n")) + "\n```"
		}

		// TODO: figure out why this doesn't work and the command blocks aren't replaced with code blocks
		newSlide = strings.Replace(newSlide, commandString, codeBlock.code, 1)
	}

	// TODO: figure out how to use glamour.WithAutoStyle() for this
	rendered, err := glamour.Render(newSlide, "dark")

	if err != nil {
		panic(err)
	}

	for i, codeBlock := range p.codeBlocks {
		if i == p.codeCursor {
			if codeBlock.inline {
				p.codeBlockToRender = codeBlock.code
			} else {
				p.codeBlockToRender = strings.Join(codeBlock.commands, " ")
			}
		}
	}

	if p.codeBlockToRender != "" {
		commandBlockPrefix := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0000")).Render("> ")

		var renderedCommandBlock string

		if p.commandOutput != "" {
			renderedCommandBlock = lipgloss.NewStyle().Foreground(lipgloss.Color("#9cc6ab")).Render(p.codeBlockToRender)
			renderedCommandBlock = lipgloss.NewStyle().MarginLeft(5).MarginTop(1).Render(commandBlockPrefix + renderedCommandBlock)
			
			renderedOutput := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ff00")).Render(p.commandOutput)
			renderedOutput = lipgloss.NewStyle().MarginLeft(5).MarginTop(1).Render(commandBlockPrefix + renderedOutput)

			renderedCommandBlock = renderedCommandBlock + "\n" + renderedOutput

			p.viewport.SetContent(rendered + "\n" + renderedCommandBlock)
		} else {
			renderedCommandBlock = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ff00")).Render(p.codeBlockToRender)

			renderedCommandBlock = lipgloss.NewStyle().MarginLeft(5).MarginTop(1).Render(commandBlockPrefix + renderedCommandBlock)

			p.viewport.SetContent(rendered + "\n" + renderedCommandBlock)
		}
	} else {
		p.viewport.SetContent(rendered)
	}

	p.viewport, cmd = p.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p Pages) View() string {
	if !p.ready {
		return lipgloss.NewStyle().Width(p.viewport.Width).Height(p.viewport.Height).Align(lipgloss.Center).Render("\n  Initializing...")
	}

	var b strings.Builder

	b.WriteString(" " + lipgloss.NewStyle().Width(p.viewport.Width).Height(p.viewport.Height).Align(lipgloss.Center).Render(p.viewport.View()))
	b.WriteString("\n\n" + lipgloss.NewStyle().Width(p.viewport.Width).Align(lipgloss.Center).Render(p.paginator.View()))

	return b.String()
}
