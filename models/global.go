package models

import (
	"sync"
	"logs"
	"runtime/debug"
)

const cacheLength = 10000

var (
	ApiUpdatePoll     int
	ApiCacheMin		 int
	ApiCacheMax       int
	ApiCacheThreshold float64
)

type UserCache struct {
	InfoCache [cacheLength]*User
	LRU       [cacheLength]int64
	Index     int
}

func (u *UserCache) Check(token string) {

}

var ApiCacheInstance apiCache

type apiCache struct {
	api     map[string]bool
	once    sync.Once
	mutex   sync.Mutex
	updated bool
}

func (a *apiCache) Init() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var err error
	a.api, err = GetApiList()
	if err != nil {
		panic(err)
	}
}

func (a *apiCache) SetUpdated(b bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var err error
	a.api, err = GetApiList()
	if err != nil {
		logs.Error(err, debug.Stack())
	}
}

func (a *apiCache) GetApi(api string) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.api[api] {
		return true
	}
	return false
}

/**

// 功能废弃
func (a *apiCache) Init() {
	a.once.Do(func() {
		a.api = make(map[string]bool)
	})
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.updated || len(a.api) == 0 {
		ret, err := utils.ServiceCall(conf.GetString("client-center.app"), "getWxClientApi", "get", nil, nil)
		if err != nil {
			return
		}
		for k := range ret {
			a.api[k] = true
		}
		a.updated = false
	}
}
func (a *apiCache) GetApi(api string) bool {
	//废弃
	//a.Init()
	a.mutex.Lock()
	defer a.mutex.Unlock()
	poll, ok := a.api[api]
	if ok && poll == ApiUpdatePoll {
		return true
	}
	return false
}

func (a *apiCache) SetUpdated(b bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.updated = b
}

func PollingApi() {
	ApiCacheMin = conf.GetIntegerW("api.cache.min", 8)
	ApiCacheMax = conf.GetIntegerW("api.cache.max", 16)
	ApiCacheThreshold = conf.GetFloatW("api.cache.threshold", 1.4)
	ApiCacheInstance.api = make(map[string]int, ApiCacheMin)
	apiConfig := conf.GetStringW("api.file", "apiConf")
	logs.Info("api min:", ApiCacheMin)
	logs.Info("api max:", ApiCacheMax)
	logs.Info("threshold:", ApiCacheThreshold)
	for {
		logs.Info("polling api ...")
		file, err := os.Open(apiConfig)
		if err != nil {
			logs.Error(string(debug.Stack()))
			return
		}
		reader := bufio.NewReader(file)
		line, _, err := reader.ReadLine()
		update := string(line)
		logs.Info(update)
		if err != nil || strings.TrimSpace(update) == "" || !strings.HasPrefix(update, "update") {
			if err != nil {
				logs.Error(err, string(debug.Stack()))
			}
			return
		}
		data := strings.Split(update, "=")
		if len(data) != 2 {
			logs.Error("format error [eg: update=2]")
			return
		}
		newPoll := utils.ParseInt(strings.TrimSpace(data[1]))
		if newPoll > ApiUpdatePoll {
			ApiUpdatePoll = newPoll
			updateApi(reader)
		}
		file.Close()
		//<- time.Tick(time.Second * 60)
		time.Sleep(time.Minute)
	}
	return
}

func updateApi(reader *bufio.Reader) {
	ApiCacheInstance.mutex.Lock()
	defer ApiCacheInstance.mutex.Unlock()
	line, _, err := reader.ReadLine()
	var thisTimeNum int
	for line != nil && err == nil {
		api := strings.TrimSpace(string(line))
		if api != "" && !strings.HasPrefix(api, "#") {
			logs.Info(api)
			ApiCacheInstance.api[api] = ApiUpdatePoll
			thisTimeNum++
		}
		line, _, err = reader.ReadLine()
	}

	// 缓存容量达到最大值
	if len(ApiCacheInstance.api) > ApiCacheMax {
		if thisTimeNum >= ApiCacheMax{
			logs.Info("api cache max increase")
			ApiCacheMax = int(float64(thisTimeNum) * ApiCacheThreshold)
		} else {
			logs.Info("delete invalid api")
			for k, v := range ApiCacheInstance.api{
				if v < ApiUpdatePoll {
					delete(ApiCacheInstance.api, k)
				}
			}
		}
	}
}

*/