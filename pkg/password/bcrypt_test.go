package password

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/taraslis453/solid-software-test/pkg/logging"
)

func Test_bcrypt_GenerateHashFromPassword(t *testing.T) {
	logger := logging.NewZapLogger("error")
	passwordHasher := NewBcrypt(logger)

	hashedPassword, err := passwordHasher.GenerateHashFromPassword("password")
	require.NoError(t, err, "failed to generate hashed password")
	require.NotEmpty(t, hashedPassword, "hashed password is empty")
}

func Test_bcrypt_CompareHashAndPassword(t *testing.T) {
	logger := logging.NewZapLogger("error")
	passwordHasher := NewBcrypt(logger)

	type args struct {
		password        string
		comparePassword string
	}
	testCases := []struct {
		name                        string
		args                        args
		expectedPasswordCorrectness bool
	}{
		{
			name: "positive:generate",
			args: args{
				password:        "password",
				comparePassword: "password",
			},
			expectedPasswordCorrectness: true,
		},
		{
			name: "negative:invalid password",
			args: args{
				password:        "password",
				comparePassword: "wrong",
			},
			expectedPasswordCorrectness: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := passwordHasher.GenerateHashFromPassword(tc.args.password)
			require.NoError(t, err, "failed to generate hashed password")
			require.NotEmpty(t, hashedPassword, "hashed password is empty")

			isPasswordCorrect, err := passwordHasher.CompareHashAndPassword(&CompareHashAndPasswordOptions{
				Hashed:   hashedPassword,
				Password: tc.args.comparePassword,
			})
			require.Equal(t, tc.expectedPasswordCorrectness, isPasswordCorrect)
			require.NoError(t, err, "failed to check password correctness")
		})
	}
}
