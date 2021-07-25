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

CREATE TABLE "social_demos" (
	"id" TEXT NOT NULL,
	"description" TEXT DEFAULT '',
	PRIMARY KEY ("id")
);

CREATE TABLE "banners_rotation" (
	"slot_id" TEXT NOT NULL,
	"banner_id" TEXT NOT NULL
);

CREATE TABLE "session_clicks" (
	"slot_id" TEXT NOT NULL,
	"banner_id" TEXT NOT NULL,
	"social_demo_id" TEXT NOT NULL,
	"date" TEXT NOT NULL
);

CREATE TABLE "session_views" (
	"slot_id" TEXT NOT NULL,
	"banner_id" TEXT NOT NULL,
	"social_demo_id" TEXT NOT NULL,
	"date" TEXT NOT NULL
);
