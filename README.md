# TView Chooser

TView Chooser is a simple directory and file chooser GUI GOlang library built with [tview](https://github.com/rivo/tview).

![ScreenShot](https://raw.githubusercontent.com/AEROGU/tvchooser/main/ScreenShot.png)

## Features
- Navigate through directories
- Select files and directories
- Windows support
- It can show hidden files

## Usage

To add this package to your project:

```
go get github.com/aerogu/tvchooser
```

Get the library and call FileChooser or DirectoryChooser passing your tview `app` as parameter to pause it while the chooser interface is in use, or pass `nil` if you are building a non GUI console application but want the user to select a file or directory at some point.

## Hello World

Simple file chooser, English is the default language.

```go
package main

import (
	"fmt"
	"github.com/aerogu/tvchooser"
)

func main() {
	path := tvchooser.FileChooser(nil)
	fmt.Print("RUTA: " + path)
}
```

Spanish language is included too, but you can set your own language by filling your own `tvclang.Texts` object and set it with `tvclang.SetTranslations(yourTranslations)`

```go
package main

import (
	"fmt"

	"github.com/aerogu/tvchooser"
	"github.com/aerogu/tvchooser/tvclang"
)

func main() {
	tvclang.SetTranslations(tvclang.LangSpanish())
	path := tvchooser.FileChooser(nil)

	fmt.Print("RUTA: " + path)
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
	}
}
```

For more deep modifications, new features or improvements send a pull request.
