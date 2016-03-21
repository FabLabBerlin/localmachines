
-- No need to populate user_locations because they can't login anyways...

START TRANSACTION;

INSERT INTO user (id, 
                  first_name, 
                  last_name, 
                  username, 
                  email, 
                  user_role,
                  client_id,
                  b2b,
                  vat_rate,
                  test_user)
SELECT user.id + 10000, 
            CONCAT('TEST ', REVERSE(first_name)), 
            CONCAT('TEST ', REVERSE(last_name)), 
            CONCAT('test_demo_', REVERSE(username)), 
            CONCAT('test_demo_', user.id, '@example.com'), 
            'member',
            -1,
            0,
            19,
            1
FROM user
JOIN user_locations ON user.id = user_locations.user_id
AND user_locations.location_id = 1
AND user_locations.user_role = "member";


INSERT INTO purchases (location_id, TYPE, user_id, time_start, time_end, quantity, price_per_unit, price_unit, machine_id)
SELECT 2,
       purchases.TYPE,
                 purchases.user_id + 10000,
                 purchases.time_start,
                 purchases.time_end,
                 10 * rand() * purchases.quantity,
                 10 * rand() * purchases.price_per_unit,
                 purchases.price_unit,
                 46
FROM purchases
JOIN user_locations ON purchases.user_id = user_locations.user_id
AND user_locations.location_id = 1
AND user_locations.user_role = "member";


INSERT INTO membership (
id,
location_id,
title,
short_name,
duration_months,
monthly_price,
machine_price_deduction,
affected_machines
)
SELECT
membership.id + 10000,
2,
CONCAT('TEST ', REVERSE(title)),
CONCAT('TEST_', REVERSE(short_name)),
duration_months,
50 * rand() * monthly_price,
machine_price_deduction,
'[46]'
FROM membership;


INSERT INTO user_membership
(
user_id,
membership_id,
start_date,
end_date,
auto_extend,
is_terminated
)
SELECT
user_membership.user_id + 10000,
membership_id + 10000,
start_date,
end_date,
auto_extend,
is_terminated
FROM user_membership
JOIN user_locations ON user_membership.user_id = user_locations.user_id
AND user_locations.location_id = 1
AND user_locations.user_role = "member";


COMMIT;
