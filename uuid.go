// Universally Unique IDentifier kit
//
// Currently parses standard UID's and can generate
// new V4 (random) UID's.
//
// Uses crypto/rand for seed.
package uuid

import "crypto/rand"
import "fmt"
import "encoding/json"
import "strings"

// A universally unique identifier.
type UUID [16]byte

// Constructs a new V4 (random) UUID.  NewV4 can panic
// if there is an error reading from the random source.
func NewV4() (u *UUID) {
	u, err := V4()
	if err != nil {
		panic(err)
	}
	return
}

// Constructs a new V4 (random) UUID.  Error is returned
// iff there is an error reading from the random source.
func V4() (u *UUID, err error) {
	u = new(UUID)
	_, err = rand.Read(u[0:16])
	if err != nil {
		return
	}
	// Set bits 6&7 of byte 8 to 0 and 1 respectively
	// 0x80
	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F
	return
}

// Formats a UUID as a standard UUID string.
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// Marshal a UUID to a UUID string so as to
// avoid byte-format marshalling.
func (u *UUID) MarshalJSON() (buff []byte, err error) {
	return json.Marshal(u.String())
}

// Parse an UUID string from JSON.
func (u *UUID) UnmarshalJSON(buff []byte) (err error) {
	ustr := ""
	err = json.Unmarshal(buff, &ustr)
	if err == nil {
		err = u.parse(ustr)
	}
	return
}

// Parse an UUID string and return a new object.
func Parse(s string) (u *UUID, err error) {
	u = new(UUID)
	u.parse(s)
	return
}

func hexValue(c byte) byte {
	switch {
	// `c` is a number
	case '0' <= c && c <= '9':
		return c - '0'
	// `c` is an char [a-f]
	case 'a' <= c && c <= 'f':
		return 10 + (c - 'a')
	// `c` is a char [A-F]
	case 'A' <= c && c <= 'F':
		return 10 + (c - 'A')
	}
	return 0xff
}

func (u *UUID) parse(s string) (err error) {
	if u == nil {
		u = new(UUID)
	}
	//fmt.Printf("UUID Unmarshal: [%s]\n", s)
	blks := strings.SplitN(s, "-", 5)
	hexstr := strings.Join(blks, "")
	var value byte = 0
	for bi := range hexstr {
		if bi%2 == 0 {
			value |= hexValue(hexstr[bi]) << 4
		} else {
			value |= hexValue(hexstr[bi])
			u[bi/2] = value
			value = 0
		}
	}
	//fmt.Printf("UUID Unmarshaled: [%s]\n", u.String())
	return
}
