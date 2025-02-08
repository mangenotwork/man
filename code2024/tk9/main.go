package main

import (
	"fmt"
	tk "modernc.org/tk9.0"
	"runtime"
)

func main() {
	menubar := tk.Menu()

	fileMenu := menubar.Menu()
	fileMenu.AddCommand(tk.Lbl("New"), tk.Underline(0), tk.Accelerator("Ctrl+N"))
	fileMenu.AddCommand(tk.Lbl("Open..."), tk.Underline(0), tk.Accelerator("Ctrl+O"), tk.Command(func() { tk.GetOpenFile() }))
	tk.Bind(tk.App, "<Control-o>", tk.Command(func() { fileMenu.Invoke(1) }))
	fileMenu.AddCommand(tk.Lbl("Save"), tk.Underline(0), tk.Accelerator("Ctrl+S"))
	fileMenu.AddCommand(tk.Lbl("Save As..."), tk.Underline(5))
	fileMenu.AddCommand(tk.Lbl("Close"), tk.Underline(0), tk.Accelerator("Crtl+W"))
	fileMenu.AddSeparator()
	fileMenu.AddCommand(tk.Lbl("Exit"), tk.Underline(1), tk.Accelerator("Ctrl+Q"), tk.ExitHandler())
	tk.Bind(tk.App, "<Control-q>", tk.Command(func() { fileMenu.Invoke(6) }))
	menubar.AddCascade(tk.Lbl("File"), tk.Underline(0), tk.Mnu(fileMenu))

	editMenu := menubar.Menu()
	editMenu.AddCommand(tk.Lbl("Undo"))
	editMenu.AddSeparator()
	editMenu.AddCommand(tk.Lbl("Cut"))
	editMenu.AddCommand(tk.Lbl("Copy"))
	editMenu.AddCommand(tk.Lbl("Paste"))
	editMenu.AddCommand(tk.Lbl("Delete"))
	editMenu.AddCommand(tk.Lbl("Select All"))
	menubar.AddCascade(tk.Lbl("Edit"), tk.Underline(0), tk.Mnu(editMenu))

	helpMenu := menubar.Menu()
	helpMenu.AddCommand(tk.Lbl("Help Index"))
	helpMenu.AddCommand(tk.Lbl("About..."))
	menubar.AddCascade(tk.Lbl("Help"), tk.Underline(0), tk.Mnu(helpMenu))

	tk.App.WmTitle(fmt.Sprintf("%s on %s", tk.App.WmTitle(""), runtime.GOOS))
	tk.ActivateTheme("azure light")
	tk.App.Configure(tk.Mnu(menubar), tk.Width("8c"), tk.Height("6c")).Wait()
}
