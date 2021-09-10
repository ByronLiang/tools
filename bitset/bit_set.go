package bitset

import "fmt"

type BitUnit uint8

func AddBit(origin, new BitUnit) BitUnit {
    return origin | new
}

func Exist(origin, compare BitUnit) bool {
    return origin & compare == compare
}

func DelBit(origin, bit BitUnit) BitUnit {
    return origin & (^bit)
}

func BitString(bit BitUnit) string {
    return fmt.Sprintf("%b", bit)
}
