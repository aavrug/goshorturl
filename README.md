# goshorturl

## Record Table

Create this table in postgres database.

```
    DROP TABLE IF EXISTS "records";

    CREATE TABLE "public"."records" (
        "id" integer DEFAULT nextval('records_id_seq') NOT NULL,
        "code" character(100) NOT NULL,
        "url" character(250) NOT NULL,
        "created_time" integer NOT NULL,
        CONSTRAINT "records_code" UNIQUE ("code"),
        CONSTRAINT "records_id" PRIMARY KEY ("id"),
        CONSTRAINT "records_url" UNIQUE ("url")
    ) WITH (oids = false);
```


