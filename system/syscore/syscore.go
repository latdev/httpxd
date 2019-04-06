package syscore

import (
	"encoding/json"
	"encoding/xml"
	"github.com/jinzhu/gorm"
	"github.com/latdev/httpxd/system/models"
	"github.com/latdev/httpxd/system/syshelper"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const SessionCookieName = "gs"

// CoreSettings contains all settings stored in xml file
type CoreSettings struct {
	XMLName xml.Name `xml:"settings"`
	Server struct {
		Binding        string `xml:"binding"`
		ReadTimeout    uint16 `xml:"readTimeout"`
		WriteTimeout   uint16 `xml:"writeTimeout"`
		MaxHeaderBytes uint32 `xml:"MaxHeaderBytes"`
	} `xml:"server"`
	DatabaseString string `xml:"mysql"`
}

// CoreDirector allows to engine components interact with other modules
type CoreDirector struct {
	*gorm.DB
	Settings *CoreSettings
}

type SessionManager struct {
	lock      *sync.Mutex
	cd        *CoreDirector
	wr        http.ResponseWriter
	req       *http.Request
	sessionId string
	data      *syshelper.Serialized
}

func New() (*CoreDirector, error) {
	var result = &CoreDirector{
		Settings: &CoreSettings{},
	}

	var loadDefaultSettings = true
	if file, err := os.Open("settings.xml"); err == nil {
		defer file.Close()
		if buffer, err := ioutil.ReadAll(file); err == nil {
			if err := xml.Unmarshal(buffer, result.Settings); err == nil {
				loadDefaultSettings = false
			}
		}
	}
	if loadDefaultSettings {
		func (s *CoreSettings) {
			s.Server.Binding = "0.0.0.0:10000"
			s.Server.ReadTimeout = 15
			s.Server.WriteTimeout = 30
			s.Server.MaxHeaderBytes = 2048
			s.DatabaseString = ""
		}(result.Settings)
	}

	dbc, err := gorm.Open("mysql", result.Settings.DatabaseString)
	if err != nil {
		return nil, err
	}
	if err = models.ForceModelsCreation(dbc); err != nil {
		return nil, err
	}
	result.DB = dbc

	return result, nil
}

// Close will close all links (use only once at all end)
func (cd *CoreDirector) Close() {
}

// MuxServeNotFound exports error 404 handler for mux.NewRouter().NotFoundHandler = director.MuxServeNotFound()
func (*CoreDirector) MuxServeNotFound() http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		bytes, err := json.Marshal(struct {
			Success int    `json:"success"`
			Message string `json:"message"`
		}{
			Success: 0,
			Message: "application not found",
		})
		if err == nil {
			wr.WriteHeader(404)
			wr.Header().Set("Content-Type", "application/json")
			wr.Write(bytes)
		} else {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (cd *CoreDirector) Session(wr http.ResponseWriter, req *http.Request) (*SessionManager) {
	var result = &SessionManager{
		lock:      &sync.Mutex{},
		cd:        cd,
		wr:        wr,
		req:       req,
		sessionId: "#",
		data:      &syshelper.Serialized{},
	}

	var cookieUnsafe = true
	cookie, err := req.Cookie(SessionCookieName)
	if err == nil {
		result.sessionId = cookie.Value
		var session = &models.Session{}
		if err := cd.Where("session = ?", result.sessionId).First(session).Error; err == nil {
			if data, err := syshelper.DeserializeStruct(session.Value); err == nil {
				cookieUnsafe = false
				result.data = data
			}
		}
	}
	if cookieUnsafe {
		result.sessionId = syshelper.GenerateNewSessionId()
		http.SetCookie(wr, &http.Cookie{
			Name: SessionCookieName,
			Value: result.sessionId,
			Expires: time.Now().Add(time.Hour),
			HttpOnly: true,
		})
		result.data = &syshelper.Serialized{}
	}
	return result
}

func (sm *SessionManager) Save() (error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if data64, err := syshelper.SerializeStruct(sm.data); err == nil {
		var session = &models.Session{}
		if err := sm.cd.Where("session = ?", sm.sessionId).First(session).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				session.Session = sm.sessionId
				session.LastWrite = time.Now()
				session.Value = data64
				sm.cd.Create(session)
			} else {
				return errors.Wrap(err, "db write error")
			}
		} else {
			sm.cd.Model(&session).Where("session = ?", sm.sessionId).Update("value", data64)
		}
	} else {
		return errors.Wrap(err, "serialize error")
	}
	return nil
}

func (sm *SessionManager) Get(name string) (value interface{}, ok bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	value, ok = map[string]interface{}(*sm.data)[name]
	return
}

func (sm *SessionManager) Set(name string, value interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	map[string]interface{}(*sm.data)[name] = value
}


