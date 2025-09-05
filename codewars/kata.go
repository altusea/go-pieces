package codewars

import "net"

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


func Is_valid_ip(ip string) bool {
	res := net.ParseIP(ip)
	return res != nil
}

func Factorial(n int) int {
    // Handle base cases
    if n < 0 {
        return 0 // Factorial is not defined for negative numbers
    }
    if n == 0 || n == 1 {
        return 1
    }
    
    // Calculate factorial iteratively to avoid stack overflow for large n
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    
    return result
}