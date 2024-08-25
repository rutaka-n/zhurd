package label

import (
	"bytes"
	"errors"
	"testing"
)

func TestPrint(t *testing.T) {
	tcs := []struct {
		desc         string
		tplt         []byte
		placeholders map[string]string
		expected     []byte
		expectedErr  error
	}{
		{
			desc: "label w/o placeholders",
			tplt: []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`),
			placeholders: map[string]string{},
			expected: []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`),
			expectedErr: nil,
		},
		{
			desc: "label with placeholders",
			tplt: []byte(`^XA
^FX _comment1_
^CF_font1_
^FO_identation_left_top1_^GB_graphical_box1_^FS
^FO_identation_left_top2_^FR^GB_graphical_box2_^FS
^FO_identation_left_top3_^GB_graphical_box3_^FS
^FO_identation_left_top4_^FD_company_name_^FS
^CF_font2_
^FO_identation_left_mid1_^FD_company_address_building_^FS
^FO_identation_left_mid2_^FD_company_address_city_ _company_address_postcode_^FS
^FO_identation_left_mid3_^FD_company_address_country_^FS
^FO_identation_left_mid4_^GB_graphical_box4_^FS
^FX _comment2_
^BY_global_bar_code_
^FO_identation_left_botttom1_^BC^FD_field_data_string_^FS
^XZ
`),
			placeholders: map[string]string{
				"_comment1_":                 "Top section with logo, name and address.",
				"_font1_":                    "0,60",
				"_identation_left_top1_":     "50,50",
				"_graphical_box1_":           "100,100,100",
				"_identation_left_top2_":     "75,75",
				"_graphical_box2_":           "100,100,100",
				"_identation_left_top3_":     "93,93",
				"_graphical_box3_":           "40,40,40",
				"_identation_left_top4_":     "220,50",
				"_company_name_":             "Intershipping, Inc.",
				"_font2_":                    "0,30",
				"_identation_left_mid1_":     "220,115",
				"_company_address_building_": "1000 Shipping Lane",
				"_identation_left_mid2_":     "220,155",
				"_company_address_city_":     "Shelbyville TN",
				"_company_address_postcode_": "38102",
				"_company_address_country_":  "United States (USA)",
				"_identation_left_mid3_":     "220,195",
				"_identation_left_mid4_":     "50,250",
				"_graphical_box4_":           "700,3,3",
				"_comment2_":                 "Section with bar code.",
				"_global_bar_code_":          "5,2,270",
				"_field_data_string_":        "12345678",
				"_identation_left_botttom1_": "100,550",
			},
			expected: []byte(`^XA
^FX Top section with logo, name and address.
^CF0,60
^FO50,50^GB100,100,100^FS
^FO75,75^FR^GB100,100,100^FS
^FO93,93^GB40,40,40^FS
^FO220,50^FDIntershipping, Inc.^FS
^CF0,30
^FO220,115^FD1000 Shipping Lane^FS
^FO220,155^FDShelbyville TN 38102^FS
^FO220,195^FDUnited States (USA)^FS
^FO50,250^GB700,3,3^FS
^FX Section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`),
			expectedErr: nil,
		},
		{
			desc: "missing placeholder",
			tplt: []byte(`^XA
^FX _comment1_
^BY_global_bar_code_
^FO_identation_^BC^FD_field_data_string_^FS
^XZ
`),
			placeholders: map[string]string{
				"_comment1_":          "bar code section",
				"_global_bar_code_":   "5,2,270",
				"_field_data_string_": "12345678",
			},
			expected:    nil,
			expectedErr: MissingPlaceholderError,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			tplt, err := NewTemplate("ZPL", tc.tplt)
			if err != nil {
				t.Fatalf("got err: %v\n", err)
			}

			result, err := tplt.Print(tc.placeholders)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected err: %v, got err: %v\n", tc.expectedErr, err)
			}
			if !bytes.Equal(tc.expected, result) {
				t.Errorf("expected: %s, got: %s\n", tc.expected, result)
			}
		})
	}
}
