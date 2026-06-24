package prefix

// Resolve picks the issue-ID prefix from the first non-empty source in
// precedence order: the --prefix flag, then the BEES_PREFIX env var, then the
// caller-supplied fallback (today the per-repo config prefix). When nothing is
// set it defaults to "bees". The env and flag values are passed in so the
// choice stays pure and unit-testable without touching the environment.
func Resolve(flag, env, fallback string) string {
	if flag != "" {
		return flag
	}

	if env != "" {
		return env
	}

	if fallback != "" {
		return fallback
	}

	return "bees"
}
