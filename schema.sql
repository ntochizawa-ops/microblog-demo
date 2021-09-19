CREATE TABLE Messages (
  MessageId     STRING(36)  NOT NULL,
  CreatedAt     TIMESTAMP   NOT NULL,
  Name          STRING(MAX) NOT NULL,
  Body          STRING(MAX) NOT NULL,
  WrittenAt     STRING(MAX) NOT NULL,
) PRIMARY KEY (MessageId, CreatedAt DESC);
