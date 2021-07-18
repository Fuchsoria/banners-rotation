CREATE TABLE "slots" (
	"id" TEXT NOT NULL,
	"description" TEXT DEFAULT '',
	PRIMARY KEY ("id")
);

CREATE TABLE "banners" (
	"id" TEXT NOT NULL,
	"description" TEXT DEFAULT '',
	PRIMARY KEY ("id")
);

CREATE TABLE "social-demos" (
	"id" TEXT NOT NULL,
	"description" TEXT DEFAULT '',
	PRIMARY KEY ("id")
);

CREATE TABLE "banners-rotation" (
	"id" TEXT NOT NULL,
	"slot_id" TEXT NOT NULL,
	"banner_id" TEXT NOT NULL,
	PRIMARY KEY ("id")
);