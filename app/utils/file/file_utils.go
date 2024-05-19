package file

import "os"

func WriteCsvFile(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}
