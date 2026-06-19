-- MST USER
INSERT INTO "user" (id, nama, email, created_by) VALUES
(replace(gen_random_uuid()::text, '-', ''), 'Bambang Pamungkas', 'bambang@email.com', 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'Siti Badriah', 'siti@email.com', 'system');

-- USER BALANCE
INSERT INTO public.user_balance(
	id, user_id, idr_balance, gold_balance, version, created_by)
	VALUES 
	(replace(gen_random_uuid()::text, '-', ''), 'a1b2c3d4e5f67a8b9c0d1e2f3a4b5c6d', 30000000, 1, 1, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'e02b2c3d479f47ac10b58cc4372a5670', 60000000, 2, 1, 'system');

-- MST GOLD
INSERT INTO mst_gold (id, code, gold_gram, active, created_by) VALUES
(replace(gen_random_uuid()::text, '-', ''), 'G0001', 0.0100, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G0002', 0.0200, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G0005', 0.0500, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G0010', 0.1000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G0025', 0.2500, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G0050', 0.5000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G1000', 1.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G2000', 2.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G3000', 3.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G5000', 5.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G10K0', 10.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G25K0', 25.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G50K0', 50.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G100K', 100.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G250K', 250.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G500K', 500.0000, true, 'system'),
(replace(gen_random_uuid()::text, '-', ''), 'G1000K', 1000.0000, true, 'system');

-- GOLD PRICES
INSERT INTO public.gold_prices(
	id, mst_gold_id, buy_price, sell_price, buy_price_per_gram, sell_price_per_gram, version, created_by)
	VALUES 
	(replace(gen_random_uuid()::text, '-', ''), '62bc723e66bc4d309e5a5ca42f7738af', 30000, 29000, 3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), '958509dd4002418089c295af134bf09b', 60000, 58000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'bc473713cc704ceeb76b6fe380999046', 150000, 149000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), '38ea35922d4a4566b7bbfd6255a5df4f', 300000, 290000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), '2e4e4620227348ffb817f18b66a4a525', 500000, 490000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'ea2b30be700e495ca451a7732e8f1faf', 15000000, 14500000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), '296423ede7684a8fa77392c03708b87f', 3000000, 2900000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'f78115647ddd458eb12570be3a8196ce', 6000000, 5900000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'deca0a227a26412e8d11ddafc2207818', 9000000, 8900000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'e922082173044ec39b22a5fc9cc2af08', 15000000, 14900000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'f3d80796732a4539936188661d73a67a', 30000000, 29000000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'a937b1d4f6964ecb93c2cd25bbb0200a', 61500000, 61400000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), '91e2f2a57c204ac1a0a1c0d557a205c0', 150000000, 149000000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'b17d9cbd7ff44598b3f83830d2943ac9', 300000000, 299000000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), '6847834ea43d4b509da8f6ed9dda3164', 615000000, 614000000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'b1b3bfe74dc14b51a10b99ba67783a2e', 1500000000, 1499000000,  3000000, 2900000, 2, 'system'),
	(replace(gen_random_uuid()::text, '-', ''), 'aaab9d5a8c9e467289385678d57bd7dc', 3000000000, 2999000000,  3000000, 2900000, 2, 'system');
