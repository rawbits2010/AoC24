package fileout

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// FileOut is for dumping text into a series of files in a dedicated directory.
// Can just dump text or continuously appnd to the current file.
// Also can add some basic technical parameters in the first lines.
type FileOut struct {
	Name            string
	basePath        string
	count           uint
	paramCount      uint
	paramLineLength uint
	currFile        *os.File
}

// NewFileOut creates a new FileOut structure while ensuring the path is exists.
func NewFileOut(basePath, name string) (*FileOut, error) {

	filePath := path.Join(basePath, name)
	exists, err := validateDirectory(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating new FileOut: %w", err)
	}

	if !exists {
		err := os.Mkdir(filePath, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating directory for new FileOut: %w", err)
		}
	}

	tempFileOut := FileOut{
		Name:     name,
		basePath: basePath,
		count:    0,
		currFile: nil,
	}

	return &tempFileOut, nil
}

// ReserveParamLines sets how many parameter lines to reserve
// in the begining of the file and how long are they.
// NOTE: No sanity check here!
func (fo *FileOut) ReserveParamLines(quantity, lineLength uint) {
	fo.paramCount = quantity
	fo.paramLineLength = lineLength
}

// StartFile creates and opens the next file. Call EndFile when you finished!
func (fo *FileOut) StartFile() error {

	filePath := fo.getFilePath()
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("couldn't create file '%s': %w", filePath, err)
	}

	for paramLineIdx := uint(0); paramLineIdx < fo.paramCount; paramLineIdx++ {
		line := make([]byte, fo.paramLineLength)

		line[0] = '/'
		line[1] = '/'

		for charIdx := uint(2); charIdx < fo.paramLineLength; charIdx++ {
			line[charIdx] = ' '
		}

		outFile.WriteString(string(line))
	}

	fo.currFile = outFile

	return nil
}

// EndFile closes the file opened by StartFile.
func (fo *FileOut) EndFile() {

	if fo.currFile != nil {

		fo.currFile.Close()
		fo.currFile = nil

		fo.count++
	}
}

func (fo *FileOut) UpdateParameter(name, value string) error {

	if fo.currFile == nil {
		return fmt.Errorf("no file is open for writing")
	}

	_, err := fo.currFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking to the top in file: '%s'", fo.getFilePath())
	}

	scanner := bufio.NewScanner(fo.currFile)

	foundParam := false
	for scanner.Scan() {
		line := scanner.Text()

		if line[:2] != "//" {
			break
		}

		tokens := strings.Split(line, " ")
		if len(tokens) < 2 {
			continue // just ignore what happened
		}

		// remove old
		if tokens[1] == name {
			fo.currFile.tru
		}
	}

}

// DumpToFile makes sure the file exists and writes 'text' into it,
// ensuring every string in the array is end in '\n'.
// It will advance the file counter.
func (fo *FileOut) DumpToFile(text []string) error {

	if fo.currFile == nil {
		return fmt.Errorf("no file is open for writing")
	}

	for _, line := range text {
		_, err := fo.currFile.WriteString(line)
		if err != nil {
			return fmt.Errorf("couldn't write to file '%s': %w", fo.getFilePath(), err)
		}
		if line[len(line)-1] != '\n' {
			_, err := fo.currFile.Write([]byte{'\n'})
			if err != nil {
				return fmt.Errorf("couldn't write to file '%s': %w", fo.getFilePath(), err)
			}
		}
	}

	return nil
}

func (fo *FileOut) getFilePath() string {
	return path.Join(fo.basePath, fo.Name, fmt.Sprintf("%s_%03d", fo.Name, fo.count))
}

//-Utils-----------------------------------------------------------------------

// validateDirectory checks if a directory exists at the given path.
func validateDirectory(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking directory '%s': %w", dirPath, err)

	}
	if !info.IsDir() {
		return false, fmt.Errorf("directory is a file '%s': %w", dirPath, err)
	}
	return true, nil
}
