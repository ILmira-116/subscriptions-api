-- 002_add_unique_index_subscriptions.up.sql
CREATE UNIQUE INDEX uniq_user_service_start
ON subscriptions(user_id, service, start_date); 