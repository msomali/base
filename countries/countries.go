/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package countries

import (
	"fmt"
	"strings"
)

const (
	Uganda                  = "UGANDA"
	Nigeria                 = "NIGERIA"
	Tanzania                = "TANZANIA"
	Kenya                   = "KENYA"
	Rwanda                  = "RWANDA"
	Zambia                  = "ZAMBIA"
	Gabon                   = "GABON"
	Niger                   = "NIGER"
	Brazzaville             = "CONGO-BRAZZAVILLE"
	DrCongo                 = "DR CONGO"
	CHAD                    = "CHAD"
	Seychelles              = "SEYCHELLES"
	Madagascar              = "MADAGASCAR"
	Malawi                  = "MALAWI"
	UgandaCodeName          = "UG"
	NigeriaCodeName         = "NG"
	TanzaniaCodeName        = "TZ"
	KenyaCodeName           = "KE"
	RwandaCodeName          = "RW"
	ZambiaCodeName          = "ZM"
	GabonCodeName           = "GA"
	NigerCodeName           = "NE"
	BrazzavilleCodeName     = "CG"
	DrCongoCode             = "CD"
	ChadCodeName            = "CFA"
	SeychellesCodeName      = "SC"
	MadagascarCodeName      = "MG"
	MalawiCodeName          = "MW"
	UgandaCurrencyCode      = "UGX"
	NigeriaCurrencyCode     = "NGN"
	TanzaniaCurrencyCode    = "TZS"
	KenyaCurrencyCode       = "KES"
	RwandaCurrencyCode      = "RWF"
	ZambiaCurrencyCode      = "ZMW"
	GabonCurrencyCode       = "CFA"
	NigerCurrencyCode       = "XOF"
	BrazzavilleCurrencyCode = "XAF"
	DrCongoCurrencyCode     = "CDF"
	ChadCurrencyCode        = "XAF"
	SeychellesCurrencyCode  = "SCR"
	MadagascarCurrencyCode  = "MGA"
	MalawiCurrencyCode      = "MWK"
	UgandaCurrency          = "Ugandan shilling"
	NigeriaCurrency         = "Nigerian naira"
	TanzaniaCurrency        = "Tanzanian shilling"
	KenyaCurrency           = "Kenyan shilling"
	RwandaCurrency          = "Rwandan franc"
	ZambiaCurrency          = "Zambian kwacha"
	GabonCurrency           = "CFA franc BEAC"
	NigerCurrency           = "CFA franc BCEAO"
	BrazzavilleCurrency     = "CFA franc BCEA"
	DrCongoCurrency         = "Congolese franc"
	ChadCurrency            = "CFA franc BEAC"
	SeychellesCurrency      = "Seychelles rupee"
	MadagascarCurrency      = "Malagasy ariary"
	MalawiCurrency          = "Malawian Kwacha"
)

type (
	Country struct {
		CommonName   string
		CodeName     string
		CurrencyName string
		CurrencyCode string
	}
)

func Names() []string {
	names := []string{
		Uganda, Niger, Nigeria, Tanzania, Kenya, Rwanda,
		Zambia, Gabon, Brazzaville, DrCongo, CHAD,
		Seychelles, Madagascar, Malawi,
	}

	return names
}

func Get(name string) (Country, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	countries := list()
	for _, country := range countries {
		n := strings.TrimSpace(strings.ToLower(country.CommonName))
		if name == n {
			return country, nil
		}
	}

	return Country{}, fmt.Errorf("error: the country %s is not supported", name)
}

func list() []Country {

	var countries []Country
	var (
		uganda = Country{
			CommonName:   Uganda,
			CodeName:     UgandaCodeName,
			CurrencyName: UgandaCurrency,
			CurrencyCode: UgandaCurrencyCode,
		}

		sych = Country{
			CommonName:   Seychelles,
			CodeName:     SeychellesCodeName,
			CurrencyName: SeychellesCurrency,
			CurrencyCode: SeychellesCurrencyCode,
		}

		brazzaville = Country{
			CommonName:   Brazzaville,
			CodeName:     BrazzavilleCodeName,
			CurrencyName: BrazzavilleCurrency,
			CurrencyCode: BrazzavilleCurrencyCode,
		}

		kenya = Country{
			CommonName:   Kenya,
			CodeName:     KenyaCodeName,
			CurrencyName: KenyaCurrency,
			CurrencyCode: KenyaCurrencyCode,
		}

		nigeria = Country{
			CommonName:   Nigeria,
			CodeName:     NigeriaCodeName,
			CurrencyName: NigeriaCurrency,
			CurrencyCode: NigeriaCurrencyCode,
		}

		rwanda = Country{
			CommonName:   Rwanda,
			CodeName:     RwandaCodeName,
			CurrencyName: RwandaCurrency,
			CurrencyCode: RwandaCurrencyCode,
		}

		niger = Country{
			CommonName:   Nigeria,
			CodeName:     NigerCodeName,
			CurrencyName: NigerCurrency,
			CurrencyCode: NigerCurrencyCode,
		}

		chad = Country{
			CommonName:   CHAD,
			CodeName:     ChadCodeName,
			CurrencyName: ChadCurrency,
			CurrencyCode: ChadCurrencyCode,
		}

		congo = Country{
			CommonName:   DrCongo,
			CodeName:     DrCongoCode,
			CurrencyName: DrCongoCurrency,
			CurrencyCode: DrCongoCurrencyCode,
		}

		madagascar = Country{
			CommonName:   Madagascar,
			CodeName:     MadagascarCodeName,
			CurrencyName: MadagascarCurrency,
			CurrencyCode: MadagascarCurrencyCode,
		}

		zambia = Country{
			CommonName:   Zambia,
			CodeName:     ZambiaCodeName,
			CurrencyName: ZambiaCurrency,
			CurrencyCode: ZambiaCurrencyCode,
		}

		gabon = Country{
			CommonName:   Gabon,
			CodeName:     GabonCodeName,
			CurrencyName: GabonCurrency,
			CurrencyCode: GabonCurrencyCode,
		}
		tz = Country{
			CommonName:   Tanzania,
			CodeName:     TanzaniaCodeName,
			CurrencyName: TanzaniaCurrency,
			CurrencyCode: TanzaniaCurrencyCode,
		}
		malawi = Country{
			CommonName:   Malawi,
			CodeName:     MalawiCodeName,
			CurrencyName: MalawiCurrency,
			CurrencyCode: MalawiCurrencyCode,
		}
	)

	countries = append(countries, uganda, malawi, tz, gabon, congo,
		brazzaville, rwanda, kenya, madagascar, zambia, sych, chad, niger, nigeria)

	return countries
}

func Search(countryName string) bool {
	countryName = strings.TrimSpace(strings.ToLower(countryName))
	countries := list()
	for _, country := range countries {
		n := strings.TrimSpace(strings.ToLower(country.CommonName))
		if countryName == n {
			return true
		}
	}
	return false
}

func GetCodeName(countryName string) (string, error) {
	countryName = strings.TrimSpace(strings.ToLower(countryName))
	countries := list()
	for _, country := range countries {
		n := strings.TrimSpace(strings.ToLower(country.CommonName))
		if countryName == n {
			return country.CodeName, nil
		}
	}
	return "", fmt.Errorf("error: the country %s is not supported", countryName)
}

func GetCurrencyCode(countryName string) (string, error) {
	countryName = strings.TrimSpace(strings.ToLower(countryName))
	countries := list()
	for _, country := range countries {
		n := strings.TrimSpace(strings.ToLower(country.CommonName))
		if countryName == n {
			return country.CurrencyCode, nil
		}
	}
	return "", fmt.Errorf("error: the country %s is not supported", countryName)
}
