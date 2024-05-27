--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.2 (Ubuntu 16.2-1ubuntu4)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: petra_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.petra_user (
    id bigint NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL
);


ALTER TABLE public.petra_user OWNER TO postgres;

--
-- Name: petra_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.petra_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.petra_user_id_seq OWNER TO postgres;

--
-- Name: petra_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.petra_user_id_seq OWNED BY public.petra_user.id;


--
-- Name: petra_user id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.petra_user ALTER COLUMN id SET DEFAULT nextval('public.petra_user_id_seq'::regclass);


--
-- Data for Name: petra_user; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.petra_user (id, username, password) VALUES (1, 'kyle', 'hunter2');
INSERT INTO public.petra_user (id, username, password) VALUES (2, 'jack', 'baby');


--
-- Name: petra_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.petra_user_id_seq', 2, true);


--
-- PostgreSQL database dump complete
--

PGPASSWORD="hunter2" pg_dump -h localhost -p 5432 -U postgres --column-inserts postgres
--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.2 (Ubuntu 16.2-1ubuntu4)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: petra_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.petra_user (
    id bigint NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL
);


ALTER TABLE public.petra_user OWNER TO postgres;

--
-- Name: petra_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.petra_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.petra_user_id_seq OWNER TO postgres;

--
-- Name: petra_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.petra_user_id_seq OWNED BY public.petra_user.id;


--
-- Name: petra_user id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.petra_user ALTER COLUMN id SET DEFAULT nextval('public.petra_user_id_seq'::regclass);


--
-- Data for Name: petra_user; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.petra_user (id, username, password) VALUES (1, 'kyle', 'hunter2');
INSERT INTO public.petra_user (id, username, password) VALUES (2, 'jack', 'baby');


--
-- Name: petra_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.petra_user_id_seq', 2, true);


--
-- PostgreSQL database dump complete
--