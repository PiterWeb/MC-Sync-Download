package main

import (
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"os"
	"strings"
	"github.com/pterm/pterm"
)

type File struct {
	Name       string
	Data       []byte
	FolderName string
}

func cloudMethod() (File, error) {

	credentialsFile, err := os.ReadFile("./serviceAccountKey.json")

	if err != nil {
		return File{}, err
	}

	storageBucket, err := getEnv("ST_BUCKET")

	var config = &firebase.Config{
		StorageBucket: storageBucket,
	}
	var opt = option.WithCredentialsJSON(credentialsFile)
	firebaseApp, err := firebase.NewApp(ctx, config, opt)

	if err != nil {
		return File{}, err
	}

	client, err := firebaseApp.Storage(ctx)
	if err != nil {
		return File{}, err
	}

	worlds, err := getWorldList(client)

	if err != nil {
		return File{}, err
	}

	pterm.DefaultSection.Println("📚 Listado de archivos")
	pterm.Println()

	var fileName string

	fileName, _ = pterm.DefaultInteractiveSelect.WithDefaultText("-> Seleccione el mundo").
	WithOptions(worlds).
	Show()

	// fmt.Println("------ 🏷 Ingrese el nombre del archivo a descargar: ------")

	// fmt.Printf("Introduce comillas dobles para seleccionar el archivo: %c %c\n", rune(34), rune(34))

	// fmt.Scanf("%q", &fileName)

	fileData, err := downloadWorld(client, fileName)

	folderName := strings.Split(fileName, "-")[0]

	if _, err = os.Stat(folderName); !os.IsNotExist(err) {
		pterm.Info.Println("El directorio ya existe, este se sobreescribirá") 
		os.RemoveAll(folderName)
	}

	file := File{
		Name:       fileName,
		Data:       fileData,
		FolderName: folderName,
	}

	return file, nil

}

func directMethod() (File, error) {

	var url string

	url, err := pterm.DefaultInteractiveTextInput.WithDefaultText("📚 Introduzca la URL").WithTextStyle(pterm.NewStyle(pterm.FgLightYellow)).Show()

	if err != nil {
		return File{}, err
	}

	file, err := downloadWorldFromUrl(url)

	if err != nil {
		return File{}, err
	}

	pterm.Println()
	pterm.Success.Println("Mundo " + file.Name + " descargandose correctamente")

	if _, err = os.Stat(file.FolderName); !os.IsNotExist(err) {
		pterm.Info.Println("El directorio ya existe, este se sobreescribirá") 
		os.RemoveAll(file.FolderName)
	}

	return file, nil

}
