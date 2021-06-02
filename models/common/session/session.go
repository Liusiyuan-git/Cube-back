package session

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/server/web/session"
	"net/http"
)

var globalSessions *session.Manager

func SetSession(ctx context.Context, w http.ResponseWriter, r *http.Request, CubeId string) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(ctx, w)
	sess.Set(ctx, "CubeId", CubeId)
}

func GetSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(ctx, w)
	username := sess.Get(ctx, "username")
	fmt.Print(username)
}

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "CubeSessionId",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}
