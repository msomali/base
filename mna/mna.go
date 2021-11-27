package mna

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var _ operator = (*Operator)(nil)

const (
	Tigo Operator = iota
	Vodacom
	TTCL
	Airtel
	Zantel
	Smile
	MoMobile
	Halotel
	Mkulima
	WiAfrica
)

const (
	registeredTigoName     = "MIC Tanzania PLC"
	registeredVodacomName  = "Vodacom Tanzania PLC"
	registeredTTCLName     = "Tanzania Telecommunications Corporation"
	registeredAirtelName   = "Airtel Tanzania PLC"
	registeredZantelName   = "Zanzibar Telecom PLC"
	registeredSmileName    = "Smile Communications Tanzania Limited"
	registeredMoMobileName = "MO Mobile Holding Limited"
	registeredHalotelName  = "Viettel Tanzania PLC"
	registeredMkulimaName  = "Mkulima African Telecommunication Company Limited"
	registeredWiAfricaName = "Wiafrica Tanzania Limited"
	commonTigoName         = "Tigo"
	commonVodacomName      = "Vodacom"
	commonTTCLName         = "TTCL"
	commonAirtelName       = "Airtel"
	commonZantelName       = "Zantel"
	commonSmileName        = "Smile"
	commonMoName           = "Mo Mobile"
	commonHalotelName      = "Halotel"
	commonMkulimaName      = "Mkulima"
	commonWiAfricaName     = "Wiafrica"
	StatusOperational      = "Operational"
)

var (
	ErrOperatorNotFound = errors.New("mobile operator not found")
	ErrInvalidFormat    = errors.New("invalid format, correct formats are +255[9-digits], 255[9-digits], 0[9-digits] or just last 9 digits")
	ErrNumericOnly      = errors.New("phone numbers should contains numeric characters only")

	tigoPrefixes     = []string{"071", "065", "067"}
	vodaPrefixes     = []string{"074", "075", "076"}
	ttclPrefixes     = []string{"073"}
	airtelPrefixes   = []string{"078", "068", "069"}
	zantelPrefixes   = []string{"077"}
	smilePrefixes    = []string{"066"}
	moPrefixes       = []string{"072"}
	viettelPrefixes  = []string{"061", "062"}
	mkulimaPrefixes  = []string{"063"}
	wiAfricaPrefixes = []string{"064"}

)

type (
	Operator int8

	// Info contains basic details of a phone number include the mno
	Info struct {
		Operator        Operator `json:"operator"`
		FormattedNumber string   `json:"formatted_number"`
	}

	FilterOperatorFunc func(op Operator) bool
	FilterPhoneFunc    func(phone string) bool

	operator interface {
		fmt.Stringer
		Prefixes() []string
		RegisteredName() string
		CommonName() string
		Status() string
	}
)

func (op Operator) Prefixes() []string {
	prefixes := [][]string{
		tigoPrefixes,
		vodaPrefixes,
		ttclPrefixes,
		zantelPrefixes,
		airtelPrefixes,
		smilePrefixes,
		moPrefixes,
		viettelPrefixes,
		mkulimaPrefixes,
		wiAfricaPrefixes,
	}

	return prefixes[op]
}

func (op Operator) RegisteredName() string {
	registeredNames := []string{
        registeredTigoName,
        registeredVodacomName,
        registeredTTCLName,
        registeredZantelName,
        registeredAirtelName,
        registeredSmileName,
        registeredMoMobileName,
        registeredHalotelName,
        registeredMkulimaName,
        registeredWiAfricaName,
    }

	return registeredNames[op]
}

func (op Operator) CommonName() string {
	commonNames := []string{
        commonTigoName,
        commonVodacomName,
        commonTTCLName,
        commonZantelName,
        commonAirtelName,
        commonSmileName,
        commonMoName,
        commonHalotelName,
        commonMkulimaName,
        commonWiAfricaName,
    }
    return commonNames[op]
}

func (op Operator) Status() string {
	return StatusOperational
}

func (op Operator) String() string {
	return fmt.Sprintf("registered name: %s, common name :%s, status: %s, prefixes :%v\n",
		op.RegisteredName(),op.CommonName(),op.Status(),op.Prefixes())
}


// format return a phone number starting with 255 or error
// if it can not be formatted.
// It tries it best to remove white spaces or hyphens put in between
// numbers and a plus sign at the beginning then replace 0 with 255
// if need be
func format(phoneNumber string) (string, error) {
	phoneNumber = strings.TrimSpace(phoneNumber)
	replacer := strings.NewReplacer(" ", "", "-", "", "+", "")
	phoneNumber = replacer.Replace(phoneNumber)
	numericOnlyRegexStr := "^[0-9]+$"
	match, err := regexp.MatchString(numericOnlyRegexStr, phoneNumber)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNumericOnly, err)
	}
	if !match {
		return "", ErrNumericOnly
	}
	phoneNumberLen := len(phoneNumber)

	withoutZero := phoneNumberLen == 9 && !strings.HasPrefix(phoneNumber, "0")
	startsWith255 := strings.HasPrefix(phoneNumber, "255") && phoneNumberLen == 12
	startsWithZero := strings.HasPrefix(phoneNumber, "0") && phoneNumberLen == 10

	if withoutZero {
		return fmt.Sprintf("255%s", phoneNumber), nil
	}

	if startsWith255 {
		return phoneNumber, nil
	}

	if startsWithZero {
		return fmt.Sprintf("255%s", phoneNumber[1:]), nil
	}

	return "", fmt.Errorf("pass the correct format")
}

func mergePrefixes() map[string]Operator {
	var m map[string]Operator
	m = make(map[string]Operator)
	for _, prefix := range tigoPrefixes {
		m[prefix] = Tigo
	}

	for _, prefix := range vodaPrefixes {
		m[prefix] = Vodacom
	}

	for _, prefix := range ttclPrefixes {
		m[prefix] = TTCL
	}

	for _, prefix := range zantelPrefixes {
		m[prefix] = Zantel
	}

	for _, prefix := range airtelPrefixes {
		m[prefix] = Airtel
	}

	for _, prefix := range mkulimaPrefixes {
		m[prefix] = Mkulima
	}

	for _, prefix := range smilePrefixes {
		m[prefix] = Smile
	}

	for _, prefix := range moPrefixes {
		m[prefix] = MoMobile
	}

	for _, prefix := range viettelPrefixes {
		m[prefix] = Halotel
	}

	for _, prefix := range wiAfricaPrefixes {
		m[prefix] = WiAfrica
	}

	return m
}

// sanitize takes a user input and tries to figure out if its in the
// right format. sanitize returns a 3 character prefix and nil error
// if the number is in the right format or an error if the number is
// bogus
func sanitize(phoneNumber string) (string, error) {
	phoneNumber = strings.TrimSpace(phoneNumber)
	replacer := strings.NewReplacer(" ", "", "-", "", "+", "")
	phoneNumber = replacer.Replace(phoneNumber)
	numericOnlyRegexStr := "^[0-9]+$"
	match, err := regexp.MatchString(numericOnlyRegexStr, phoneNumber)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNumericOnly, err)
	}
	if !match {
		return "", ErrNumericOnly
	}
	phoneNumberLen := len(phoneNumber)

	// the number len is not correct
	isWrongLen := phoneNumberLen != 9 && phoneNumberLen != 10 && phoneNumberLen != 12
	if isWrongLen {
		return "", ErrInvalidFormat
	}

	withoutZero := phoneNumberLen == 9 && !strings.HasPrefix(phoneNumber, "0")
	startsWith255 := strings.HasPrefix(phoneNumber, "255") && phoneNumberLen == 12
	startsWithZero := strings.HasPrefix(phoneNumber, "0") && phoneNumberLen == 10

	if startsWithZero {
		chars := []rune(phoneNumber)
		prefix := string(chars[0:3])
		return prefix, err
	} else if startsWith255 {
		chars := []rune(phoneNumber)
		prefix := "0" + string(chars[3:5])
		return prefix, err
	} else if withoutZero {
		chars := []rune(phoneNumber)
		prefix := string(chars[0:2])
		return fmt.Sprintf("0%s", prefix), err
	} else {
		return "", ErrInvalidFormat
	}

}

func findUsingPrefix(prefix string) (op Operator, err error) {

	m := mergePrefixes()

	op, ok := m[prefix]

	if !ok {
        return -1, ErrOperatorNotFound
    }

	return op, nil

}

func Get(phoneNumber string)(Operator,error){

    prefix, err := sanitize(phoneNumber)
    if err != nil {
        return -1, err
    }

    op, err := findUsingPrefix(prefix)
    if err != nil {
        return -1, err
    }

    return op, nil

}

func GetAndFilter(phoneNumber string, f1 FilterPhoneFunc, f2 FilterOperatorFunc)(Operator,error){

	var (
		passFilterOne bool
		passFilterTwo bool
	)

	s, err := format(phoneNumber)
	if err != nil {
		return -1, err
	}

	if f1 != nil {
        passFilterOne = f1(s)
		if !passFilterOne{
			return -1, errors.New("could not pass set filters")
		}
    }
	op, err := Get(s)
	if err != nil{
		return -1, err
	}
	if f2 != nil{
		passFilterTwo = f2(op)
		if !passFilterTwo{
			return -1, errors.New("could not pass set filters")
		}
	}

	return op,nil
}

func Information(phoneNumber string)(*Info,error){

    fmtNumber, err := format(phoneNumber)
	if err != nil{
		return nil,err
	}
    op, err := Get(phoneNumber)
    if err != nil {
        return nil, err
    }


	info := &Info{
		Operator:        op,
		FormattedNumber: fmtNumber,
	}

    return info, nil

}

func InfoAfterFilters(phoneNumber string, f1 FilterPhoneFunc, f2 FilterOperatorFunc)(*Info,error){

	var (
		passFilterOne bool
		passFilterTwo bool
	)
	fmtNumber, err := format(phoneNumber)
	if err != nil{
		return nil,err
	}

	if f1 != nil {
		passFilterOne = f1(fmtNumber)
		if !passFilterOne{
			return nil, errors.New("could not pass set filters")
		}
	}
	op, err := Get(fmtNumber)
	if err != nil {
		return nil, err
	}

	if f2 != nil{
		passFilterTwo = f2(op)
		if !passFilterTwo{
			return nil, errors.New("could not pass set filters")
		}
	}
	info := &Info{
		Operator:        op,
		FormattedNumber: fmtNumber,
	}

	return info, nil
}
