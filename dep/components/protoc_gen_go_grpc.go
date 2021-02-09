package components

import (
	"go.amplifyedge.org/booty-v2/pkg/downloader"
	"go.amplifyedge.org/booty-v2/pkg/fileutil"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"go.amplifyedge.org/booty-v2/pkg/store"

	"fmt"
	"os"
	"path/filepath"
)

const (
	// version -- version -- os.arch
	genGrpcUrlFormat = "https://github.com/grpc/grpc-go/releases/download/cmd/protoc-gen-go-grpc/v%s/protoc-gen-go-grpc.v%s.%s.tar.gz"
)

type ProtocGenGoGrpc struct {
	version string
	dlPath  string
	db      *store.DB
}

func (p *ProtocGenGoGrpc) Name() string {
	return "protoc-gen-go-grpc"
}

func (p *ProtocGenGoGrpc) Version() string {
	return p.version
}

func (p *ProtocGenGoGrpc) Download(targetDir string) error {
	if osutil.GetArch() != "amd64" {
		return fmt.Errorf("error: unsupported arch: %v", osutil.GetArch())
	}
	osName := fmt.Sprintf("%s.%s", osutil.GetOS(), osutil.GetArch())
	fetchUrl := fmt.Sprintf(genGrpcUrlFormat, p.version, p.version, osName)
	fmt.Println("DOWNLOADING FROM: ", fetchUrl)
	target := filepath.Join(targetDir, p.Name()+"-"+p.version)
	err := downloader.Download(fetchUrl, target)
	if err != nil {
		return err
	}
	p.dlPath = target
	return nil
}

func (p *ProtocGenGoGrpc) Install() error {
	var err error
	// install to path
	goBinDir := filepath.Join(osutil.GetGoPath(), "bin")
	_ = os.MkdirAll(goBinDir, 0755)

	executableName := p.Name()
	if osutil.GetOS() == "windows" {
		executableName += ".exe"
	}

	// all files that are going to be installed
	filesMap := map[string][]interface{}{
		filepath.Join(p.dlPath, executableName): {filepath.Join(goBinDir, executableName), 0755},
	}

	ip := store.InstalledPackage{
		Name:     p.Name(),
		Version:  p.version,
		FilesMap: map[string]int{},
	}

	// copy file to the bin directory
	for k, v := range filesMap {
		if err = fileutil.Copy(k, v[0].(string)); err != nil {
			return err
		}
		installedName := v[0].(string)
		installedMode := v[1].(int)
		if err = os.Chmod(installedName, os.FileMode(installedMode)); err != nil {
			return err
		}
		ip.FilesMap[installedName] = installedMode
	}
	if err = p.db.New(&ip); err != nil {
		return err
	}
	return os.RemoveAll(p.dlPath)
}

func (p *ProtocGenGoGrpc) Uninstall() error {
	var err error
	var pkg *store.InstalledPackage
	pkg, err = p.db.Get(p.Name())
	if err != nil {
		return err
	}
	var filesList []string
	for k := range pkg.FilesMap {
		filesList = append(filesList, k)
	}
	for _, file := range filesList {
		if err = os.RemoveAll(file); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProtocGenGoGrpc) Run(args ...string) error {
	// Should not run as standalone, rather to compliment protoc
	return nil
}

func (p *ProtocGenGoGrpc) Update(version string) error {
	p.version = version
	targetDir := filepath.Dir(p.dlPath)
	if err := p.Uninstall(); err != nil {
		return err
	}
	if err := p.Download(targetDir); err != nil {
		return err
	}
	return p.Install()
}

func (p *ProtocGenGoGrpc) RunStop() error {
	return nil
}

func (p *ProtocGenGoGrpc) Backup() error {
	return nil
}

func NewProtocGenGoGrpc(db *store.DB, version string) *ProtocGenGo {
	return &ProtocGenGo{
		version: version,
		dlPath:  "",
		db:      db,
	}
}