-- Testing with Production data is always better but risky
-- This script tries to reduce risk as much as possible:
--
-- - Jabber IDs are nulled
-- - User Names where user_id > 20 are anonymized
-- - netswitch hosts, url ons and offs are nulled
-- - all passwords are set to 123456

UPDATE locations SET local_ip = '',
                     xmpp_id = '';

UPDATE machines SET netswitch_host = '',
                    netswitch_url_on = '',
                    netswitch_url_off = '';

UPDATE auth SET hash = 'f7c19341b9c14c27136b4653514f1b7d7ad16b1c2306181481956fb93b749c74c0337dcb2622d86644d83406e98d45b782c4588a3f94d25ce79547d26f7a11ae',
                salt = '53d2ab2f6759bf41bff8a4bbb93975fb31cd4a914a6750f3d021ba3be5ea8fd4';

UPDATE user SET client_id = 0 WHERE id <> 19;

UPDATE user SET first_name = concat('f', '-', id),
                last_name = concat('l', '-', id),
                username = concat('f', '-', id, '.', 'l', '-', id),
                email = concat('f', '-', id, '.', 'l', '-', id, '@example.com')
            WHERE (id > 20 OR id = 11 OR id = 14) AND username <> 'testuser';
