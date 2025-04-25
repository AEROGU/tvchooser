# TView Chooser

TView Chooser is a simple directory and file chooser GUI GOlang library built with [tview](https://github.com/rivo/tview).

![ScreenShot](https://raw.githubusercontent.com/AEROGU/tvchooser/main/ScreenShot.png)

## Features

- Navigate through directories
- Select files and directories
- Windows support (It should also work fine on Linux and other OS)
- It can show hidden files
- You can set custom directories and display them with another name

## Usage

To add this package to your project:

```
go get github.com/AEROGU/tvchooser
```

Get the library and call FileChooser or DirectoryChooser passing your tview `app` as parameter to pause it while the chooser interface is in use, or pass `nil` if you are building a non GUI console application but want the user to select a file or directory at some point, and `true` if you want to show hidden files.

## Hello World

Simple file chooser, English is the default language.

```go
package main

import (
	"fmt"
	"github.com/AEROGU/tvchooser"
)

func main() {
	path := tvchooser.FileChooser(nil, false)
	fmt.Print("RUTA: " + path)
}
```

Spanish language is included too, but you can set your own language by filling your own `tvclang.Texts` object and set it with `tvclang.SetTranslations(yourTranslations)`

```go
package main

import (
	"fmt"

	"github.com/AEROGU/tvchooser"
	"github.com/AEROGU/tvchooser/tvclang"
)

func main() {
	tvclang.SetTranslations(tvclang.LangSpanish())
	path := tvchooser.FileChooser(nil, false)

	fmt.Print("RUTA: " + path)
}
```

### Custom directories

![ScreenShot](https://raw.githubusercontent.com/AEROGU/tvchooser/main/ScreenShot-2.png)

You can add custom directories, and even display them with another name with the syntax "`<path>`|`<name>`" example:

```go
package main

import (
	"fmt"
	"os"

	"github.com/AEROGU/tvchooser"
)

func main() {
	curdir, _ := os.Getwd()

	path := tvchooser.FileChooser(nil, false,
		curdir+"|🏠 Current directory",
		"D:\\artur\\Music|🎵 Music",
		"D:\\artur\\Desktop|🖥️ Desktop",
		"C:\\Users\\Public|📁 |Public access|",
		"C:\\Users",
	)
	if path == "" {
		fmt.Println("No file selected")
		return
	}
	fmt.Println(path)
}
```




### Using it with your application:

```go
package main

import (
	"github.com/AEROGU/tvchooser"
	"github.com/AEROGU/tvchooser/tvclang"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Set your preferred language (If is enlish you can remove this line because it is the default)
	tvclang.SetTranslations(tvclang.LangSpanish())
	// tvclang.SetTranslations(tvclang.LangSpanish())

	app := tview.NewApplication().EnableMouse(true)
	panel := tview.NewFlex().SetDirection(tview.FlexRow)
	text := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	text.SetText("Select a file or directory")

	panel.AddItem(text, 0, 1, false)

	form := tview.NewForm()
	form.AddButton("File", func() {
		path := tvchooser.FileChooser(app, false)
		if path == "" {
			text.SetText("No file selected")
		} else {
			text.SetText(path)
		}
	})
	form.AddButton("Directory", func() {
		path := tvchooser.DirectoryChooser(app, false)
		if path == "" {
			text.SetText("No directory selected")
		} else {
			text.SetText(path)
		}
	})

	panel.SetBackgroundColor(tcell.ColorRed)
	panel.AddItem(form, 4, 0, false)

	app.SetRoot(panel, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

## Dependencies

This package is based on [github.com/rivo/tview](https://github.com/rivo/tview) (and its dependencies) as well as on [github.com/dustin/go-humanize](https://github.com/dustin/go-humanize).

## Your Feedback

Add your issue here on GitHub. Feel free to get in touch if you have any questions.

## Contributing

To request new languages to be added you can either make a pull request or simply request the feature here on GitHub by submitting your own function like this:

```go
func LangSpanish() Texts {
	return Texts{
		Cancel:       "Cancelar",
		Accept:       "Aceptar",
		Modfied:      "Modificado",
		Size:         "Tamaño",
		AccessDenied: "Acceso denegado",
		ThisPC:       "Este PC",
		HomeDir:      "Inicio",
		Devices:      "Unidades",
		Favorites:    "Favoritos",
	}
}
```

For more deep modifications, new features or improvements send a pull request.
