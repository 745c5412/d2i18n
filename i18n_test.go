package d2i18n

import (
	"os"
	"testing"
)

func openFixture(t *testing.T) *os.File {
	file, err := os.Open("fixtures/i18n_fr.d2i")
	if err != nil {
		t.Fatalf("openFixture failed : %v", err)
	}
	return file
}

func TestParse(t *testing.T) {
	reader := openFixture(t)
	defer func() { _ = reader.Close() }()

	_, err := Parse(NewReader(reader))
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
}

func Test_i18n_GetNamedText(t *testing.T) {
	reader := openFixture(t)
	defer func() { _ = reader.Close() }()

	i18n, err := Parse(NewReader(reader))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			"valid",
			args{"ui.chat.console.noHelp"},
			"Aucune aide n'est disponible pour la commande '%1'.",
			true,
			false,
		},
		{
			"unknown",
			args{"unnqnzjdqkzd"},
			"",
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := i18n.GetNamedText(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("i18n.GetNamedText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("i18n.GetNamedText() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("i18n.GetNamedText() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
