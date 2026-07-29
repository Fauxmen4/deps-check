package analyzer

func FilterUpdatable(report []DependencyReport) []DependencyReport {
	filtered := make([]DependencyReport, 0, len(report))
	for _, r := range report {
		if r.Update != UpdateNone {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
