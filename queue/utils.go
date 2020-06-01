package queue

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GenQScript(i uint64, d string, c []string) string {

	exscript := d + "/" + strconv.FormatUint(i, 10) + ".qscript"

	qscript := "#!/bin/bash\nset +x\n\n"

	AppendToFile(exscript, qscript)

	for _, cmd := range c {
		AppendToFile(exscript, cmd+"\n")
	}

	log.Info("Created qscript.")

	return exscript
}

func AppendToFile(name string, data string) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Seek(0, os.SEEK_END)
	io.WriteString(file, data)
}

func RemoveDirContents(d string) error {
	log.Debug("Going to clean: ", d)
	dir, err := os.Open(d)
	if err != nil {
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(d, name))
		if err != nil {
			return err
		}
	}
	return nil
}
