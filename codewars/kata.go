package codewars

func TwoOldestAges(ages []int) [2]int {
  a, b := 0, 0
  for _, v := range ages {
    if v > b {
      a, b = b, v
    } else if v > a {
      a = v
    }
  }
  return [2]int{a, b}
}