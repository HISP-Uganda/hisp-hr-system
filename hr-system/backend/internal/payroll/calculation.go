package payroll

func CalculateAmounts(baseSalary, allowancesTotal, deductionsTotal, taxTotal float64) (grossPay float64, netPay float64) {
	grossPay = baseSalary + allowancesTotal
	netPay = grossPay - deductionsTotal - taxTotal
	return grossPay, netPay
}
