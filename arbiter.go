package filebounce

import "crypto/rand"
import "encoding/hex"
import "errors"
import "fmt"
import "math/big"


type BounceFile struct {
  DlId, Filename, Mimetype, Token, CookieName string
  Size, sizeProgress int
  transferStarted bool
  TransferFinished bool
  targets []DownloadTarget
}

var currentBounces = map[string]*BounceFile{}

func GetBounceFile(dlId string) *BounceFile {
  bf, ex := currentBounces[dlId]
  if ex {
    return bf
  } else {
    return nil
  }
}

var idCounter = int64(1024)
var maxRand = int64(281474976710655)

func getNewIds() (dlId, cookieName string){
  dlId = ""
  for {
    _, ex := currentBounces[dlId];
    if !ex && (len(dlId) > 0) { break }

    r, _ := rand.Int(rand.Reader, big.NewInt(maxRand))
    dlId = fmt.Sprintf("%s", hex.EncodeToString(r.Bytes()))
  }
  r, _ := rand.Int(rand.Reader, big.NewInt(maxRand))
  cookieName = fmt.Sprintf("Dl-%s", hex.EncodeToString(r.Bytes()))
  return
}

func NewBounceFile(filename, mimetype, token string, size int) (bf *BounceFile) {
  bf = &BounceFile{
    Filename: filename,
    Mimetype: mimetype,
    Token: token,
    Size: size,
    sizeProgress: 0,
    transferStarted: false,
    TransferFinished: false,
    targets: []DownloadTarget{},
  }

  bf.DlId, bf.CookieName = getNewIds()
  currentBounces[bf.DlId] = bf

  return bf
}

func RegisterDownloadTarget(dlId string, target DownloadTarget) error {
  bf, ex := currentBounces[dlId]
  if !ex {
    return errors.New("Download ID does not exist")
  }

  if bf.transferStarted {
    return errors.New("File transfer already started")
  }

  bf.targets = append(bf.targets, target)
  return nil
}

func (bf *BounceFile) SendData(data []byte) error {
  if bf.TransferFinished {
    return errors.New("Transfer already finished!")
  }

  if !bf.transferStarted {
    bf.transferStarted = true
    for _, target := range bf.targets {
      target.StartFile(bf)
    }
  }

  remaining := bf.Size - bf.sizeProgress
  if len(data) > remaining {
    data = data[:remaining] // XXX: This should probably be an error
  }

  bf.sizeProgress += len(data)
  bf.TransferFinished = (bf.sizeProgress >= bf.Size)
  for _, target := range bf.targets {
    target.Stream <- data
    if (bf.TransferFinished) {
      close(target.Stream)
    }
  }
  return nil
}
