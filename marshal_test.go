package maskjson

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Auth struct {
	Account  string `json:"account"`
	Password string `json:"password" mask:"true"`
}

type Token struct {
	Token       string `json:"token,omitempty" mask:"true"`
	ExpiresTime time.Duration
}

type Secret struct {
	Key string `json:"key" mask:"true"`
	IV  string `json:"iv" mask:"true"`
}

type Profile struct {
	Name   string
	Phone  string      `mask:"true"`
	Email  string      `json:"email"`
	Auth   Auth        `json:"auth,omitempty"`
	Token  interface{} `json:"token,omitempty"`
	Secret *Secret     `json:"secret,omitempty"`
}

func TestMask_Marshal(t *testing.T) {
	m := NewMask(false, 3)

	type testCase struct {
		desc    string
		profile Profile
		result  string
	}

	for _, tt := range []testCase{
		{
			desc: "test 1",
			profile: Profile{
				Name:  "henry",
				Phone: "1234567890",
				Email: "9h4Hb@example.com",
			},
			result: `{"Name":"henry","Phone":"1234******","auth":{"account":"","password":""},"email":"9h4Hb@example.com"}`,
		},
		{
			desc: "test2",
			profile: Profile{
				Name:  "henry",
				Phone: "1234567890",
				Email: "9h4Hb@example.com",
				Auth: Auth{
					Account:  "henry",
					Password: "123456",
				},
			},
			result: `{"Name":"henry","Phone":"1234******","auth":{"account":"henry","password":"12****"},"email":"9h4Hb@example.com"}`,
		},
		{
			desc: "test 4",
			profile: Profile{
				Name:  "henry",
				Phone: "1234567890",
				Email: "9h4Hb@example.com",
				Auth: Auth{
					Account:  "henry",
					Password: "123456",
				},
				Token: Token{
					Token:       "123456",
					ExpiresTime: time.Second * 10,
				},
			},
			result: `{"Name":"henry","Phone":"1234******","auth":{"account":"henry","password":"12****"},"email":"9h4Hb@example.com","token":{"ExpiresTime":10000000000,"token":"12****"}}`,
		},
		{
			desc: "test 4",
			profile: Profile{
				Name:  "henry",
				Phone: "1234567890",
				Email: "9h4Hb@example.com",
				Auth: Auth{
					Account:  "henry",
					Password: "abcd1234",
				},
				Secret: &Secret{
					Key: "my_secret_key",
					IV:  "my_secret_iv",
				},
			},
			result: `{"Name":"henry","Phone":"1234******","auth":{"account":"henry","password":"abc*****"},"email":"9h4Hb@example.com","secret":{"iv":"my_s********","key":"my_se********"}}`,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := m.Marshal(tt.profile)
			if err != nil {
				t.Fatal(err)
			}
			if string(result) != tt.result {
				t.Errorf("want %s\ngot %s", tt.result, string(result))
			}

			ori, err := json.Marshal(tt.profile)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("origin:", string(ori))
			fmt.Println("masked:", string(result))
		})
	}
}
