--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1 (Ubuntu 14.1-1.pgdg21.10+1)
-- Dumped by pg_dump version 14.1 (Ubuntu 14.1-1.pgdg21.10+1)

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
-- Name: prices; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.prices (
    id integer NOT NULL,
    winter_price integer,
    summer_price integer,
    room_id integer NOT NULL
);


ALTER TABLE public.prices OWNER TO chrismo;

--
-- Name: prices_id_seq; Type: SEQUENCE; Schema: public; Owner: chrismo
--

CREATE SEQUENCE public.prices_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.prices_id_seq OWNER TO chrismo;

--
-- Name: prices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrismo
--

ALTER SEQUENCE public.prices_id_seq OWNED BY public.prices.id;


--
-- Name: reservations; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.reservations (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    start_date date NOT NULL,
    end_date date NOT NULL,
    user_id integer NOT NULL,
    rooms_id integer NOT NULL
);


ALTER TABLE public.reservations OWNER TO chrismo;

--
-- Name: reservations_id_seq; Type: SEQUENCE; Schema: public; Owner: chrismo
--

CREATE SEQUENCE public.reservations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reservations_id_seq OWNER TO chrismo;

--
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrismo
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- Name: restrictions; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.restrictions (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    restriction_name character varying(50)
);


ALTER TABLE public.restrictions OWNER TO chrismo;

--
-- Name: restrictions_id_seq; Type: SEQUENCE; Schema: public; Owner: chrismo
--

CREATE SEQUENCE public.restrictions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.restrictions_id_seq OWNER TO chrismo;

--
-- Name: restrictions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrismo
--

ALTER SEQUENCE public.restrictions_id_seq OWNED BY public.restrictions.id;


--
-- Name: rooms; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.rooms (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    room_name character varying(30)
);


ALTER TABLE public.rooms OWNER TO chrismo;

--
-- Name: rooms_id_seq; Type: SEQUENCE; Schema: public; Owner: chrismo
--

CREATE SEQUENCE public.rooms_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rooms_id_seq OWNER TO chrismo;

--
-- Name: rooms_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrismo
--

ALTER SEQUENCE public.rooms_id_seq OWNED BY public.rooms.id;


--
-- Name: rooms_restrictions; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.rooms_restrictions (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    start_date date NOT NULL,
    end_date date NOT NULL,
    room_id integer NOT NULL,
    reservations_id integer NOT NULL,
    restrictions_id integer NOT NULL
);


ALTER TABLE public.rooms_restrictions OWNER TO chrismo;

--
-- Name: rooms_restrictions_id_seq; Type: SEQUENCE; Schema: public; Owner: chrismo
--

CREATE SEQUENCE public.rooms_restrictions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rooms_restrictions_id_seq OWNER TO chrismo;

--
-- Name: rooms_restrictions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrismo
--

ALTER SEQUENCE public.rooms_restrictions_id_seq OWNED BY public.rooms_restrictions.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO chrismo;

--
-- Name: users; Type: TABLE; Schema: public; Owner: chrismo
--

CREATE TABLE public.users (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    phone character varying(25),
    email character varying(40) NOT NULL,
    password character varying(60),
    access_level integer
);


ALTER TABLE public.users OWNER TO chrismo;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: chrismo
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO chrismo;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrismo
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: prices id; Type: DEFAULT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.prices ALTER COLUMN id SET DEFAULT nextval('public.prices_id_seq'::regclass);


--
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- Name: restrictions id; Type: DEFAULT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.restrictions ALTER COLUMN id SET DEFAULT nextval('public.restrictions_id_seq'::regclass);


--
-- Name: rooms id; Type: DEFAULT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms ALTER COLUMN id SET DEFAULT nextval('public.rooms_id_seq'::regclass);


--
-- Name: rooms_restrictions id; Type: DEFAULT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms_restrictions ALTER COLUMN id SET DEFAULT nextval('public.rooms_restrictions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: prices prices_pkey; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.prices
    ADD CONSTRAINT prices_pkey PRIMARY KEY (id);


--
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- Name: restrictions restrictions_pkey; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.restrictions
    ADD CONSTRAINT restrictions_pkey PRIMARY KEY (id);


--
-- Name: rooms rooms_pkey; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_pkey PRIMARY KEY (id);


--
-- Name: rooms_restrictions rooms_restrictions_pkey; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms_restrictions
    ADD CONSTRAINT rooms_restrictions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: rooms_restrictions_reservations_id_idx; Type: INDEX; Schema: public; Owner: chrismo
--

CREATE INDEX rooms_restrictions_reservations_id_idx ON public.rooms_restrictions USING btree (reservations_id);


--
-- Name: rooms_restrictions_room_id_idx; Type: INDEX; Schema: public; Owner: chrismo
--

CREATE INDEX rooms_restrictions_room_id_idx ON public.rooms_restrictions USING btree (room_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: chrismo
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: prices prices_room_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.prices
    ADD CONSTRAINT prices_room_id_fkey FOREIGN KEY (room_id) REFERENCES public.rooms(id) ON DELETE CASCADE;


--
-- Name: reservations reservations_rooms_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_rooms_id_fkey FOREIGN KEY (rooms_id) REFERENCES public.restrictions(id) ON DELETE CASCADE;


--
-- Name: reservations reservations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: rooms_restrictions rooms_restrictions_reservations_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms_restrictions
    ADD CONSTRAINT rooms_restrictions_reservations_id_fkey FOREIGN KEY (reservations_id) REFERENCES public.reservations(id) ON DELETE CASCADE;


--
-- Name: rooms_restrictions rooms_restrictions_restrictions_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms_restrictions
    ADD CONSTRAINT rooms_restrictions_restrictions_id_fkey FOREIGN KEY (restrictions_id) REFERENCES public.restrictions(id) ON DELETE CASCADE;


--
-- Name: rooms_restrictions rooms_restrictions_room_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrismo
--

ALTER TABLE ONLY public.rooms_restrictions
    ADD CONSTRAINT rooms_restrictions_room_id_fkey FOREIGN KEY (room_id) REFERENCES public.rooms(id);


--
-- PostgreSQL database dump complete
--

