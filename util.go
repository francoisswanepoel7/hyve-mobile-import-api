package main

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/jpeg"
	"net"
	"os"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

func getDomainFromEmail(email string) string {
	index := strings.Index(email, "@")
	if index > -1 {
		return email[index+1:]
	}
	return ""
}

func getTZConstituents(tz string) (string, string) {
	index := strings.Index(tz, "/")
	if index > -1 {
		return tz[:index], tz[index+1:]
	}
	return "", ""
}

func ipLookup(domain string) string {
	addr, err := net.LookupIP(domain)
	if err != nil {
		return "127.0.0.1"
	}

	return addr[0].String()
}

func addText(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func createImage(email string, fn string, ln string) {
	img := image.NewRGBA(image.Rect(0, 0, 600, 600))
	addText(img, 40, 40, email)
	//addText(img, 20, 50, fn + " " + ln)

	f, err := os.Create("images/" + email + ".jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 1}); err != nil {
		panic(err)
	}
}
