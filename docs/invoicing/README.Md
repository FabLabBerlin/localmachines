Names:

- Current User Invoice = Null if there's none in Fastbill for the current month
                         Last invoice otherwise



Methods:

List of all invoicable users  GET /api/billing/months/:yyyy/:mm
and the total amount incl VAT

List of all invoices from     GET /api/billing/months/:yyyy/:mm/invoices
this month, including
drafts and canceled ones

Getting invoice               GET /api/billing/months/:yyyy/:mm/invoices/:invId

Purchases of user (like) in   GET /api/billing/months/:yyyy/:mm/users/:uid
current user invoice
                              (returns null if there's no user invoice yet?!)

Update current user invoice   POST /api/billing/months/:yyyy/:mm/users/:uid/update
                              (maybe better: POST /api/billing/months/:yyyy/:mm/users/:uid)

Create draft of current user  POST /api/billing/months/:yyyy/:mm/users/:uid/draft
invoice                       (maybe better: /api/billing/months/:yyyy/:mm/users/:uid,
                              thus coinciding this 'update')

Send                          POST /api/billing/months/:yyyy/:mm/users/:uid/send



Need to save Stuff that is sent to Fastbill for the historical read-only access, i.e. the user invoices that are not current.

Purchases table contains 2 additional columns: invoice_id and invoice_status, where invoice_status = 'draft' / 'credit' / 'outgoing'. only invoice_id = 'draft' is used for current stuff. invoice_id corresponds to the Fastbill invoice number (which does NOT equal the Fastbill invoice id) for the moment but it'll be normalized at some point, i.e. easylab.invoice_id <> fastbill.invoice_id, thus allowing for other billing backends in the future. Items with invoice_status = 'draft' are read-only. This column will be normalized at some point as well. The fields will also be added to the user_membership table.
