## 1.1.1

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
