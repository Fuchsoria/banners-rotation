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

CREATE TABLE "clicks" (
	"slot_id" TEXT NOT NULL,
	"banner_id" TEXT NOT NULL,
	"social_demo_id" TEXT NOT NULL,
	"date" TEXT NOT NULL
);

CREATE TABLE "views" (
	"slot_id" TEXT NOT NULL,
	"banner_id" TEXT NOT NULL,
	"social_demo_id" TEXT NOT NULL,
	"date" TEXT NOT NULL
);

INSERT INTO "banners" ("id","description") VALUES ('banner1','description');
INSERT INTO "banners" ("id","description") VALUES ('banner2','description');
INSERT INTO "banners" ("id","description") VALUES ('banner3','description');
INSERT INTO "banners" ("id","description") VALUES ('banner4','description');
INSERT INTO "banners" ("id","description") VALUES ('banner5','description');

INSERT INTO "slots" ("id","description") VALUES ('slot1','description');
INSERT INTO "slots" ("id","description") VALUES ('slot2','description');
INSERT INTO "slots" ("id","description") VALUES ('slot3','description');
INSERT INTO "slots" ("id","description") VALUES ('slot4','description');
INSERT INTO "slots" ("id","description") VALUES ('slot5','description');

INSERT INTO "social_demos" ("id","description") VALUES ('social_demo1','description');
INSERT INTO "social_demos" ("id","description") VALUES ('social_demo2','description');
INSERT INTO "social_demos" ("id","description") VALUES ('social_demo3','description');
INSERT INTO "social_demos" ("id","description") VALUES ('social_demo4','description');
INSERT INTO "social_demos" ("id","description") VALUES ('social_demo5','description');

INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot1','banner1');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot1','banner2');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot1','banner3');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot1','banner4');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot1','banner5');

INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot2','banner1');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot2','banner2');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot2','banner3');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot2','banner4');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot2','banner5');

INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot3','banner1');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot3','banner2');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot3','banner3');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot3','banner4');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot3','banner5');

INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot4','banner1');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot4','banner2');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot4','banner3');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot4','banner4');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot4','banner5');

INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot5','banner1');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot5','banner2');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot5','banner3');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot5','banner4');
INSERT INTO "banners_rotation" ("slot_id","banner_id") VALUES ('slot5','banner5');

-- INSERT INTO "views" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner1','social_demo1','TEST');
-- INSERT INTO "views" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner2','social_demo2','TEST');
-- INSERT INTO "views" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner3','social_demo3','TEST');
-- INSERT INTO "views" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner4','social_demo4','TEST');
-- INSERT INTO "views" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner5','social_demo5','TEST');

-- INSERT INTO "clicks" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner1','social_demo1','TEST');
-- INSERT INTO "clicks" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner2','social_demo2','TEST');
-- INSERT INTO "clicks" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner3','social_demo3','TEST');
-- INSERT INTO "clicks" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner4','social_demo4','TEST');
-- INSERT INTO "clicks" ("slot_id","banner_id","social_demo_id","date") VALUES ('slot3','banner5','social_demo5','TEST');
