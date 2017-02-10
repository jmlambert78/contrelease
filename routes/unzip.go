package routes

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jmlambert78/contrelease/models"
)

func GetImageFromYamlInZip(file string) string {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	imageName := ""
	// Iterate through the files in the archive,
	// printing some of their contents.
	var validYAML = regexp.MustCompile(`^[^/]+\.yaml$`)
	for _, f := range r.File {
		//fmt.Printf("Contents of %s:\n", f.Name)
		if validYAML.MatchString(f.Name) {
			//fmt.Println("YAML File found :", f.Name)

			// Search for the Image: tag in the YAML
			rc, err := f.Open()
			defer rc.Close()
			if err != nil {
				log.Fatal(err)
			}
			var validImageLine = regexp.MustCompile(`image:`)
			// create a new scanner and read the file line by line
			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				if validImageLine.MatchString(scanner.Text()) {
					a := strings.Split(scanner.Text(), "\"")
					imageName = a[1]
				}
			}
			// check for errors
			if err = scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
	}
	// Output:
	return imageName
}
func DownloadFromUrl(url string) (error, string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return err, fileName
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return err, fileName
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return err, fileName
	}

	fmt.Println(n, "bytes downloaded.")
	return err, fileName
}
func CheckZipURLForImage(url string) (error, string) {
	//_, fileName := routes.DownloadFromUrl("https://dockerhub.gemalto.com/repository/docker-delivery/risk-engine/re-cci/1.1.1.0-670/re-cci-1.1.1.0-670.zip")
	imageName := ""
	err, fileName := DownloadFromUrl(url)
	if err == nil {
		imageName = GetImageFromYamlInZip(fileName)
		fmt.Println(imageName)
	}
	return err, imageName
}
func checkAllZipImages(releases []models.Release) {
	for i, r := range releases {
		err, imageName := CheckZipURLForImage(r.CentralZipURL)
		fmt.Println("Err:", err, "image:", imageName)
		if err != nil {
			releases[i].CentralZipURL = releases[i].CentralZipURL + "(ERROR)"
		} else {
			releases[i].CentralImage = releases[i].CentralImage + "\n" + imageName
		}
	}
}
