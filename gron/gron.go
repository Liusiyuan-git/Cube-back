package gron

import (
	"github.com/roylee0704/gron"
	"time"
)

func init() {
	c := gron.New()
	c.AddFunc(gron.Every(3*time.Second), func() {
		cubeTalkNewUpdate()
		cubeTalkHotUpdate()
	})
	c.AddFunc(gron.Every(360*time.Second), func() {
		cubeBlogNewUpdate()
		cubeBlogHotUpdate()
		cubeBlogCollectUpdate()
	})
	c.AddFunc(gron.Every(420*time.Second), func() {
		userProfileUpdate()
	})
	c.AddFunc(gron.Every(86400*time.Second), func() {
		cubeBlogDetailClean()
		cubeTalkDetailClean()
	})
	c.Start()
}
