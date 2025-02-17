CREATE TABLE applicants (
    id UUID PRIMARY KEY ,
    name VARCHAR(255),
    employment_status VARCHAR(50),
    marital_status VARCHAR(50),
    sex VARCHAR(10),
    date_of_birth DATE
);

CREATE TABLE household_members (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    employment_status VARCHAR(50),
    sex VARCHAR(10),
    date_of_birth DATE,
    relation VARCHAR(50),
    school_level VARCHAR(50),
    applicant_id UUID 
    FOREIGN KEY (applicant_id) REFERENCES applicants(id)
);

CREATE TABLE schemes (
    id UUID PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE criteria (
    scheme_id UUID,
    employment_status VARCHAR(50),
    marital_status VARCHAR(50),
    has_children BOOLEAN,
    FOREIGN KEY (scheme_id) REFERENCES schemes(id)
);

CREATE TABLE benefits (
    id UUID PRIMARY KEY,
    scheme_id UUID,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(10,2),
    FOREIGN KEY (scheme_id) REFERENCES schemes(id)
);

CREATE TABLE applications (
    application_id UUID PRIMARY KEY,
    applicant_id UUID,
    scheme_id UUID,
    status VARCHAR(50),
    FOREIGN KEY (applicant_id) REFERENCES applicants(id),
    FOREIGN KEY (scheme_id) REFERENCES schemes(id)
);
