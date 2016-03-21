START TRANSACTION;


DELETE purchases
FROM purchases
JOIN user ON user.id = purchases.user_id
WHERE test_user = 1;


DELETE FROM membership
WHERE id >= 10000 AND location_id = 2;


DELETE user_membership
FROM user_membership
JOIN user ON user.id = user_membership.user_id
WHERE test_user = 1;


DELETE FROM user WHERE test_user = 1;


COMMIT;
