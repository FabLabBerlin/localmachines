Names:

Purchases table contains 2 additional columns: `invoice_id` and `invoice_status`, where `invoice_status = {"draft" | "credit" | "outgoing"}`. only `invoice_id = "draft"` is used for current stuff. `invoice_id` corresponds to the `invoice` table's column `id`. Items with `invoice_status <> 'draft'` are read-only. This column will be normalized at some point as well. The fields are also in the `user_membership` table.

In the table `invoices` only the combination `(location_id, user_id, month, year)` is guaranteed to be unique.
