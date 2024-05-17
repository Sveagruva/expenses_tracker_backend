CREATE TABLE "Transactions" (
    "Id" INTEGER PRIMARY KEY,
    "Price" INTEGER,
    "CategoryId" INTEGER,
    "CreatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "UserId" INTEGER,
    FOREIGN KEY ("CategoryId") REFERENCES "TransactionCategories"("Id"),
    FOREIGN KEY ("UserId") REFERENCES "Users"("Id")
);

CREATE TABLE "TransactionCategories" (
    "Id" INTEGER PRIMARY KEY,
    "UserId" INTEGER,
    "Name" TEXT,
    "Color" TEXT, -- Assuming hexadecimal color codes are stored as text
    FOREIGN KEY ("UserId") REFERENCES "Users"("Id")
);
