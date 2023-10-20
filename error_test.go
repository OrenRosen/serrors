package serrors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/FTBpro/OrenRosen/serrors"
)

func TestSerror(t *testing.T) {
	testCases := []struct {
		desc            string
		errFunc         func() error
		expectedKeyvals []any
		expectedMessage string
	}{
		{
			desc:            "no error",
			errFunc:         func() error { return nil },
			expectedKeyvals: nil,
			expectedMessage: "",
		},
		{
			desc: "new error",
			errFunc: func() error {
				return serrors.New("new error",
					"key1", "val1",
					"key2", "val2",
				)
			},
			expectedKeyvals: []any{"key1", "val1", "key2", "val2"},
			expectedMessage: "new error",
		},
		{
			desc: "new and wrap error",
			errFunc: func() error {
				err := serrors.New("new error",
					"key1", "val1",
					"key2", "val2",
				)
				return serrors.Wrap(err, "wrap msg",
					"key3", "val3",
				)
			},
			expectedKeyvals: []any{"key1", "val1", "key2", "val2", "key3", "val3"},
			expectedMessage: "wrap msg: new error",
		},
		{
			desc: "with standard error in between",
			errFunc: func() error {
				err := serrors.New("new error",
					"key1", "val1",
					"key2", "val2",
				)

				err = fmt.Errorf("fmt error: %w", err)

				return serrors.Wrap(err, "wrap msg",
					"key3", "val3",
				)
			},
			expectedKeyvals: []any{"key1", "val1", "key2", "val2", "key3", "val3"},
			expectedMessage: "wrap msg: fmt error: new error",
		},
		{
			desc: "duplicate keys",
			errFunc: func() error {
				err := serrors.New("new error",
					"key1", "val1",
					"key2", "val2",
				)
				return serrors.Wrap(err, "wrap msg",
					"key1", "other-val1",
					"key3", "val3",
				)
			},
			expectedKeyvals: []any{"key1", "val1", "key2", "val2", "key3", "val3", "key1", "other-val1"},
			expectedMessage: "wrap msg: new error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.errFunc()
			keyvals := serrors.KeyVals(err)
			if tc.expectedMessage == "" {
				require.Nil(t, err)
			} else {
				require.Equal(t, tc.expectedMessage, err.Error())
			}

			if tc.expectedKeyvals == nil {
				require.Nil(t, keyvals)
			} else {
				require.Len(t, keyvals, len(tc.expectedKeyvals))
				require.ElementsMatch(t, keyvals, tc.expectedKeyvals)
			}
		})
	}
}
