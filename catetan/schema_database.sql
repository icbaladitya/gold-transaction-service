-- MST USER
CREATE TABLE IF NOT EXISTS public.mst_user
(
    id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    nama character varying(100) COLLATE pg_catalog."default" NOT NULL,
    email character varying(150) COLLATE pg_catalog."default" NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_email_key UNIQUE (email)
);


-- MST GOLD
CREATE TABLE IF NOT EXISTS public.mst_gold
(
    id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    gold_gram numeric(12,2) NOT NULL,
    active boolean DEFAULT true,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(50) COLLATE pg_catalog."default",
    updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_by character varying(50) COLLATE pg_catalog."default",
    stock numeric NOT NULL DEFAULT 0,
    code character varying(8) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT mst_gold_pkey PRIMARY KEY (id)
);


-- USER BALANCE
CREATE TABLE IF NOT EXISTS public.user_balance
(
    id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    user_id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    idr_balance numeric(15,2) DEFAULT 0.00,
    gold_balance numeric(12,2) DEFAULT 0.00,
    version integer DEFAULT 1,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT user_balance_pkey PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES public.mst_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
CREATE INDEX IF NOT EXISTS idx_user_balance_user_id
    ON public.user_balance USING btree
    (user_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;


-- GOLD PRICES
CREATE TABLE IF NOT EXISTS public.gold_prices
(
    id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    mst_gold_id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    buy_price numeric(15,2) NOT NULL,
    sell_price numeric(15,2) NOT NULL,
    buy_price_per_gram numeric(15,2) NOT NULL,
    sell_price_per_gram numeric(15,2) NOT NULL,
    version integer DEFAULT 1,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(32) COLLATE pg_catalog."default",
    CONSTRAINT gold_prices_pkey PRIMARY KEY (id),
    CONSTRAINT fk_gold_prices_mst FOREIGN KEY (mst_gold_id)
        REFERENCES public.mst_gold (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
CREATE INDEX IF NOT EXISTS gold_prices_idx1
    ON public.gold_prices USING btree
    (mst_gold_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;


-- GOLD TRANSACTION HEADER
CREATE TABLE IF NOT EXISTS public.gold_trx_hdr
(
    id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    user_id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    type character varying(20) COLLATE pg_catalog."default" NOT NULL,
    gold_gram numeric(12,2) DEFAULT 0.00,
    gold_idr numeric(15,2) DEFAULT 0.00,
    status character varying(30) COLLATE pg_catalog."default" DEFAULT 'PENDING'::character varying,
    description text COLLATE pg_catalog."default",
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT gold_trx_hdr_pkey PRIMARY KEY (id),
    CONSTRAINT fk_gold_trx_hdr_user FOREIGN KEY (user_id)
        REFERENCES public.mst_user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
)
CREATE INDEX IF NOT EXISTS idx_gold_trx_hdr_user_id
    ON public.gold_trx_hdr USING btree
    (user_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;


-- GOLD TRANSACTION DETAIL
CREATE TABLE IF NOT EXISTS public.gold_trx_dtl
(
    id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    gold_trx_hdr_id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    gold_prices_id character varying(32) COLLATE pg_catalog."default" NOT NULL,
    gold_gram numeric(12,4) NOT NULL,
    buy_price numeric(15,2) NOT NULL,
    sell_price numeric(15,2) NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT gold_trx_dtl_pkey PRIMARY KEY (id),
    CONSTRAINT fk_gold_trx_dtl_hdr FOREIGN KEY (gold_trx_hdr_id)
        REFERENCES public.gold_trx_hdr (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_gold_trx_dtl_prices FOREIGN KEY (gold_prices_id)
        REFERENCES public.gold_prices (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
)
CREATE INDEX IF NOT EXISTS idx_gold_trx_dtl_hdr_id
    ON public.gold_trx_dtl USING btree
    (gold_trx_hdr_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_gold_trx_dtl_prices_id
    ON public.gold_trx_dtl USING btree
    (gold_prices_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;