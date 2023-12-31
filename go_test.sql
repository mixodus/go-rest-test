PGDMP                  	    {            go_test    16.0    16.0     4           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            5           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            6           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            7           1262    16398    go_test    DATABASE     i   CREATE DATABASE go_test WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'C';
    DROP DATABASE go_test;
                postgres    false            �            1259    16432    banks    TABLE       CREATE TABLE public.banks (
    id text DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    bank_code text,
    bank_initial text
);
    DROP TABLE public.banks;
       public         heap    postgres    false            �            1259    16410    players    TABLE     l  CREATE TABLE public.players (
    id text DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_name text,
    last_name text,
    password text,
    email text NOT NULL,
    phone text,
    balance bigint,
    players_bank_id text
);
    DROP TABLE public.players;
       public         heap    postgres    false            �            1259    16441    players_banks    TABLE     Q  CREATE TABLE public.players_banks (
    id text DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    player_id text NOT NULL,
    bank_id text NOT NULL,
    bank_account_number text,
    bank_account_name text
);
 !   DROP TABLE public.players_banks;
       public         heap    postgres    false            �            1259    16421    token_sessions    TABLE       CREATE TABLE public.token_sessions (
    id text DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    player_id text NOT NULL,
    token text NOT NULL
);
 "   DROP TABLE public.token_sessions;
       public         heap    postgres    false            �            1259    16456    transactions    TABLE     �  CREATE TABLE public.transactions (
    id text DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    player_id uuid NOT NULL,
    players_bank_id uuid NOT NULL,
    amount bigint,
    status text DEFAULT 'pending'::text,
    transaction_type text NOT NULL,
    file_name character varying(255),
    notes character varying(255)
);
     DROP TABLE public.transactions;
       public         heap    postgres    false            /          0    16432    banks 
   TABLE DATA           f   COPY public.banks (id, created_at, updated_at, deleted_at, name, bank_code, bank_initial) FROM stdin;
    public          postgres    false    217   y       -          0    16410    players 
   TABLE DATA           �   COPY public.players (id, created_at, updated_at, deleted_at, first_name, last_name, password, email, phone, balance, players_bank_id) FROM stdin;
    public          postgres    false    215           0          0    16441    players_banks 
   TABLE DATA           �   COPY public.players_banks (id, created_at, updated_at, deleted_at, player_id, bank_id, bank_account_number, bank_account_name) FROM stdin;
    public          postgres    false    218   �!       .          0    16421    token_sessions 
   TABLE DATA           b   COPY public.token_sessions (id, created_at, updated_at, deleted_at, player_id, token) FROM stdin;
    public          postgres    false    216   g"       1          0    16456    transactions 
   TABLE DATA           �   COPY public.transactions (id, created_at, updated_at, deleted_at, player_id, players_bank_id, amount, status, transaction_type, file_name, notes) FROM stdin;
    public          postgres    false    219   �"       �           2606    16440    banks banks_pkey 
   CONSTRAINT     N   ALTER TABLE ONLY public.banks
    ADD CONSTRAINT banks_pkey PRIMARY KEY (id);
 :   ALTER TABLE ONLY public.banks DROP CONSTRAINT banks_pkey;
       public            postgres    false    217            �           2606    16420    players idx_players_email 
   CONSTRAINT     U   ALTER TABLE ONLY public.players
    ADD CONSTRAINT idx_players_email UNIQUE (email);
 C   ALTER TABLE ONLY public.players DROP CONSTRAINT idx_players_email;
       public            postgres    false    215            �           2606    16449     players_banks players_banks_pkey 
   CONSTRAINT     ^   ALTER TABLE ONLY public.players_banks
    ADD CONSTRAINT players_banks_pkey PRIMARY KEY (id);
 J   ALTER TABLE ONLY public.players_banks DROP CONSTRAINT players_banks_pkey;
       public            postgres    false    218            �           2606    16417    players players_pkey 
   CONSTRAINT     R   ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.players DROP CONSTRAINT players_pkey;
       public            postgres    false    215            �           2606    16429 "   token_sessions token_sessions_pkey 
   CONSTRAINT     `   ALTER TABLE ONLY public.token_sessions
    ADD CONSTRAINT token_sessions_pkey PRIMARY KEY (id);
 L   ALTER TABLE ONLY public.token_sessions DROP CONSTRAINT token_sessions_pkey;
       public            postgres    false    216            �           2606    16431 '   token_sessions token_sessions_token_key 
   CONSTRAINT     c   ALTER TABLE ONLY public.token_sessions
    ADD CONSTRAINT token_sessions_token_key UNIQUE (token);
 Q   ALTER TABLE ONLY public.token_sessions DROP CONSTRAINT token_sessions_token_key;
       public            postgres    false    216            �           2606    16465    transactions transactions_pkey 
   CONSTRAINT     \   ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
 H   ALTER TABLE ONLY public.transactions DROP CONSTRAINT transactions_pkey;
       public            postgres    false    219            �           1259    16455    idx_players_banks_deleted_at    INDEX     \   CREATE INDEX idx_players_banks_deleted_at ON public.players_banks USING btree (deleted_at);
 0   DROP INDEX public.idx_players_banks_deleted_at;
       public            postgres    false    218            �           2606    16450 #   players_banks fk_players_banks_bank    FK CONSTRAINT     �   ALTER TABLE ONLY public.players_banks
    ADD CONSTRAINT fk_players_banks_bank FOREIGN KEY (bank_id) REFERENCES public.banks(id);
 M   ALTER TABLE ONLY public.players_banks DROP CONSTRAINT fk_players_banks_bank;
       public          postgres    false    217    218    3478            �           2606    16521 %   players_banks fk_players_players_bank    FK CONSTRAINT     �   ALTER TABLE ONLY public.players_banks
    ADD CONSTRAINT fk_players_players_bank FOREIGN KEY (player_id) REFERENCES public.players(id);
 O   ALTER TABLE ONLY public.players_banks DROP CONSTRAINT fk_players_players_bank;
       public          postgres    false    215    218    3472            /   ~   x��ͱ
�  �Y�½\�;5���\2g�#ZJ�Bڥ_:vj��1Ik�B$T-����7��!r���=G�������'���E��J��).�����&1��d����Y�x��a���n����h4�      -   �  x���Mo�@��>p�֙����rjD�G Pb�ˮ�`Z�K	_��5���&Q��0�;�y�IAY��`(bd"�	f'���N�(  ������!(2q,.@}�*iN-  {M�Q���$�U�-���ؠN��P��w[LF/��M�x��2+~v����	�s;�'��6�a���p,vvuܮ���̗vQDi����6J ��_����1�Un`B��Y���[�i��3�<�#���?��O��l�lߒwF�d��ݽV��c�&�<�;���M��S��wë8�׃�Hz�0��? �k����&e�d"�)��)��p)�أ��,
	\6�G���-�J\-HIk�)�]�&�C���-�KPŝ�#�L��vG�L��-�N����d7��q\��c��ۯ�}�9z�$&�K����i@�8�%��䵬=E�Z�7�.�C      0   �   x�uɻ�0 �ڞ��SN��J�L�Ʋ�;
�/8�/��4��Ԓ�!%a�y��
#ٍ�`<�7��w����|�ɂs���@/�<�5�>)� F_k&a�� �L�]Vr�_��Zy������Z���,�      .      x������ � �      1   u  x�ő�nAF�ݧ�yc�x��< D�4�ɕr�Fw7����E��F�d�|k���8�
&1��y	\�28t��D�PXPf�*ߡ�i�&?���<V��fn�@Ԭ���T*eVN�0x��2Bf���7l-��@�߰����}�-��������~w\Ӷ�r�.�����(�\[�����3�yC��|���q���~�r?�\�$M@%��b	bZ\�(��PZH��D�[�x]����w�v��ׯ��q�/�v<�߷�5�-��{���Ɣ�aM�bdъ�67/�ׇ�G��T)���_�1��FP)����ZW�Ps�z�����+��XƙD���*w�^�Si�#���'=��<��w���     