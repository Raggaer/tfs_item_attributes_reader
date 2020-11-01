package tfsreader

import (
	"io/ioutil"
	"testing"
)

var hexItems = []string{
	"220100000000000000110072756E65656D626C656D63686172676573028813000000000000",
	"1837007265666F72676564206D6173746572637261667465642068656C6D6574202B33302041544B2C202B3330205245532C202B333020534B4C1C000000002204000000000000000A00736B696C6C76616C75650241000000000000000D00707670726573697374616E6365020000000000000000070072656772616465011E007B5B226C6576656C225D203D20342C5B2267656D73225D203D207B7D2C7D0B006D6167696361747461636B020000000000000000",
	"160100181C0061646F726E6564207370656369616C69737420656D626C656D202B3822020000000000000005006C6576656C020800000000000000110072756E65656D626C656D63686172676573028813000000000000",
	"182900696365666F726765642061726D6F72202B31352041544B2C202B3135205245532C202B313520534B4C1C000000002204000000000000000B006D6167696361747461636B020000000000000000070072656772616465011E007B5B226C6576656C225D203D20342C5B2267656D73225D203D207B7D2C7D0D00707670726573697374616E63650200000000000000000A00736B696C6C76616C7565021B00000000000000",
	"182A00696365666F726765642068656C6D6574202B31352041544B2C202B3135205245532C202B313520534B4C1C000000002204000000000000000B006D6167696361747461636B020000000000000000070072656772616465011E007B5B226C6576656C225D203D20342C5B2267656D73225D203D207B7D2C7D0D00707670726573697374616E63650200000000000000000A00736B696C6C76616C7565021B00000000000000",
	"160100181900676F6C64207370656369616C69737420656D626C656D202B38220200000000000000110072756E65656D626C656D6368617267657302881300000000000005006C6576656C020800000000000000",
	"18190073616D7572616920736F756C2072756E652028543129202B3822010000000000000005006C6576656C020800000000000000",
	"0F1E241E00",
	"0F14241400",
}

func TestUnserializeItemHex(t *testing.T) {
	for _, i := range hexItems {
		_, err := UnserializeHexString(i)
		if err != nil {
			t.Fatalf("Unable to unserialize item attributes - %v", err)
			return
		}
	}
}

func TestUnserializeCountHex1(t *testing.T) {
	c := "0F05"
	item, err := UnserializeHexString(c)
	if err != nil {
		t.Fatalf("Unable to unserialize item attributes - %v", err)
		return
	}

	if item.Count != 5 {
		t.Fatalf("Invalid item count, expected 5 but got %d", item.Count)
		return
	}
}

func TestUnserializeTextHex1(t *testing.T) {
	c := "061D004B6E656B726F206D616E64612079206E6F2074752070616E64610A73641279BB985F130600416C7661726F"
	item, err := UnserializeHexString(c)
	if err != nil {
		t.Fatalf("Unable to unserialize item attributes - %v", err)
		return
	}

	if item.WrittenBy != "Alvaro" {
		t.Fatalf("Invalid written by value, expected %s but got %s", "Alvaro", item.WrittenBy)
		return
	}
	if item.Text != "Knekro manda y no tu panda\nsd" {
		t.Fatalf("Invalid text value, expected %s but got %s", "Knekro manda y no tu panda\nsd", item.Text)
		return
	}
}

func TestUnserialize(t *testing.T) {
	data, err := ioutil.ReadFile("test/player_items-1772_103_112_2304_30_-----.attributes")
	if err != nil {
		t.Fatalf("Unable to open testing file - %v", err)
		return
	}

	_, err = Unserialize(data)
	if err != nil {
		t.Fatalf("Unable to unserialize item attributes - %v", err)
		return
	}
}
