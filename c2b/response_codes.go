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

package c2b

const (
	success               = "00"
	authenticationFailure = "36"
	incorrectMSISDN       = "41"
	missingParameters     = "40"
	transactionMismatch   = "12"
	duplicateRequest      = "35"
	invalidAmount         = "31"
	maxAmountReached      = "04"
	minAmountReached      = "03"
	maxDailyAmountReached = "08"
	invalidMerchID        = "09"
	insufficientFunds     = "01"
	transactionExpired    = "05"
	confirmFailure        = "06"
	resolveMSISDNFailed   = "10"
	unableToComplete      = "11"
	downtime              = "29"
	missingReferenceID    = "30"
	inactiveAccount       = "32"
	unapprovedAccount     = "33"
	processingDelay       = "34"
)

//ResolveCode returns the Mpesa documented message for the given response codes
func ResolveCode(code string) string {
	switch code {
	case processingDelay:
		return "System processing delay"
	case unapprovedAccount:
		return "Customer account not approved to transact"
	case inactiveAccount:
		return "Customer account inactive"
	case downtime:
		return "System experiening downtime"
	case unableToComplete:
		return "Unable to complete transaction"
	case resolveMSISDNFailed:
		return "Customer MSISDN not on mpesa"
	case confirmFailure:
		return "Failed to confirm transaction"
	case transactionExpired:
		return "Transaction has expired"
	case insufficientFunds:
		return "Customer has insufficient funds"
	case invalidMerchID:
		return "Unknown Paybill/BuyGoods Number"
	case maxDailyAmountReached:
		return "Customer maximum daily amount reached"
	default:
		return "Unknown error code"

	}
}
