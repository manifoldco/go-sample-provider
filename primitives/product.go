package primitives

// Product is the only product supported by this provider.
const Product = "bear"

// Region is the only region supported by this provider.
const Region = "all::global"

// Plans contains a list of plans supported by this provider.
var Plans = []string{"ursa-minor", "ursa-major"}

// ValidPlan return whether a plan is contained in the list of available plans.
func ValidPlan(plan string) bool {
	for _, p := range Plans {
		if p == plan {
			return true
		}
	}

	return false
}
