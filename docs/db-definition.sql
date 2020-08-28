CREATE TABLE IF NOT EXISTS Tags
(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name Text NOT NULL UNIQUE,
    EnglishName Text
);

CREATE TABLE IF NOT EXISTS VoiceActors
(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name Text NOT NULL
);

CREATE TABLE IF NOT EXISTS Circles
(
    ID INTEGER PRIMARY KEY,
    Name Text NOT NULL
);

-- Relationship tables
CREATE TABLE IF NOT EXISTS Works
(
    ID INTEGER PRIMARY KEY,
    sfw INTEGER NOT NULL,
    Name TEXT NOT NULL,
    Filepath TEXT,
    CircleID INTEGER,

    FOREIGN KEY (CircleID) REFERENCES Circles(ID)
);

CREATE TABLE IF NOT EXISTS WorkTag
(
    WorkID INTEGER,
    TagID INTEGER,

    FOREIGN KEY (WorkID) REFERENCES Works(ID),
    FOREIGN KEY (TagID) REFERENCES Tags(ID)
);

CREATE TABLE IF NOT EXISTS WorkVoiceActor
(
    WorkID INTEGER,
    VoiceActorID INTEGER,

    FOREIGN KEY (WorkID) REFERENCES Works(ID),
    FOREIGN KEY (VoiceActorID) REFERENCES VoiceActors(ID)
);

-- Insert test data
INSERT INTO Tags
    (Name)
VALUES
    ('バイノーラル/ダミヘ'),
    ('ASMR'),
    ('ラブラブ/あまあま'),
    ('言葉責め'),
    ('耳舐め'),
    ('同級生/同僚')

INSERT INTO VoiceActors
    (Name)
VALUES
    ('みやぢ')

INSERT INTO Circles
    (ID, Name)
VALUES
    (38835, 'みやぢ屋')

INSERT INTO Works
    (ID, sfw, Name, CircleID)
VALUES
    (293003, 1, '僕は彼女の耳奴隷。~いじめっ子JKが僕に夢中になるまで。~', 38835)

-- The id depends on database
INSERT INTO WorkVoiceActor
(WorkID, VoiceActorID)
VALUES
(293003, 1)

INSERT INTO WorkTag
(WorkID, TagID)
VALUES
(293003, 1),
(293003, 2),
(293003, 3),
(293003, 4),
(293003, 5),
(293003, 6)

