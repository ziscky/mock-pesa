/*
Copyright (C) 2016  Eric Ziscky

    This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, write to the Free Software Foundation, Inc.,
    51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package common

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

//letterBytes for generating similar to mpesa transaction IDS
const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//Middleware used for pre request processing
type Middleware func(http.ResponseWriter, *http.Request) bool

//CheckHeader checks if the Content-Type header is correct
//most SOAP servers require this strictly
func CheckHeader() Middleware {
	return func(rw http.ResponseWriter, r *http.Request) bool {
		if !strings.Contains(r.Header.Get("Content-Type"), "text/xml") {
			log.Println(fmt.Sprintf("Invalid Header: %s Safaricom API expects: %s", r.Header.Get("Content-Type"), "text/xml"))
			return false
		}
		return true
	}
}

//GenerateMpesaTrx generate mpesa transaction id
func GenerateMpesaTrx() string {
	return fmt.Sprintf("%s%d%s%d", randStringBytes(3), rand.Intn(10), randStringBytes(4), rand.Intn(100))
}

//randStringBytes generates random n number of letters from letterBytes
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
