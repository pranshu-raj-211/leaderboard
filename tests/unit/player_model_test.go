package unit_test

import (
	"errors"
	"leaderboard/src/models"
	"strings"
	"testing"
)

func TestPlayerValidate(t *testing.T) {
	cases := []struct {
		name         string
		input        models.Player
		wantErr      bool
		wantErrField string
	}{
		{
			name:    "valid player",
			input:   models.Player{PlayerID: "p_0001", Name: "SwiftFalcon"},
			wantErr: false,
		},
		{
			name:         "empty player id",
			input:        models.Player{PlayerID: "", Name: "SwiftFalcon"},
			wantErr:      true,
			wantErrField: "player_id",
		},
		{
			name:         "empty name",
			input:        models.Player{PlayerID: "p_0001", Name: ""},
			wantErr:      true,
			wantErrField: "name",
		},
		{
			name:         "name too long",
			input:        models.Player{PlayerID: "p_0001", Name: strings.Repeat("a", 65)},
			wantErr:      true,
			wantErrField: "name",
		},
		{
			name:    "name at max length",
			input:   models.Player{PlayerID: "p_0001", Name: strings.Repeat("a", 64)},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Validate()

			if !tc.wantErr {
				if got != nil {
					t.Fatalf("unexpected validation error: %v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected error, got none")
			}
			var ve *models.ValidationError
			if !errors.As(got, &ve) {
				t.Fatalf("expected *models.ValidationError, got %T", got)
			}
			if ve.Field != tc.wantErrField {
				t.Fatalf("expected error field %q, got %q", tc.wantErrField, ve.Field)
			}
		})
	}
}
