CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE key_value_pair AS (
  data_key VARCHAR(25),
  data_value VARCHAR(250)
);

CREATE TABLE documents (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  origin_id VARCHAR(128) NOT NULL,
  origin_name VARCHAR(128)  NOT NULL,
  class VARCHAR(50),
  paragraphs TEXT[],
  recognized_data key_value_pair[]
);

CREATE TYPE analysis_status AS ENUM ('pending','extracting','classifying','postprocessing','done');

CREATE TABLE analysis_processes (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  document_id uuid references documents(id),
  document_size INT NOT NULL,
  analysis_status analysis_status NOT NULL,
  remote_results_uri VARCHAR(250),
  remote_status VARCHAR(12),
  remote_api_version VARCHAR(25),
  remote_model VARCHAR(50),
  remote_created_datetime TIMESTAMP,
  remote_updated_datetime TIMESTAMP
);