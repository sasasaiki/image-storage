package imageStorage

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

//ImgStorageI イメージ保管に必要なメソッドを持つ
type ImgStorageI interface {
	SaveWithFileHeader(multipart.File, *multipart.FileHeader, string, string) error
	SaveWithOriginFileName(multipart.File, string, string, string) error
}

//DirImgStorage ディレクトリに保存する
type DirImgStorage struct {
}

//SaveWithFileHeader 拡張子を指定せず、渡したファイルの拡張子を使用して保存する
func (im *DirImgStorage) SaveWithFileHeader(file multipart.File, fileHeader *multipart.FileHeader, newFileName string, directory string) error {

	e := im.SaveAsItIs(file, fileHeader.Filename, newFileName, directory)
	if e != nil {
		printError("AddWithAutoExtension()", e)
		return e
	}

	return nil
}

//SaveAsItIs 新しく保存するファイル名には拡張子を指定せず、渡したファイルの拡張子を使用して保存する
//そのまま何もせず保存しているだけなので多分何でも保存できる
// example:
// originFileName = "hoge.png"
// newFileName = "fuga"
func (im *DirImgStorage) SaveAsItIs(file multipart.File, originFileName string, newFileName string, directory string) error {
	defer file.Close()
	var storageFilePath string
	if newFileName == "" {
		storageFilePath = filepath.Join(directory, originFileName)
	} else {
		storageFilePath = filepath.Join(directory, newFileName+filepath.Ext(originFileName))
	}

	data, e := ioutil.ReadAll(file)
	if e != nil {
		printError("Add()でReadALL(file)に失敗", e)
		return e
	}

	e = ioutil.WriteFile(storageFilePath, data, 0600)
	if e != nil {
		printError("Add()でfileの保存に失敗", e)
		return e
	}

	return nil
}

// SavePngToJpeg jpegに変換して保存
// TODO:透明部分が黒くなってしまうので一旦置いとく
func SavePngToJpeg(file multipart.File, originFileExtension string, newFileName string, directory string, quality int) (*os.File, error) {
	var img image.Image
	var err error
	switch originFileExtension {
	case ".jpeg", ".jpg":
		err = errors.New("jpg don't need to jpeg")
		printError("jpgなので変換の必要がありません", err)
		return nil, err
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			printError("pngのデコードに失敗しました", err)
			return nil, err
		}
	default:
		err = errors.New("Not compatible")
		printError("png以外のファイルです", err)
		return nil, err
	}

	storageFilePath := filepath.Join(directory, newFileName+".jpg")
	out, err := os.Create(storageFilePath)
	if err != nil {
		printError("imageのCreateに失敗しました", err)
		return nil, err
	}
	defer out.Close()

	opts := &jpeg.Options{Quality: quality}
	err = jpeg.Encode(out, img, opts)
	if err != nil {
		printError("jpegへのEncodeに失敗しました", err)
		return nil, err
	}

	return out, nil
}

func (im *DirImgStorage) Update() {
}

func (im *DirImgStorage) Delete() {

}

func (im *DirImgStorage) Get() {

}

func printError(message string, e error) {
	fmt.Println("in image-storage ", message, " error occurred", e)
}