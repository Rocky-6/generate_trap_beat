package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Rocky-6/trap/service"
	"github.com/rs/cors"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET"},
	})

	//http.HandleFunc("/download", handleDownload)
	http.Handle("/download", corsHandler.Handler(http.HandlerFunc(handleDownload)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	key := r.FormValue("key")

	// MIDI データのスライスを取得
	kickRepository := service.NewKick()
	kickData, err := kickRepository.MakeSMF(ctx)

	clapRepository := service.NewClap()
	clapData, err := clapRepository.MakeSMF(ctx)

	hihatRepository := service.NewHihat()
	hihatData, err := hihatRepository.MakeSMF(ctx)

	dbRepository, err := service.NewSqliteClient("trap.db")
	chordInformation, err := dbRepository.Scan(ctx)

	bassRepositroy := service.NewBass(key, chordInformation)
	bassData, err := bassRepositroy.MakeSMF(ctx)

	chordRepository := service.NewChord(key, chordInformation)
	chordData, err := chordRepository.MakeSMF(ctx)

	melodyRepository := service.NewMelody(key)
	melodyData, _ := melodyRepository.MakeSMF(ctx)

	// 一時的なディレクトリを作成
	tempDir, err := os.MkdirTemp("", "midi-files")
	if err != nil {
		log.Println("Failed to create temporary directory:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// MIDI データをファイルに保存
	err = saveMIDIFile(tempDir+"/kick.mid", kickData)
	if err != nil {
		log.Println("Failed to save MIDI file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = saveMIDIFile(tempDir+"/clap.mid", clapData)
	if err != nil {
		log.Println("Failed to save MIDI file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = saveMIDIFile(tempDir+"/hihat.mid", hihatData)
	if err != nil {
		log.Println("Failed to save MIDI file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = saveMIDIFile(tempDir+"/bass.mid", bassData)
	if err != nil {
		log.Println("Failed to save MIDI file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = saveMIDIFile(tempDir+"/chord.mid", chordData)
	if err != nil {
		log.Println("Failed to save MIDI file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = saveMIDIFile(tempDir+"/melody.mid", melodyData)
	if err != nil {
		log.Println("Failed to save MIDI file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
		log.Println("Failed to create ZIP file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ZIP ファイルをクライアントに送信
	sendZIP(w, zipFile)
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

func sendZIP(w http.ResponseWriter, zipFile string) {
	file, err := os.Open(zipFile)
	if err != nil {
		log.Println("Failed to open ZIP file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "midi-files.zip"))

	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Failed to send ZIP file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
