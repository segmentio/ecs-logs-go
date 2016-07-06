package ecslogs

import "testing"

var levelTests = []struct {
	lvl   Level
	str   string
	gostr string
}{
	{
		lvl:   EMERG,
		str:   "EMERG",
		gostr: "Level(0)",
	},
	{
		lvl:   ALERT,
		str:   "ALERT",
		gostr: "Level(1)",
	},
	{
		lvl:   CRIT,
		str:   "CRIT",
		gostr: "Level(2)",
	},
	{
		lvl:   ERROR,
		str:   "ERROR",
		gostr: "Level(3)",
	},
	{
		lvl:   WARN,
		str:   "WARN",
		gostr: "Level(4)",
	},
	{
		lvl:   NOTICE,
		str:   "NOTICE",
		gostr: "Level(5)",
	},
	{
		lvl:   INFO,
		str:   "INFO",
		gostr: "Level(6)",
	},
	{
		lvl:   DEBUG,
		str:   "DEBUG",
		gostr: "Level(7)",
	},
}

func TestParseLevelSuccess(t *testing.T) {
	for _, test := range levelTests {
		if lvl, err := ParseLevel(test.str); err != nil {
			t.Errorf("%s: error: %s", test.str, err)
		} else if lvl != test.lvl {
			t.Errorf("%s: invalid level: %s", test.str, lvl)
		}
	}
}

func TestParseLevelFailure(t *testing.T) {
	if _, err := ParseLevel(""); err == nil {
		t.Error("no error returned when parsing an invalid log level")
	} else if s := err.Error(); s != "invalid message level \"\"" {
		t.Error("invalid error message returned when parsing an invalid log level:", s)
	}
}

func TestLevelString(t *testing.T) {
	for _, test := range levelTests {
		if s := test.lvl.String(); s != test.str {
			t.Errorf("%s: invalid string: %s", test.lvl, s)
		}
	}
}

func TestLevelGoString(t *testing.T) {
	for _, test := range levelTests {
		if s := test.lvl.GoString(); s != test.gostr {
			t.Errorf("%s: invalid go string: %s", test.lvl, s)
		}
	}
}

func TestLevelJSON(t *testing.T) {
	for _, test := range levelTests {
		if b, err := test.lvl.MarshalJSON(); err != nil {
			t.Errorf("%s: %s", test.lvl, err)
		} else {
			var lvl Level
			if err = lvl.UnmarshalJSON(b); err != nil {
				t.Errorf("%s: %s", lvl, err)
			}
		}
	}
}

func TestLevelText(t *testing.T) {
	for _, test := range levelTests {
		if b, err := test.lvl.MarshalText(); err != nil {
			t.Errorf("%s: %s", test.lvl, err)
		} else {
			var lvl Level
			if err = lvl.UnmarshalText(b); err != nil {
				t.Errorf("%s: %s", lvl, err)
			}
		}
	}
}
