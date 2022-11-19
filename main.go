package main

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/beito123/nbt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

var users Users

var ctx = context.Background()

func main() {

	accountData, err := os.ReadFile("./accounts.toml")

	if err != nil {
		pterm.Fatal.Println(err)
	}

	_, err = toml.Decode(string(accountData), &users)

	if err != nil {
		pterm.Fatal.Println(err)
	}

	pterm.Println()
	pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("MC-Sync-Download", pterm.NewStyle(pterm.FgCyan))).Render()

	var persona string

	pterm.Println()

	persona, _ = pterm.DefaultInteractiveSelect.WithDefaultText("-> Seleccione su usuario 👤").
		WithOptions(getUsersNames()).
		WithMaxHeight(len(getUsersNames())).
		Show()

	_, err = getUserByName(persona)

	if err != nil {
		pterm.Fatal.Println(err)
		time.Sleep(time.Second * 2)
		os.Exit(1)
	}

	var option string

	var options = []string{
		"Desde nube ☁",
		"Desde url 🔗",
	}

	pterm.Println()
	option, _ = pterm.DefaultInteractiveSelect.WithDefaultText("-> Seleccione un método de descarga ⬇").
		WithOptions(options).
		Show()

	var file File

	if option == options[0] { // Descargar mundo desde nube

		file, err = cloudMethod()

		if err != nil {
			pterm.Fatal.Println(err)
		}

	} else if option == options[1] { // Descargar mundo desde url

		file, err = directMethod()

		if err != nil {
			pterm.Fatal.Println(err)
		}

	} else { // Opción no válida
		pterm.Error.Println("Opción no valida")
		time.Sleep(time.Second * 2)
		os.Exit(1)
	}

	err = os.WriteFile(file.Name, file.Data, os.ModePerm)

	if err != nil {
		pterm.Fatal.Println(err)
	}

	err = unzipSource(file.Name, "")

	if err != nil {
		pterm.Fatal.Println(err)
	}

	err = os.Remove(file.Name)

	if err != nil {
		pterm.Fatal.Println(err)
	}

	for _, user := range users.Users {

		if user.Name == persona {

			streamPlayer, err := nbt.FromFile(path.Join(file.FolderName, "playerdata", user.OfflineUuid+".dat"), nbt.BigEndian)

			if err != nil {
				pterm.Fatal.Println(err)
			}

			playerdata, err := streamPlayer.ReadTag()

			if err != nil {
				pterm.Fatal.Println(err)
			}

			streamWorld, err := nbt.FromFile(path.Join(file.FolderName, "level.dat"), nbt.BigEndian)

			if err != nil {
				pterm.Fatal.Println(err)
			}

			worldData, err := streamWorld.ReadTag()

			if err != nil {
				pterm.Fatal.Println(err)
			}

			tag := &nbt.Compound{
				Value: map[string]nbt.Tag{
					"Data": worldData,
				},
			}

			if err != nil {
				playerdata.Name()
			}

			dataCompound, err := tag.GetCompound("Data")

			if err != nil {
				pterm.Fatal.Println(err)
			}

			dataCompound, err = dataCompound.GetCompound("Data")

			if err != nil {
				pterm.Fatal.Println(err)
			}

			dataCompound.Value["Player"] = playerdata

			finalCompound := &nbt.Compound{
				Value: map[string]nbt.Tag{
					"Data": dataCompound,
				},
			}

			stream := nbt.NewStream(nbt.BigEndian)

			err = stream.WriteTag(finalCompound)

			if err != nil {
				pterm.Fatal.Println(err)
			}

			data, err := nbt.Compress(stream, nbt.CompressGZip, nbt.DefaultCompressionLevel)

			if err != nil {
				pterm.Fatal.Println(err)
			}

			err = os.WriteFile(path.Join(file.FolderName, "level.dat"), data, os.ModePerm)

			if err != nil {
				pterm.Fatal.Println(err)
			}

		}

	}

	pterm.Success.Println("------ ✨ Archivo descargado correctamente ✨ ------")

	time.Sleep(time.Second * 5)

	pterm.Info.Println("------ 🏁 Saliendo del programa 🏁 ------")

	time.Sleep(time.Second * 2)

	os.Exit(0)
}
