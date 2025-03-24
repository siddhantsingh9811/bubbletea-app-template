package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/savioxavier/termlink"
)

var (
	headingStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AFFF")).MarginTop(1)
	boxStyle            = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#5F5FD7")).Width(48)
	boxDescriptionStyle = lipgloss.NewStyle().Width(48).MarginTop(1)
	boxTechStyle        = lipgloss.NewStyle().Width(48).Foreground(lipgloss.Color("#5F5FD7")).MarginTop(1).Bold(true)
	boxHeadingStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AFFF")).Margin(0)
	mainColour          = lipgloss.NewStyle().Foreground(lipgloss.Color("#00AFFF"))
	width               = lipgloss.NewStyle().Width(44)
	bold                = lipgloss.NewStyle().Bold(true)
	purple              = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A00FF"))
	insta               = lipgloss.NewStyle().Foreground(lipgloss.Color("#E72A80")).Bold(true)
	whatsapp            = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FF475")).Bold(true)
	linkedin            = lipgloss.NewStyle().Foreground(lipgloss.Color("#0CD3FF")).Bold(true)
	// subStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Underline(true).Italic(true).MarginTop(1)
)

func (m *model) getHome() string {
	banner := `
-----------------------------
Welcome to ssh l0calhost.xyz
-----------------------------`
	banner = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, bannerStyle.Render(banner))
	// banner += "\n"

	bunny := `
(\__/) ||
(•ㅅ•) ||
/    づ`
	bunny = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, bunnyStyle.Render(bunny))
	bunny += "\n"

	// text := lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, asciiArt)
	// text += "\n"
	c := lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, helpStyle.Render("Navigation: Arrow Keys + Enter • Quit: Ctrl + C or q"))
	text := banner + bunny + c
	text = lipgloss.PlaceVertical(20, lipgloss.Center, text)
	return text
}

func (m *model) getResume() string {
	t1 := "Hi I'm still trying to figure out a way to render a pdf in a terminal until then you can view it " + termlink.Link("here", "https://l0calhost.xyz/resume.html")

	t1 = lipgloss.PlaceHorizontal(0, lipgloss.Center, t1)
	t1 = lipgloss.PlaceVertical(20, lipgloss.Center, t1)
	return t1
}

type box struct {
	title        string
	description  string
	technologies string
	link         string
}
type projects []box

func (b *box) getStr() string {
	text := boxHeadingStyle.Render(b.title)
	desc := boxDescriptionStyle.Render(b.description)
	text = lipgloss.JoinVertical(0, text, desc)
	tech := fmt.Sprintf("%s %s", boxTechStyle.Render("Techonologies: "), b.technologies)
	text = lipgloss.JoinVertical(0, text, tech)
	link := termlink.ColorLink("Github Repository", b.link, "blue")
	text = lipgloss.JoinVertical(0, text, link)
	text = boxStyle.Render(text)
	return text
}
func (m *model) getProjects() string {
	text := "Here are some of the projects I have worked on.\n"
	text = headingStyle.Render(text)

	p := projects{
		box{
			title:        "DocConnect",
			description:  "A healthcare app for doctors integrating management, collaboration, and AI-driven insights.",
			technologies: "React, NodeJs, Google Cloud, Docker",
			link:         "https://github.com/siddhantsingh9811/doc-connect"},
		box{
			title:        "Zelto",
			description:  "A student-focused scooty rental app enabling easy ride booking and price comparison.",
			technologies: "React, NodeJs, ExpressJs, PWAs, Docker",
			link:         "https://github.com/siddhantsingh9811/zelto-frontend"},
		box{
			title:        "Kat Social Media",
			description:  "A simple social media application inspired by twitter (back when it wasnt.. well whatever it is now)",
			technologies: "React, NodeJs, Strapi, Vercel, Nginx",
			link:         "https://github.com/siddhantsingh9811/kat-social-media"},
		box{
			title:        "Document Classification Model",
			description:  "A CNN-based document classification model developed while at Eisenvault, achieving 92% accuracy.",
			technologies: "Tensorflow, Pandas, Google Colab, FastAPI",
			link:         "https://github.com/siddhantsingh9811/document-classification-model"},
		box{
			title:        "Pomodoro Timer",
			description:  "A pomodoro timer i coded to help with studying instead of actually studying.",
			technologies: "React, NodeJs, PWAs, Vercel",
			link:         "https://github.com/siddhantsingh9811/pomodoro"},
	}

	for _, b := range p {
		text = lipgloss.JoinVertical(0, text, b.getStr())
	}
	return text
}
func (m *model) getAbout() string {
	t1 := "Hi im "
	t1 += bold.Render(mainColour.Render("Siddhant Singh\n"))
	t2 := "I've been programming and competing in hackathons for over 7 years, mostly experienced in Web Development, Devops and AI but I'm always exploring new tech that I find interesting."
	t2 = width.Render(t2)
	t1 = lipgloss.JoinVertical(0, t1, t2)
	t3 := "Made with "
	t3 += bold.Render(purple.Render("Bubble Tea "))
	t3 += "in "
	t3 += bold.Render(mainColour.Render("Go\n"))
	t1 = lipgloss.JoinVertical(lipgloss.Center, t1, t3)
	t1 = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, t1)
	t1 = lipgloss.PlaceVertical(m.viewport.Height, lipgloss.Center, t1)
	return t1
}

func (m *model) getContact() string {
	t1 := bold.Render(mainColour.Render("Contact me:"))
	t2 := helpStyle.Render("You can reach me here\n")
	t1 = lipgloss.JoinVertical(0, t1, t2)

	t2 = "Github: " + termlink.Link("@siddhantsingh9811", "https://github.com/siddhantsingh9811")
	t1 = lipgloss.JoinVertical(0, t1, t2)
	t2 = "Email: " + termlink.Link("ssiddhant9811@gmail.com", "mailto:ssiddhant9811@gmail.com")
	t1 = lipgloss.JoinVertical(0, t1, t2)
	t2 = whatsapp.Render("Whatsapp: ") + termlink.Link("+91 9711990554", "https://api.whatsapp.com/send/?phone=919711990554&text&type=phone_number&app_absent=0")
	t1 = lipgloss.JoinVertical(0, t1, t2)
	t2 = insta.Render("Instagram: ") + termlink.Link("@siddhant_219", "https://www.instagram.com/siddhant_219")
	t1 = lipgloss.JoinVertical(0, t1, t2)
	t2 = linkedin.Render("Linkedin: ") + termlink.Link("@siddhant-singh-3b94371b2", "https://www.linkedin.com/in/siddhant-singh-3b94371b2")
	t1 = lipgloss.JoinVertical(0, t1, t2)

	t1 = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, t1)
	t1 = lipgloss.PlaceVertical(m.viewport.Height, lipgloss.Center, t1)
	return t1
}
