CREATE TABLE IF NOT EXISTS users(
   uid                      UUID PRIMARY KEY NOT NULL,
   auth_provider            varchar(255) NOT NULL,
   provider_id              varchar(255) NOT NULL,
   user_name                varchar(255) NOT NULL,   
   details                  JSONB NOT NULL,
   created_at               TIMESTAMP DEFAULT now()
   updated_at               TIMESTAMP DEFAULT now()
   CONSTRAINT provider_unique UNIQUE (auth_provider, provider_id)
);
