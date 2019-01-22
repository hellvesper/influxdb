// Code generated by go-bindata. DO NOT EDIT.
// sources:
// helloworld.json (869B)

package proto

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _helloworldJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xcd\x8e\xab\x20\x14\xde\xfb\x14\xe7\xb2\xbe\x4d\x40\x6b\xed\x65\x77\x17\xb3\x9d\x55\x93\x49\x66\xd2\x05\xd5\xd3\x62\x82\xe2\x00\xfd\x31\x8d\xef\x3e\x11\x53\x07\x49\xeb\x82\x90\xef\x8f\x0f\x39\xf7\x04\x80\xd4\x15\xe1\x40\xe8\xf2\x63\xe4\xef\xc8\xb5\xa2\xc1\x91\xb5\xbd\x75\xd8\x4c\x58\x25\xac\x3c\x68\x61\x2a\x4b\x38\x7c\x25\x00\x00\x77\xbf\x86\x1c\xe1\x33\x18\xc4\x48\x54\x4a\xfb\x94\x87\x1e\x6d\x69\xea\xce\xd5\xba\x1d\x79\x2d\xc1\x4b\xc0\x49\x34\xf8\x27\x54\x96\xa8\xd4\xef\x81\xd3\x77\x0f\xf6\xf3\x45\x2a\x51\x88\x83\xc8\x31\xaf\xb6\x0c\xf3\x9c\x06\x21\x5e\x75\x23\x1c\x58\x84\xf5\x84\x43\x1a\x61\x57\xc2\x21\x8b\x30\x49\x38\xac\x03\x68\x98\xf7\xfb\xa0\x6a\x83\x4e\x2c\xee\x3f\xd6\x37\x28\x1c\x56\xff\xdd\x58\x31\xa5\xf4\xdf\x8a\xb1\x15\xa3\xbb\x34\xe3\x94\x72\x4a\x3f\x17\x35\xc9\xb9\xab\x9e\xc8\xd9\xce\x6b\xbd\x3c\x89\x3b\x0c\x8f\x00\x72\xa9\xf1\x6a\x97\x0f\xc0\xf2\x62\x9d\x31\x4a\xb3\x0d\xdb\x6e\xf3\x75\xb1\x61\xc5\x26\xae\xf8\xea\x91\x00\x48\x67\x74\x87\xc6\xd5\x68\x23\x13\x00\xb1\x52\x74\xde\x57\x4a\xa3\x5b\x7d\x32\xe2\xb8\xba\xa4\xf1\x5f\xff\x3e\xa3\x99\xec\xed\x59\xa9\x88\x14\xb7\x57\x8c\xeb\xa7\xec\x5b\x1f\x07\x2a\x3c\x61\xeb\xc7\x6c\x88\x98\x13\xea\x66\xf4\xc4\x8e\x52\x2b\x6d\x5e\x9c\xd3\x6a\x87\xcf\x3c\x56\xea\xeb\xbb\x76\xf8\x21\xb1\x7d\x6b\x3a\x37\x0e\xca\x51\x28\x8b\x4f\x87\x60\x7e\x8a\xe4\xb1\xee\x93\x21\xf9\x09\x00\x00\xff\xff\xae\x3d\x2d\xfd\x65\x03\x00\x00")

func helloworldJsonBytes() ([]byte, error) {
	return bindataRead(
		_helloworldJson,
		"helloworld.json",
	)
}

func helloworldJson() (*asset, error) {
	bytes, err := helloworldJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "helloworld.json", size: 869, mode: os.FileMode(420), modTime: time.Unix(1547843603, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x2a, 0xbe, 0x16, 0xfd, 0xe2, 0x7d, 0x82, 0x6f, 0xd3, 0x28, 0xdf, 0x7a, 0x36, 0x5d, 0x9f, 0x81, 0xf2, 0x99, 0xb1, 0x45, 0xce, 0xec, 0x1c, 0xba, 0xe, 0x6b, 0xa5, 0xcb, 0xe5, 0x3f, 0x9e, 0x58}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"helloworld.json": helloworldJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"helloworld.json": &bintree{helloworldJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
//lint:file-ignore ST1005 Ignore error strings should not be capitalized