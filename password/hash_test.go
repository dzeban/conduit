package password

import (
	"testing"
)

func TestEncode(t *testing.T) {
	defaultParams := HashParams{
		Salt:       []byte("123"),
		Iterations: DefaultIterations,
		Memory:     DefaultMemory,
		Threads:    DefaultThreads,
		Len:        DefaultLen,
	}

	type testCase struct {
		expected string
		hash     []byte
		params   HashParams
	}

	tests := []testCase{
		{
			expected: "$argon2id$v=19$m=32768,t=5,p=1$MTIz$YWJj",
			hash:     []byte("abc"),
			params:   defaultParams,
		},
		{
			expected: "$argon2id$v=19$m=32768,t=5,p=1$b+s1YZKqMW6NOxFUt4y6Og$BuqHsbcLLZgMT9JmwMAqRioXa6GEJYTsmIjQFIgBcnh255qHdwcExJm0LAzeuPpk/qJZKekibmBElZhKJ/StSQ",
			hash: []byte{
				0x6, 0xea, 0x87, 0xb1, 0xb7, 0xb, 0x2d, 0x98,
				0xc, 0x4f, 0xd2, 0x66, 0xc0, 0xc0, 0x2a, 0x46,
				0x2a, 0x17, 0x6b, 0xa1, 0x84, 0x25, 0x84, 0xec,
				0x98, 0x88, 0xd0, 0x14, 0x88, 0x1, 0x72, 0x78,
				0x76, 0xe7, 0x9a, 0x87, 0x77, 0x7, 0x4, 0xc4, 0x99,
				0xb4, 0x2c, 0xc, 0xde, 0xb8, 0xfa, 0x64, 0xfe,
				0xa2, 0x59, 0x29, 0xe9, 0x22, 0x6e, 0x60, 0x44,
				0x95, 0x98, 0x4a, 0x27, 0xf4, 0xad, 0x49,
			},
			params: HashParams{
				Salt:       []byte{0x6f, 0xeb, 0x35, 0x61, 0x92, 0xaa, 0x31, 0x6e, 0x8d, 0x3b, 0x11, 0x54, 0xb7, 0x8c, 0xba, 0x3a},
				Iterations: DefaultIterations,
				Memory:     DefaultMemory,
				Threads:    DefaultThreads,
				Len:        DefaultLen,
			},
		},
	}

	for _, test := range tests {
		actual := Encode(test.hash, test.params)
		if actual != test.expected {
			t.Fatalf("encoded hashes dont match, expected %q, got %q\n", test.expected, actual)
		}
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		result   bool
		hasErr   bool
	}{
		{
			"Equals",
			"12345678",
			"$argon2id$v=19$m=32768,t=5,p=1$fD9/37P2wxyuXy8Uc5HsHg$daeNeJMrnjBjrwHiiIzU9GuXlEu6Ov0qyeqmpJIu690Y/1ElHeHjyJcoNCyXRiqxuqLM5PAcQEhJiweiAv7dow",
			true,
			false,
		},
		{
			"Differs",
			"1234567890",
			"$argon2id$v=19$m=32768,t=5,p=1$fD9/37P2wxyuXy8Uc5HsHg$daeNeJMrnjBjrwHiiIzU9GuXlEu6Ov0qyeqmpJIu690Y/1ElHeHjyJcoNCyXRiqxuqLM5PAcQEhJiweiAv7dow",
			false,
			false,
		},
		{
			"HashParamsDiffer",
			"12345678",
			"$argon2id$v=19$m=32768,t=2,p=3$fD9/37P2wxyuXy8Uc5HsHg$daeNeJMrnjBjrwHiiIzU9GuXlEu6Ov0qyeqmpJIu690Y/1ElHeHjyJcoNCyXRiqxuqLM5PAcQEhJiweiAv7dow",
			false,
			false,
		},
		{
			"SaltDiffer",
			"12345678",
			"$argon2id$v=19$m=32768,t=5,p=1$XXXXXXXXXXXXXXXXXXXXXX$daeNeJMrnjBjrwHiiIzU9GuXlEu6Ov0qyeqmpJIu690Y/1ElHeHjyJcoNCyXRiqxuqLM5PAcQEhJiweiAv7dow",
			false,
			false,
		},
	}

	for _, test := range tests {
		equals, err := Check(test.password, test.hash)
		hasErr := (err != nil)
		if test.hasErr != hasErr {
			t.Fatalf("%s test failed: unexpected error, want %#v, got %#v (err '%#v')\n", test.name, test.hasErr, hasErr, err)
		}

		if equals != test.result {
			t.Fatalf("%s test failed: unexpected result, want %#v, got %#v\n", test.name, test.result, equals)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		hash string
		err  error
	}{
		{
			"Empty",
			"",
			ErrHashInvalid,
		},
		{
			"Valid",
			"$argon2id$v=19$m=32768,t=5,p=1$62jx8XXcWFckh+9QjSkefA$Rav8qgyBopeMfDaBbk5jHs+UGUBPqcRqfxNNlR9RqBlV6ZrC/t8/05BYCwT6KKEbI5H4lOuqSLkyxClXm1WA4w",
			nil,
		},
		{
			"InvalidType",
			"$2b$v=19$m=32768,t=5,p=1$62jx8XXcWFckh+9QjSkefA$Rav8qgyBopeMfDaBbk5jHs+UGUBPqcRqfxNNlR9RqBlV6ZrC/t8/05BYCwT6KKEbI5H4lOuqSLkyxClXm1WA4w",
			ErrHashType,
		},
		{
			"InvalidVersion",
			"$argon2id$v=2$m=32768,t=5,p=1$62jx8XXcWFckh+9QjSkefA$Rav8qgyBopeMfDaBbk5jHs+UGUBPqcRqfxNNlR9RqBlV6ZrC/t8/05BYCwT6KKEbI5H4lOuqSLkyxClXm1WA4w",
			ErrHashVersion,
		},
		{
			"InvalidParams",
			"$argon2id$v=19$ohai$62jx8XXcWFckh+9QjSkefA$Rav8qgyBopeMfDaBbk5jHs+UGUBPqcRqfxNNlR9RqBlV6ZrC/t8/05BYCwT6KKEbI5H4lOuqSLkyxClXm1WA4w",
			ErrHashParams,
		},
		{
			"InvalidParamNames",
			"$argon2id$v=19$a=32768,b=5,c=1$62jx8XXcWFckh+9QjSkefA$Rav8qgyBopeMfDaBbk5jHs+UGUBPqcRqfxNNlR9RqBlV6ZrC/t8/05BYCwT6KKEbI5H4lOuqSLkyxClXm1WA4w",
			ErrHashParams,
		},
		{
			"InvalidSalt",
			"$argon2id$v=19$m=32768,t=5,p=1$lol-kek$Rav8qgyBopeMfDaBbk5jHs+UGUBPqcRqfxNNlR9RqBlV6ZrC/t8/05BYCwT6KKEbI5H4lOuqSLkyxClXm1WA4w",
			ErrHashSalt,
		},
		{
			"InvalidHash",
			"$argon2id$v=19$m=32768,t=5,p=1$62jx8XXcWFckh+9QjSkefA$ololololo",
			ErrHash,
		},
	}

	for _, test := range tests {
		_, _, err := Decode(test.hash)
		if test.err != err {
			t.Fatalf("%s test failed: unexpected error, want %#v, got %#v\n", test.name, test.err, err)
		}
	}
}
