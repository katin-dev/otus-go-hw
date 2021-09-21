package hw02

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type state struct {
	sym      rune
	escaping bool
	builder  *strings.Builder
}

func Unpack(in string) (string, error) {
	var num int
	var err error

	state := state{}
	state.builder = &strings.Builder{}

	// Добавим в конец пробел, чтоб цикл отработал полностью. Сам пробел будет проигнорирован, как последний символ
	in += " "
	for _, s := range in {
		if unicode.IsDigit(s) {
			if num, err = strconv.Atoi(string(s)); err != nil {
				return "", err
			}

			if err := goNum(&state, num); err != nil {
				return "", err
			}
		} else if err := goSym(&state, s); err != nil {
			return "", err
		}
	}

	return state.builder.String(), nil
}

func goNum(st *state, num int) error {
	// Включен режим экранирования. Цифра считается буквой
	if st.escaping {
		st.escaping = false
		return goSym(st, '0'+int32(num))
	}

	// Число, перед которым нет буквы - ошибка
	if st.sym == 0 {
		return ErrInvalidString
	}

	// Множим текущую букву на x(num)
	st.builder.WriteString(
		strings.Repeat(string(st.sym), num),
	)

	st.sym = 0

	return nil
}

func goSym(st *state, sym rune) error {
	if sym == '\\' {
		if st.escaping {
			st.escaping = false
		} else {
			st.escaping = true
		}
	}

	// Если в памяти уже есть буква - сбрасываем её
	if st.sym != 0 {
		st.builder.WriteRune(st.sym)
	}

	if st.escaping {
		// Экранировать можно только сам слеш
		if sym != '\\' {
			return ErrInvalidString
		}

		// Слеш, включающий экранирование игнорируется
		sym = 0
	}

	st.sym = sym

	return nil
}
