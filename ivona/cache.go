package ivona

import(
  log "github.com/Sirupsen/logrus"
  "sync"
  "crypto/sha256"
  "encoding/hex"
  "time"
  "sort"
)

type ivonaCachedItem struct{
  Hash    string
  Text    string
  Audio   []byte
  Used    time.Time
}

type ivonaCachedItems []ivonaCachedItem
type cachedByDate []ivonaCachedItem

const maxCacheLength int = 100
var cacheMutex = &sync.Mutex{}
var Cache ivonaCachedItems

func hash(text string) string{
  c := sha256.New()
  return hex.EncodeToString(c.Sum([]byte(text)))
}

func (c ivonaCachedItems) checkCache(text string) (bool, []byte){
  textHash := hash(text)
  var foundInCache bool
  item := ivonaCachedItem{}
  cacheMutex.Lock()
  for i,ci := range c{
    if ci.Hash == textHash{
      ci.Used = time.Now()
      item = ci
      foundInCache = true

      Cache[i] = ci
      break
    }
  }
  cacheMutex.Unlock()
  log.WithFields(log.Fields{
    "hash": textHash,
    "inCache": foundInCache,
    "used": item.Used,
    "cacheSize": len(c),
  }).Debug("Cache check for Ivona call")
  return foundInCache, item.Audio
}

func (c ivonaCachedItems) addCache(text string, audio []byte){
  cacheMutex.Lock()
  ci := ivonaCachedItem{}
  ci.Hash = hash(text)
  ci.Audio = audio
  ci.Used = time.Now()
  ci.Text = text

  c = append(c, ci)

  sort.Sort(cachedByDate(c))
  if len(c) >= maxCacheLength{
    for n,_ := range c{
      if n > maxCacheLength{
        c = append(c[:n], c[n+1:]...)
      }
    }
  }

  Cache = c
  cacheMutex.Unlock()
}

func (c cachedByDate) Len() int           { return len(c) }
func (c cachedByDate) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c cachedByDate) Less(i, j int) bool { return c[i].Used.Before(c[j].Used) }
