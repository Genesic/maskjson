package maskjson

import (
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

type Profile struct {
	Name  string
	Phone string      `mask:"true"`
	Email string      `json:"email"`
	Auth  Auth        `json:"auth,omitempty"`
	Token interface{} `json:"token"`
}

func TestMask_Marshal(t *testing.T) {
	m := NewMask(false, 3)

	for _, tt := range []Profile{
		{
			Name:  "henry",
			Phone: "1234567890",
			Email: "9h4Hb@example.com",
		},
		{
			Name:  "henry",
			Phone: "1234567890",
			Email: "9h4Hb@example.com",
			Auth: Auth{
				Account:  "henry",
				Password: "123456",
			},
		},
		{
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
	} {
		b, err := m.Marshal(tt)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
	}
}
