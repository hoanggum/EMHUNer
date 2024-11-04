package utility

import (
	"EMHUNer/models"
	"fmt"
)

func CalculateTransactionUtility(transaction *models.Transaction) float64 {
	totalUtility := 0.0
	for _, utility := range transaction.Utilities {
		totalUtility += utility
	}
	return totalUtility
}

func CalculateAndPrintAllTransactionUtilities(transactions []*models.Transaction) {
	for i, transaction := range transactions {
		tu := CalculateTransactionUtility(transaction)
		fmt.Printf("Transaction %d TU: %.2f\n", i+1, tu)
	}
}

func CalculateRLUForAllItemsRhoAnDenta(transactions []*models.Transaction, rho, delta map[int]bool, utilityArray *models.UtilityArray) {
	combinedSet := UnionMaps(rho, delta)

	for item := range combinedSet {
		totalRLU := 0.0
		fmt.Printf("\nCalculating RLU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsItem(transaction, item) {
				fmt.Printf("  Found item %d in transaction: %v\n", item, transaction.Items)
				rlu := CalculateRemainingResidualUtility(transaction, item)
				totalRLU += rlu
				fmt.Printf("  RLU for this transaction: %.2f (cumulative RLU: %.2f)\n", rlu, totalRLU)
			}
		}

		utilityArray.SetRLU(item, totalRLU)
		fmt.Printf("Calculated total RLU for item %d: %.2f\n", item, totalRLU)
	}
}

func CalculateRLUForAllItems(transactions []*models.Transaction, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRLU := 0.0
		fmt.Printf("\nCalculating RLU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsItem(transaction, item) {
				index := GetItemIndex(transaction, item)
				itemUtility := transaction.Utilities[index]
				remainingUtility := CalculateRemainingUtility(transaction, index+1)
				totalRLU += itemUtility + remainingUtility

				fmt.Printf("  Found item %d in transaction %v with utility: %.2f, Remaining Residual Utility: %.2f\n",
					item, transaction.Items, itemUtility, remainingUtility)
			}
		}

		utilityArray.SetRLU(item, totalRLU)
		fmt.Printf("Calculated total RLU for item %d: %.2f\n", item, totalRLU)
	}
}

func CalculateRemainingResidualUtility(transaction *models.Transaction, currentItem int) float64 {
	foundCurrentItem := false
	rru := 0.0
	fmt.Printf("    Remaining items after %d: ", currentItem)

	for i, item := range transaction.Items {
		utility := transaction.Utilities[i]

		if foundCurrentItem && utility > 0 {
			rru += utility
			fmt.Printf("%d(%.2f) ", item, utility)
		}

		if item == currentItem {
			foundCurrentItem = true
			if utility > 0 {
				rru += utility
				fmt.Printf("    Adding utility of currentItem %d: %.2f\n", currentItem, utility)
			}
		}
	}
	fmt.Println()
	return rru
}

//Hàm cũ
// func CalculateRTWUForAllItems(transactions []*models.Transaction, rho, delta, eta map[int]bool, utilityArray *models.UtilityArray) {
// 	combinedSet := UnionMaps(rho, delta)
// 	combinedSet = UnionMaps(combinedSet, eta)

// 	for item := range combinedSet {
// 		totalRTWU := 0.0
// 		for _, transaction := range transactions {
// 			if ContainsItem(transaction, item) {
// 				rtwu := CalculateRTUForTransaction(transaction)
// 				totalRTWU += rtwu
// 			}
// 		}

//			utilityArray.SetRTWU(item, totalRTWU)
//		}
//	}
//
// Hàm cải tiến
func CalculateRTWUForAllItems(itemTransactionMap map[int][]*models.Transaction, combinedSet map[int]bool, eta map[int]bool, utilityArray *models.UtilityArray) {
	finalCombinedSet := UnionMaps(combinedSet, eta)

	for item := range finalCombinedSet {
		totalRTWU := 0.0

		// Truy xuất các giao dịch chứa item từ ItemTransactionMap
		transactionsWithItem, exists := itemTransactionMap[item]
		if exists {
			for _, transaction := range transactionsWithItem {
				rtwu := CalculateRTUForTransaction(transaction)
				totalRTWU += rtwu
			}
		}

		// Cập nhật RTWU cho item trong UtilityArray
		utilityArray.SetRTWU(item, totalRTWU)
	}
}
func CalculateRTUForTransaction(transaction *models.Transaction) float64 {
	rtwu := 0.0
	for _, utility := range transaction.Utilities {
		if utility > 0 {
			rtwu += utility
		}
	}
	return rtwu
}

//Hàm cũ
// func CalculateRSUForAllItems(transactions []*models.Transaction, secondary []int, utilityArray *models.UtilityArray) {
// 	for _, item := range secondary {
// 		totalRSU := 0.0

// 		for _, transaction := range transactions {
// 			if ContainsItem(transaction, item) {
// 				index := GetItemIndex(transaction, item)
// 				itemUtility := transaction.Utilities[index]
// 				remainingUtility := CalculateRemainingUtility(transaction, index+1)
// 				totalRSU += itemUtility + remainingUtility
// 			}
// 		}

//			utilityArray.SetRSU(item, totalRSU)
//		}
//	}
//
// Hàm mới
func CalculateRSUForAllItems(itemTransactionMap map[int][]*models.Transaction, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRSU := 0.0

		transactionsWithItem, exists := itemTransactionMap[item]
		if !exists {
			continue
		}

		for _, transaction := range transactionsWithItem {
			if ContainsItem(transaction, item) {
				index := GetItemIndex(transaction, item)
				itemUtility := transaction.Utilities[index]
				remainingUtility := CalculateRemainingUtility(transaction, index+1)
				totalRSU += itemUtility + remainingUtility
			}
		}

		utilityArray.SetRSU(item, totalRSU)
	}
}

func CalculateRemainingUtility(transaction *models.Transaction, startIndex int) float64 {
	remainingUtility := 0.0
	for i := startIndex; i < len(transaction.Items); i++ {
		if transaction.Utilities[i] > 0 {
			remainingUtility += transaction.Utilities[i]
		}
	}
	return remainingUtility
}

// func CalculateRSUForAllItem(transactions []*models.Transaction, X []int, secondary []int, utilityArray *models.UtilityArray) {
// 	for _, item := range secondary {
// 		totalRSU := 0.0

// 		for _, transaction := range transactions {
// 			if ContainsAllItems(transaction, X) && ContainsItem(transaction, item) {
// 				utilityX := CalculateUtilityForSet(transaction, X)
// 				indexZ := GetItemIndex(transaction, item)
// 				utilityZ := transaction.Utilities[indexZ]
// 				rru := CalculateRemainingUtility(transaction, indexZ+1)
// 				totalRSU += utilityX + utilityZ + rru
// 			}
// 		}

//			utilityArray.SetRSU(item, totalRSU)
//		}
//	}
//
// Hàm mới
func CalculateRSUForAllItem(projectedItemTransactionMap map[int][]*models.Transaction, X []int, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRSU := 0.0
		foundInAnyTransaction := false

		// fmt.Printf("\nTính RSU cho item %d:\n", item)

		// Duyệt qua tất cả các giao dịch trong `projectedItemTransactionMap`
		for _, transactions := range projectedItemTransactionMap {
			for _, transaction := range transactions {
				// Kiểm tra nếu transaction chứa tất cả các item trong X và item hiện tại
				if ContainsAllItems(transaction, X) && ContainsItem(transaction, item) {
					foundInAnyTransaction = true // Đánh dấu item có trong ít nhất một giao dịch

					// Tính toán các giá trị cần thiết
					utilityX := CalculateUtilityForSet(transaction, X)
					indexZ := GetItemIndex(transaction, item)
					utilityZ := transaction.Utilities[indexZ]
					remainingUtility := CalculateRemainingUtility(transaction, indexZ+1)

					// Cộng tổng RSU tạm thời cho item hiện tại
					totalRSU += utilityX + utilityZ + remainingUtility
					// fmt.Printf("    Utility(X): %.2f, Utility(%d): %.2f, Remaining Utility: %.2f, RSU Tổng Tạm Thời: %.2f\n", utilityX, item, utilityZ, remainingUtility, totalRSU)
				}
			}
		}

		// Nếu item không có trong bất kỳ giao dịch nào của `projectedItemTransactionMap`, bỏ qua tính toán
		if !foundInAnyTransaction {
			// fmt.Printf("Item %d không có trong bất kỳ giao dịch nào của projectedItemTransactionMap.\n", item)
			continue
		}

		// Cập nhật RSU cho item trong utilityArray
		utilityArray.SetRSU(item, totalRSU)
		// fmt.Printf("RSU(%d) = %.2f\n", item, totalRSU)
	}
}

// func CalculateRLUForAllItem(transactions []*models.Transaction, X []int, secondary []int, utilityArray *models.UtilityArray) {
// 	for _, item := range secondary {
// 		totalRLU := 0.0

// 		for _, transaction := range transactions {
// 			if ContainsAllItems(transaction, X) && ContainsItem(transaction, item) {
// 				utilityX := CalculateUtilityForSet(transaction, X)
// 				maxIndexX := FindLocationMaxIndexForSet(transaction, X)
// 				index := GetItemIndex(transaction, maxIndexX)

// 				remainingUtility := CalculateRemainingUtility(transaction, index+1)

// 				totalRLU += utilityX + remainingUtility
// 			}
// 		}

//			utilityArray.SetRLU(item, totalRLU)
//		}
//	}
func CalculateRLUForAllItem(projectedItemTransactionMap map[int][]*models.Transaction, X []int, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRLU := 0.0
		foundInAnyTransaction := false

		// fmt.Printf("\nTính RLU cho item %d:\n", item)

		// Duyệt qua tất cả các giao dịch trong `projectedItemTransactionMap`
		for _, transactions := range projectedItemTransactionMap {
			for _, transaction := range transactions {
				// Kiểm tra nếu transaction chứa tất cả các item trong X và item hiện tại
				if ContainsAllItems(transaction, X) && ContainsItem(transaction, item) {
					foundInAnyTransaction = true

					// Tính toán các giá trị cần thiết
					utilityX := CalculateUtilityForSet(transaction, X)
					maxIndexX := FindLocationMaxIndexForSet(transaction, X)

					remainingUtility := CalculateRemainingUtility(transaction, maxIndexX+1)

					// Cộng tổng RLU tạm thời cho item hiện tại
					totalRLU += utilityX + remainingUtility
					// fmt.Printf("    Utility(X): %.2f, Remaining Utility: %.2f, RLU Tổng Tạm Thời: %.2f\n", utilityX, remainingUtility, totalRLU)
				}
			}
		}

		// Nếu item không có trong bất kỳ giao dịch nào của `projectedItemTransactionMap`, bỏ qua tính toán
		if !foundInAnyTransaction {
			// fmt.Printf("Item %d không có trong bất kỳ giao dịch nào của projectedItemTransactionMap.\n", item)
			continue
		}

		// Cập nhật RLU cho item trong utilityArray
		utilityArray.SetRLU(item, totalRLU)
		// fmt.Printf("RLU(%d) = %.2f\n", item, totalRLU) // In tổng RLU cuối cùng cho item
	}
}

func CalculateUtilityForSet(transaction *models.Transaction, X []int) float64 {
	totalUtility := 0.0
	for _, item := range X {
		if ContainsItem(transaction, item) {
			index := GetItemIndex(transaction, item)
			totalUtility += transaction.Utilities[index]
		}
	}
	return totalUtility
}

func FindLocationMaxIndexForSet(transaction *models.Transaction, X []int) int {
	maxIndex := -1
	for _, item := range X {
		index := GetItemIndex(transaction, item)
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func ContainsItem(transaction *models.Transaction, item int) bool {
	for _, tItem := range transaction.Items {
		if tItem == item {
			return true
		}
	}
	return false
}

func ContainsAllItems(transaction *models.Transaction, X []int) bool {
	for _, item := range X {
		if !ContainsItem(transaction, item) {
			return false
		}
	}
	return true
}

func GetItemIndex(transaction *models.Transaction, item int) int {
	for i, tItem := range transaction.Items {
		if tItem == item {
			return i
		}
	}
	return -1
}

func UnionMaps(a, b map[int]bool) map[int]bool {
	result := make(map[int]bool)
	for k := range a {
		result[k] = true
	}
	for k := range b {
		result[k] = true
	}
	return result
}
