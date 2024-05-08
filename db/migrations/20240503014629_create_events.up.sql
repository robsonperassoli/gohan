CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(40) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    verb VARCHAR(100) NOT NULL,
    subject_id VARCHAR(40) NOT NULL,
    direct_object_id VARCHAR(40) NOT NULL,
    indirect_object_id VARCHAR(40),
    prepositional_object_id VARCHAR(40),
    context TEXT,
    CONSTRAINT fk_subject_id
        FOREIGN KEY(subject_id)
            REFERENCES objects(id),
    CONSTRAINT fk_direct_object_id
        FOREIGN KEY(direct_object_id)
            REFERENCES objects(id),
    CONSTRAINT fk_indirect_object_id
        FOREIGN KEY(indirect_object_id)
            REFERENCES objects(id),
    CONSTRAINT fk_prepositional_object_id
        FOREIGN KEY(prepositional_object_id)
            REFERENCES objects(id)
);
