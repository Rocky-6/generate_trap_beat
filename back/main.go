package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Rocky-6/trap/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	frontHost := os.Getenv("FRONT_HOST")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontHost},
		AllowMethods: []string{http.MethodGet},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/download", handleDownload)
	e.Logger.Fatal(e.Start(":8080"))
}

func handleDownload(c echo.Context) error {
	key := c.QueryParam("key")

	// MIDI データのスライスを取得
	kickRepository := service.NewKick()
	kickData, err := kickRepository.MakeSMF(c.Request().Context())
	if err != nil {
		return err
	}

	clapRepository := service.NewClap()
	clapData, err := clapRepository.MakeSMF(c.Request().Context())
	if err != nil {
		return err
	}

	hihatRepository := service.NewHihat()
	hihatData, err := hihatRepository.MakeSMF(c.Request().Context())
	if err != nil {
		return err
	}

	dbRepository, err := service.NewSqliteClient("trap.db")
	if err != nil {
		return err
	}
	chordInformation, err := dbRepository.Scan(c.Request().Context())
	if err != nil {
		return err
	}

	bassRepositroy := service.NewBass(key, chordInformation)
	bassData, err := bassRepositroy.MakeSMF(c.Request().Context())
	if err != nil {
		return err
	}

	chordRepository := service.NewChord(key, chordInformation)
	chordData, err := chordRepository.MakeSMF(c.Request().Context())
	if err != nil {
		return err
	}

	melodyRepository := service.NewMelody(key)
	melodyData, err := melodyRepository.MakeSMF(c.Request().Context())
	if err != nil {
		return err
	}

	// 一時的なディレクトリを作成
	tempDir, err := os.MkdirTemp("", "midi-files")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// MIDI データをファイルに保存
	err = saveMIDIFile(tempDir+"/kick.mid", kickData)
	if err != nil {
		return err
	}

	err = saveMIDIFile(tempDir+"/clap.mid", clapData)
	if err != nil {
		return err
	}

	err = saveMIDIFile(tempDir+"/hihat.mid", hihatData)
	if err != nil {
		return err
	}

	err = saveMIDIFile(tempDir+"/bass.mid", bassData)
	if err != nil {
		return err
	}

	err = saveMIDIFile(tempDir+"/chord.mid", chordData)
	if err != nil {
		return err
	}

	err = saveMIDIFile(tempDir+"/melody.mid", melodyData)
	if err != nil {
		return err
	}

	// ZIP ファイルを作成
	zipFile := tempDir + "/midi-files.zip"
	err = createZIP(zipFile, []string{
		tempDir + "/kick.mid",
		tempDir + "/clap.mid",
		tempDir + "/hihat.mid",
		tempDir + "/bass.mid",
		tempDir + "/chord.mid",
		tempDir + "/melody.mid",
	})
	if err != nil {
		return err
	}

	// ZIP ファイルをクライアントに送信
	err = sendZIP(c.Response().Writer, zipFile)
	if err != nil {
		return err
	}

	return nil
}

func saveMIDIFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func createZIP(zipFile string, files []string) error {
	zipWriter, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipWriter.Close()

	archive := zip.NewWriter(zipWriter)
	defer archive.Close()

	for _, file := range files {
		err := addFileToZIP(archive, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToZIP(zipWriter *zip.Writer, file string) error {
	fileToZip, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(file)

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func sendZIP(w http.ResponseWriter, zipFile string) error {
	file, err := os.Open(zipFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "midi-files.zip"))

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}
	return nil
}
