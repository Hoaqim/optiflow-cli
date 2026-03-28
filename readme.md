# OptiFlow-CLI 🚀

OptiFlow-CLI is a blazing-fast, local CI/CD pipeline analyzer built in Go. It allows developers to validate their GitHub Actions logic, enforce supply chain security, and generate "shift-left" cost projections before committing code.

## Features

- 🔍 **Syntax Validation:** Parses and validates your YAML workflow schemas instantly.
- 🛡️ **Security Enforcement:** Scans for anti-patterns and supply chain vulnerabilities:
  - Detects non-deterministic installations (e.g., forcing `npm ci` over `npm install`).
  - Blocks mutable action references (prevents using `@main` or `@master` in favor of commit SHAs).
- 💰 **Shift-Left FinOps:** Calculates estimated and worst-case execution costs based on runner types and timeout bounds. 

## Installation

Ensure you have Go 1.26+ installed, then run:

```bash
go install [github.com/Hoaqim/optiflow-cli@latest](https://github.com/Hoaqim/optiflow-cli@latest)
(Alternatively, you can clone the repository and run go build -o optiflow main.go)
```

## Usage
Analyze a Workflow
Analyze a specific workflow file to get security violations and cost projections:

```Bash
optiflow analyze .github/workflows/main.yml
```
Example Output:

Analyzing .github/workflows/main.yml...
Syntax Validated (Found 2 jobs in 'CI Pipeline')

Security Policies Enforced: Found 1 Violations
  [HIGH] Mutable Action Reference (Job: build, Step: actions/checkout@main)
      -> Action uses a mutable branch (actions/checkout@main). Pin to a specific commit SHA to prevent supply chain attacks.

Shift-Left Cost Projections:
  - Job [build] on 'ubuntu-latest' ($0.008/min)
  - Job [deploy] on 'macos-latest' ($0.080/min)

  Est. Cost Per Run:   $0.440
  Worst-Case Timeout:  $31.680
  Est. Monthly Impact: $44.00 (Assuming 100 runs/mo)
Global Flags
-v, --verbose: Enable verbose output for debugging.

## Check Version
```Bash
optiflow version
```
## Contributing
Contributions are welcome! If you'd like to add new security rules or improve the FinOps cost matrix, please fork the repository and submit a Pull Request.