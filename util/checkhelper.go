package util

import "fmt"

func EqualInt(v1 int, v2 int, txt string) error {
	if v1 == v2 {
		return nil
	}
	return fmt.Errorf("%s check %d equal %d fail", txt, v1, v2)
}

func InEqualInt(v1 int, v2 int, txt string) error {
	if v1 != v2 {
		return nil
	}
	return fmt.Errorf("%s check %d inequal %d fail", txt, v1, v2)
}

func GreatInt(v1 int, v2 int, txt string) error {
	if v1 > v2 {
		return nil
	}
	return fmt.Errorf("%s check %d great %d fail", txt, v1, v2)
}

func LittleInt(v1 int, v2 int, txt string) error {
	if v1 < v2 {
		return nil
	}
	return fmt.Errorf("%s check %d little %d", txt, v1, v2)
}

func GreatEqualInt(v1 int, v2 int, txt string) error {
	if v1 >= v2 {
		return nil
	}
	return fmt.Errorf("%s check %d greatequal %d fail", txt, v1, v2)
}

func LittleEqualInt(v1 int, v2 int, txt string) error {
	if v1 <= v2 {
		return nil
	}
	return fmt.Errorf("%s check %d littleequal %d", txt, v1, v2)
}

func EqualStr(v1 string, v2 string, txt string) error {
	if v1 == v2 {
		return nil
	}
	return fmt.Errorf("%s check %s equal %s fail", txt, v1, v2)
}

func InEqualStr(v1 string, v2 string, txt string) error {
	if v1 != v2 {
		return nil
	}
	return fmt.Errorf("%s check %s inequal %s fail", txt, v1, v2)
}
