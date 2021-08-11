package utility

import (
    "bytes"
    "io/ioutil"
)

// Util type define the Util struct
//
// OsReadFile defines a abstract function that is used for reading files. This makes it possible to overwrite the default used ioutil.ReadFile function with an other implementation that for example just returns static string without reading files or reading remote files for example.
type Util struct {
    OsReadFile func(filename string) ([]byte, error)
}

// InMemoryFileReader define type for reading files from memory instead of the filesystem
type InMemoryFileReader struct {
    Str string
}

// ReadFile read the file content from a string in the memory instead of the filesystem
func (inMemReader InMemoryFileReader) ReadFile(filename string) ([]byte, error) {
    buf := bytes.NewBufferString(inMemReader.Str)
    return ioutil.ReadAll(buf)
}

//NewUtil instantiate the default Util type.
func NewUtil() Util {
    return Util{
        OsReadFile: ioutil.ReadFile,
    }
}

// NewInMemoryUtil instantiate the Util type to read from memory instead from the file system
func NewInMemoryUtil(inMemReader InMemoryFileReader) Util {
    return Util{
        OsReadFile: inMemReader.ReadFile,
    }
}

//ReadFile returns the file content of the passed fileName.
func (util Util) ReadFile(fileName string) (*string, error) {
    b, err := util.OsReadFile(fileName) // just pass the file name
    if err != nil {
        return nil, err
    }
    yaml := string(b)
    return &yaml, nil
}

//ReadFile returns the file content of the passed fileName.
func ReadFile(fileName string) (*string, error) {
    return NewUtil().ReadFile(fileName) // just pass the file name
}
