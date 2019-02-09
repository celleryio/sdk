/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fatih/color"

	"github.com/cellery-io/sdk/components/cli/pkg/constants"
)

var Cyan = color.New(color.FgCyan)
var CyanF = color.New(color.FgCyan).SprintFunc()
var CyanBold = Cyan.Add(color.Bold).SprintFunc()
var Bold = color.New(color.Bold).SprintFunc()
var Faint = color.New(color.Faint).SprintFunc()
var Green = color.New(color.FgGreen)
var GreenBold = Green.Add(color.Bold).SprintFunc()
var Yellow = color.New(color.FgYellow)
var YellowBold = Yellow.Add(color.Bold).SprintFunc()

// ParseImage parses the given image name string and returns a CellImage struct with the relevant information.
func ParseImage(cellImageString string) (parsedCellImage *CellImage, err error) {
	cellImage := &CellImage{
		constants.CENTRAL_REGISTRY_HOST,
		"",
		"",
		"",
	}

	if cellImageString == "" {
		return cellImage, errors.New("no cell image specified")
	}

	const IMAGE_FORMAT_ERROR_MESSAGE = "incorrect image name format. Image name should be " +
		"[REGISTRY[:REGISTRY_PORT]/]ORGANIZATION/IMAGE_NAME:VERSION"

	// Parsing the cell image string
	strArr := strings.Split(cellImageString, "/")
	if len(strArr) == 3 {
		cellImage.Registry = strArr[0]
		cellImage.Organization = strArr[1]
		imageTag := strings.Split(strArr[2], ":")
		if len(imageTag) != 2 {
			return cellImage, errors.New(IMAGE_FORMAT_ERROR_MESSAGE)
		}
		cellImage.ImageName = imageTag[0]
		cellImage.ImageVersion = imageTag[1]
	} else if len(strArr) == 2 {
		cellImage.Organization = strArr[0]
		imageNameSplit := strings.Split(strArr[1], ":")
		if len(imageNameSplit) != 2 {
			return cellImage, errors.New(IMAGE_FORMAT_ERROR_MESSAGE)
		}
		cellImage.ImageName = imageNameSplit[0]
		cellImage.ImageVersion = imageNameSplit[1]
	} else {
		return cellImage, errors.New(IMAGE_FORMAT_ERROR_MESSAGE)
	}

	return cellImage, nil
}

func PrintWhatsNextMessage(action string, cmd string) {
	fmt.Println()
	fmt.Println(Bold("What's next?"))
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("Execute the following command to %s:\n", action)
	fmt.Println("  $ " + cmd)
	fmt.Println("--------------------------------------------------------")
}

//func CopyFile(oldFile string, newFile string) {
//	input, err := os.Open(oldFile)
//	if err != nil {
//		panic(err)
//	}
//	defer input.Close()
//
//	output, err := os.Create(newFile)
//	if err != nil {
//		panic(err)
//	}
//	defer output.Close()
//
//	_, err = io.Copy(output, input)
//	if err != nil {
//		panic(err)
//	}
//}

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func Trim(stream string) string {
	var trimmedString string
	if strings.Contains(stream, ".cell.balx") {
		trimmedString = strings.Replace(stream, ".cell.balx", ".celx", -1)
	} else if strings.Contains(stream, ".bal") {
		trimmedString = strings.Replace(stream, ".bal", "", -1)
	} else if strings.Contains(stream, ".cell") {
		trimmedString = strings.Replace(stream, ".cell", "", -1)
	} else {
		trimmedString = stream
	}
	return trimmedString
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func RecursiveZip(files []string, folders []string, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	for _, folder := range folders {
		err = filepath.Walk(folder, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			relPath := strings.TrimPrefix(filePath, filepath.Dir(folder))

			zipFile, err := myZip.Create(relPath)
			if err != nil {
				return err
			}

			fsFile, err := os.Open(filePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, fsFile)

			if err != nil {
				return err
			}
			return nil
		})
	}

	// Copy files
	for _, file := range files {
		zipFile, err := myZip.Create(file)
		files, err := os.Open(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, files)
	}
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func Unzip(zipFolderName string, destinationFolderName string) error {
	var fileNames []string
	zipFolder, err := zip.OpenReader(zipFolderName)
	if err != nil {
		return err
	}
	defer zipFolder.Close()

	for _, file := range zipFolder.File {
		fileContent, err := file.Open()
		if err != nil {
			return err
		}
		defer fileContent.Close()

		fpath := filepath.Join(destinationFolderName, file.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(destinationFolderName)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		fileNames = append(fileNames, fpath)
		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, fileContent)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func FindInDirectory(directory, suffix string) []string {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil
	}
	fileList := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			fileList = append(fileList, filepath.Join(directory, f.Name()))
		}
	}
	return fileList
}

func GetDuration(startTime time.Time) string {
	duration := ""
	var year, month, day, hour, min, sec int
	currentTime := time.Now()
	if startTime.Location() != currentTime.Location() {
		currentTime = currentTime.In(startTime.Location())
	}
	if startTime.After(currentTime) {
		startTime, currentTime = currentTime, startTime
	}
	startYear, startMonth, startDay := startTime.Date()
	currentYear, currentMonth, currentDay := currentTime.Date()

	startHour, startMinute, startSecond := startTime.Clock()
	currentHour, currentMinute, currentSecond := currentTime.Clock()

	year = int(currentYear - startYear)
	month = int(currentMonth - startMonth)
	day = int(currentDay - startDay)
	hour = int(currentHour - startHour)
	min = int(currentMinute - startMinute)
	sec = int(currentSecond - startSecond)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(startYear, startMonth, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	numOfTimeUnits := 0
	if year > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(year) + " years "
		numOfTimeUnits++
	}
	if month > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(month) + " months "
		numOfTimeUnits++
	}
	if day > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(day) + " days "
		numOfTimeUnits++
	}
	if hour > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(hour) + " hours "
		numOfTimeUnits++
	}
	if min > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(min) + " minutes "
		numOfTimeUnits++
	}
	if sec > 0 && numOfTimeUnits < 2 {
		duration += strconv.Itoa(sec) + " seconds"
		numOfTimeUnits++
	}
	return duration
}

func ConvertStringToTime(timeString string) time.Time {
	convertedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error parsing time: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	return convertedTime
}

// Creates a new file upload http request with optional extra params
func FileUploadRequest(uri string, params map[string]string, paramName, path string, secure bool) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	if !secure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func DownloadFile(filepath string, url string) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(url)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		out, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// This returns the ballerina home directory in the machine.
// This expects the BALLERINA_HOME environment variable to be set.
func BallerinaHomeDir() string {
	return os.Getenv("BALLERINA_HOME")
}

func CreateDir(dirPath string) error {
	dirExist, _ := FileExists(dirPath)
	if !dirExist {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func GetSubDirectoryNames(path string) ([]string, error) {
	directoryNames := []string{}
	subdirectories, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, subdirectory := range subdirectories {
		directoryNames = append(directoryNames, subdirectory.Name())
	}
	return directoryNames, nil
}

func GetFileSize(path string) (int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return file.Size(), nil
}

func RenameFile(oldName, newName string) error {
	err := os.Rename(oldName, newName)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceFile(fileToBeReplaced, fileToReplace string) error {
	errRename := RenameFile(fileToBeReplaced, fileToBeReplaced+"-old")
	if errRename != nil {
		return errRename
	}
	errCopy := CopyFile(fileToReplace, fileToBeReplaced)
	if errCopy != nil {
		return errCopy
	}
	return nil
}

func ExecuteCommand(cmd *exec.Cmd, errorMessage string) error {
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)

	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cellery : %v: %v \n", errorMessage, err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("cellery : %v: %v \n", errorMessage, err)
		os.Exit(1)
	}
	return nil
}

func DownloadFromS3Bucket(bucket, item, path string) {
	file, err := os.Create(filepath.Join(path, item))
	if err != nil {
		fmt.Printf("Error in downloading from file: %v \n", err)
		os.Exit(1)
	}

	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(constants.AWS_REGION), Credentials: credentials.AnonymousCredentials},
	)

	// Create a downloader with the session and custom options
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
		d.Concurrency = 6
	})

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Printf("Error in downloading from file: %v \n", err)
		os.Exit(1)
	}

	fmt.Println("Download completed", file.Name(), numBytes, "bytes")
}

func ExtractTarGzFile(extractTo, archive_name string) error {
	cmd := exec.Command("tar", "-zxvf", archive_name)
	cmd.Dir = extractTo

	ExecuteCommand(cmd, "Error occured in extracting file :"+archive_name)

	return nil
}

// AddImageToBalPath extracts the cell image in a temporary location and copies the relevant ballerina files to the
// ballerina repo directory. This expects the BALLERINA_HOME environment variable to be set in th developer machine.
func AddImageToBalPath(cellImage *CellImage) {
	cellImageFile := filepath.Join(UserHomeDir(), ".cellery", "repos", cellImage.Registry,
		cellImage.Organization, cellImage.ImageName, cellImage.ImageVersion,
		cellImage.ImageName+constants.CELL_IMAGE_EXT)

	// Create temp directory
	currentTIme := time.Now()
	timestamp := currentTIme.Format("27065102350415")
	tempPath := filepath.Join(UserHomeDir(), ".cellery", "tmp", timestamp)
	err := CreateDir(tempPath)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error while saving cell image to local repo: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	defer func() {
		err = os.RemoveAll(tempPath)
		if err != nil {
			fmt.Printf("\x1b[31;1m Error while cleaning up: \x1b[0m %v \n", err)
			os.Exit(1)
		}
	}()

	// Unzipping Cellery Image
	err = Unzip(cellImageFile, tempPath)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error while saving cell image to local repo: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	balRepoDir := filepath.Join(BallerinaHomeDir(), "lib", "repo", cellImage.Organization, cellImage.ImageName,
		cellImage.ImageVersion)

	// Cleaning up the old image bal files if it already exists
	hasOldImage, err := FileExists(balRepoDir)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error occurred while removing the old cell image: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	if hasOldImage {
		err = os.RemoveAll(balRepoDir)
		if err != nil {
			fmt.Printf("\x1b[31;1m Error while cleaning up: \x1b[0m %v \n", err)
			os.Exit(1)
		}
	}

	// Copying the ballerina files to bal path
	balArtifactsDir := filepath.Join(tempPath, "artifacts", "bal")
	err = CopyDir(balArtifactsDir, balRepoDir)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error while saving cell image to local repo: \x1b[0m %v \n", err)
		os.Exit(1)
	}
}
