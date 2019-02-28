package expr

import (
	"testing"
)

func Test_View(t *testing.T) {
	exp, err := Translate("(+ a (/ (+ (/ b c) (/ (+ d e) f)) (/ (/ g h) i)))")
	if err != nil {
		t.Fatalf("should translate:\n%s", exp)
	}
	prettystr := Pretty(exp)
	t.Error(exp.String())
	t.Error(prettystr)

	exp, err = Translate("(/ (/ (+ a b) c) (/ d (* e f)))")
	if err != nil {
		t.Fatalf("should translate:\n%s", exp)
	}
	prettystr = Pretty(exp)
	t.Error(exp.String())
	t.Error(prettystr)

}

func Test_vertspace(t *testing.T) {
	exp, err := Translate("(+ a (/ (+ e (/ (/ f g) (/ h i))) (/ b c)))")
	if err != nil {
		t.Fatalf("should translate:\n%s", exp)
	}
	res := vertspace(exp)
	req := 11
	if res != req {
		t.Errorf("%s\nshould return %d but is %d", exp, req, res)
	}
}

func Test_downspace(t *testing.T) {
	exp, err := Translate("(+ a (/ (/ b c) (+ e (/ (/ f g) (/ h i)))))")
	if err != nil {
		t.Fatalf("should translate:\n%s", exp)
	}
	res := downspace(exp)
	req := 7
	if res != req {
		t.Errorf("%s\nshould return %d but is %d", exp, req, res)
	}
}

func Test_upspace(t *testing.T) {
	exp, err := Translate("(+ a (/ (+ e (/ (/ f g) (/ h i))) (/ b c)))")
	if err != nil {
		t.Fatalf("should translate:\n%s", exp)
	}
	res := upspace(exp)
	req := 7
	if res != req {
		t.Errorf("%s\nshould return %d but is %d", exp, req, res)
	}
}
