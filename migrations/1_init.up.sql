DROP SEQUENCE IF EXISTS "public"."links_id_seq";
CREATE SEQUENCE "public"."links_id_seq"
    INCREMENT 1
    MINVALUE  1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1
;
ALTER SEQUENCE "public"."links_id_seq" OWNER TO "postgres";

DROP TABLE IF EXISTS "public"."links";
CREATE TABLE "public"."links" (
    "id" int8 NOT NULL DEFAULT nextval('links_id_seq'::regclass),
    "full_link" text COLLATE "pg_catalog"."default",
    "short_link" varchar(10) COLLATE "pg_catalog"."default"
);

ALTER TABLE "public"."links" ADD CONSTRAINT "full_link" UNIQUE ("full_link");
ALTER TABLE "public"."links" ADD CONSTRAINT "links_pkey" PRIMARY KEY ("id");
