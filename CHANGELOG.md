## 1.2.0

- Add resources and single/list lookups for periodic jobs, health notebooks, SRE patrols, code analysis projects, and repository and container image scanner targets.
- Add monitor folders, monitor discovery sources, log analysis sources, FinOps tenants, cost budgets, and saved cost views as resources with single/list lookups.
- Add single/list lookups for Code Analysis quality gates and code profiles.
- Extend workflows with AI input extraction, schedules with timezones and input values, and defensive event filters.
- Configure Sanity checks with default input variables, application links, run timeouts, and cron schedules using `argonix_synthetic_test`.
- Add import support for workflows, synthetic tests, and the new resources. Follow API pagination and preserve equivalent JSON when refreshing state.
- Allow environment-only API key configuration through `ARGONIX_API_KEY`.

## 1.1.6

BREAKING CHANGES:

- Removed `remediation_strategy` attribute from `argonix_monitor` resource. Auto-remediation now always executes immediately when `auto_remediate = true`.

---

## 1.1.2

BREAKING CHANGES:

- Renamed `argonix_alert_rule` → `argonix_notification_rule` (resource + data sources)
- API endpoint changed from `/alert-rules/` to `/notification-rules/`

FEATURES:

- New trigger condition `test_run_complete` — fires after each synthetic test run (useful for CI/CD webhooks)
- Notification rule creation modal now has a toggle between Monitor and Synthetic Test scopes

## 1.0.0

FEATURES:

- Monitor CRUD
- Synthetic testing CRUD
- Groups CRUD
- Notification rules CRUD
- Notifications CRUD
- Status pages CRUD
- Test Management CRUD
