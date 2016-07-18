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

import "fmt"

//Config stores the apis configurable info
//Get it to behave more like the MPESA sys
type Config struct {
	MaxAmount                    float64
	MinAmount                    float64
	MerchantID                   string
	CallBackDelay                int
	SAGPasskey                   string
	MaxCustomerTransactionPerDay float64
	EnabledAPIS                  []string
}

//ToString prints a pretty representation of the conf struct
func (conf Config) ToString() string {
	return fmt.Sprintf("%+v", conf)
}
