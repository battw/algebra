package expr

import ()

type runegrid struct {
	grd    [][]rune
	exp    *Expr
	starti int
}

func newrunegrid(e *Expr) *runegrid {
	height := vertspace(e)
	length := len(e.String())
	rg := runegrid{grd: make([][]rune, height), exp: e, starti: upspace(e)}
	for i, _ := range rg.grd {
		rg.grd[i] = make([]rune, length)
		for j, _ := range rg.grd[i] {
			rg.grd[i][j] = ' '
		}
	}
	return &rg
}

func Pretty(e *Expr) string {
	rg := newrunegrid(e)
	var rec func(y, x int, e *Expr) int
	rec = func(y, x int, e *Expr) int {
		switch e.typ {
		case OP:
			rg.grd[y][x] = '('
			x++
			rg.grd[y][x] = e.sym
			x++
			var xat int //current x position (returned by rec)
			switch e.sym {
			case '/':
				xat = max(rec(y-downspace(e.l)-1, x, e.l),
					rec(y+upspace(e.r)+1, x, e.r))
				rg.grd[y][xat] = ')'
				return xat + 1
			default:
				xat = rec(y, x+1, e.l)
				xat = rec(y, xat+1, e.r)
				rg.grd[y][xat] = ')'
				return xat + 1
			}

		default:
			rg.grd[y][x] = e.sym
			return x + 1
		}
	}
	xend := rec(rg.starti, 0, e)
	prettystr := ""
	for _, line := range rg.grd {
		prettystr += string(line[:xend]) + "\n"
	}
	//Draw horizontal lines separating numerators and denominators
	newstr := ""
	drawing := false
	for _, r := range prettystr {
		drawing = drawing && r != ')'
		if r == '/' {
			drawing = true
		}
		if drawing {
			newstr += string('-')
		} else {
			newstr += string(r)
		}
	}
	prettystr = newstr

	return prettystr
}

//vertspace - the vertical space taken up by the pretty print version of the expression.
func vertspace(e *Expr) int {
	switch e.typ {
	case OP:
		switch e.sym {
		case '/':
			return 1 + vertspace(e.l) + vertspace(e.r)
		default:
			return max(vertspace(e.l), vertspace(e.r))
		}
	default:
		return 1
	}
}

//downspace - the space taken by the pretty print version below the first operator.
func downspace(e *Expr) int {
	switch e.typ {
	case OP:
		switch e.sym {
		case '/':
			return vertspace(e.r)
		default:
			return max(downspace(e.l), downspace(e.r))
		}
	default:
		return 0
	}

}

//upspace - the space taken by the pretty print version above the first operator.
func upspace(e *Expr) int {
	switch e.typ {
	case OP:
		switch e.sym {
		case '/':
			return vertspace(e.l)
		default:
			return max(upspace(e.l), upspace(e.r))
		}
	default:
		return 0
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
