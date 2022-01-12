package gron

import (
	"github.com/roylee0704/gron"
	"time"
)

func init() {
	c := gron.New()
	c.AddFunc(gron.Every(300*time.Second), func() {
		cubeTalkNewUpdate()
		cubeTalkHotUpdate()
		cubeTalkCleanAll()
	})
	c.AddFunc(gron.Every(360*time.Second), func() {
		cubeBlogNewUpdate()
		cubeBlogHotUpdate()
		cubeBlogDetailUpdate()
		cubeBlogCollectUpdate()
	})
	c.AddFunc(gron.Every(4*time.Second), func() {
		userProfileUpdate()
	})
	c.AddFunc(gron.Every(86400*time.Second), func() {
		cubeBlogCleanAll()
		cubeTalkCleanAll()
	})
	c.Start()
}
