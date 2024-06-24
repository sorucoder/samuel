-- +migrate Up
INSERT INTO `roles` (`id`, `name`, `priority`)
VALUES
    ('administrator', 'Administrator', 1),
    ('instructor',    'Instructor',    2),
    ('supervisor',    'Supervisor',    3),
    ('student',       'Student',       4);

INSERT INTO `users` (`identity`, `password_hash`, `role_id`)
VALUES
    ('paulmazza',                  NULL,                                                           'administrator'),
    ('gsantella',                  NULL,                                                           'instructor'),
    ('bselfridge',                 NULL,                                                           'instructor'),
    ('npage',                      NULL,                                                           'instructor'),
    ('rgority',                    NULL,                                                           'instructor'),
    ('rgallagher@pssolutions.net', '$2y$10$A207LVcfHDrmWyfTzV1TQugv2Yo0FPKYTwECbk3QImQfV0TiLdAqK', 'supervisor'),
    ('jack@jackchristensen.com',   '$2y$05$0Lv4E0TaKKfPSi0GVCgoaOFZoe5kUSa2DwA09L3XByi04MPAiQPNy', 'supervisor'),
    ('mgermano79',                 NULL,                                                           'student'),
    ('bkletzing32',                NULL,                                                           'student');

INSERT INTO `administrators` (`user_uuid`, `first_name`, `last_name`, `email`, `phone`)
SELECT
    `users`.`uuid`             AS `user_uuid`,
    'Paul'                     AS `first_name`,
    'Mazza'                    AS `last_name`,
    'paulmazza@southhills.edu' AS `email`,
    '8142347755'               AS `phone`
FROM `users`
WHERE `identity` = 'paulmazza';

INSERT INTO `campuses` (`id`, `name`, `address`, `city`, `state`, `zip`, `phone`)
VALUES
    ('sce', 'State College Campus', '480 Waupelani Drive', 'State College', 'PA', '16801', '8142377755'),
    ('alt', 'Altoona Campus',       '508 58th Street',     'Altoona',       'PA', '16602', '8149446134');

INSERT INTO `instructors` (`user_uuid`, `first_name`, `last_name`, `email`, `phone`, `campus_id`)
SELECT
    `users`.`uuid`             AS `user_uuid`,
    'Guido'                    AS `first_name`,
    'Santella'                 AS `last_name`,
    'gsantella@southhills.edu' AS `email`,
    '8149446134'               AS `phone`,
    'alt'                      AS `campus_id`
FROM `users`
WHERE `users`.`identity` = 'gsantella';

INSERT INTO `instructors` (`user_uuid`, `first_name`, `last_name`, `email`, `phone`, `campus_id`)
SELECT
    `users`.`uuid`              AS `user_uuid`,
    'Bob'                       AS `first_name`,
    'Selfridge'                 AS `last_name`,
    'bselfridge@southhills.edu' AS `email`,
    '8149446134'                AS `phone`,
    'alt'                       AS `campus_id`
FROM `users`
WHERE `users`.`identity` = 'bselfridge';

INSERT INTO `instructors` (`user_uuid`, `first_name`, `last_name`, `email`, `phone`, `campus_id`)
SELECT
    `users`.`uuid`         AS `user_uuid`,
    'Nicholas'             AS `first_name`,
    'Page'                 AS `last_name`,
    'npage@southhills.edu' AS `email`,
    '8142377755'           AS `phone`,
    'sce'                  AS `campus_id`
FROM `users`
WHERE `users`.`identity` = 'npage';

INSERT INTO `instructors` (`user_uuid`, `first_name`, `last_name`, `email`, `phone`, `campus_id`)
SELECT
    `users`.`uuid`           AS `user_uuid`,
    'Rick'                   AS `first_name`,
    'Gority'                 AS `last_name`,
    'rgority@southhills.edu' AS `email`,
    '8142377755'             AS `phone`,
    'sce'                    AS `campus_id`
FROM `users`
WHERE `users`.`identity` = 'rgority';

INSERT INTO `companies` (`name`, `address`, `unit`, `city`, `state`, `zip`, `phone`)
VALUE
    ('PS Solutions', '350 Lakemont Park Blvd',  'Unit 2A', 'Altoona',      'PA', '16602', '8149427888'),
    ('CCSalesPro',   '117 Olde Farm Office Rd', NULL,      'Duncansville', 'PA', '16635', '7083075250');

INSERT INTO `supervisors` (`user_uuid`, `first_name`, `last_name`, `title`, `email`, `phone`, `company_uuid`)
SELECT
    `users`.`uuid`               AS `user_uuid`,
    'Ry'                         AS `first_name`,
    'Gallagher'                  AS `last_name`,
    'Programmer'                 AS `title`,
    'rgallagher@pssolutions.net' AS `email`,
    '8149427888'                 AS `phone`,
    `companies`.`uuid`           AS `company_uuid`
FROM `users`, `companies`
WHERE `users`.`identity` = 'rgallagher@pssolutions.net' AND `companies`.`name` = 'PS Solutions';

INSERT INTO `supervisors` (`user_uuid`, `first_name`, `last_name`, `title`, `email`, `phone`, `company_uuid`)
SELECT
    `users`.`uuid`               AS `user_uuid`,
    'Jack'                       AS `first_name`,
    'Christensen'                AS `last_name`,
    'Programmer'                 AS `title`,
    'jack@jackchristensen.com'   AS `email`,
    '7083075250'                 AS `phone`,
    `companies`.`uuid`           AS `company_uuid`
FROM `users`, `companies`
WHERE `users`.`identity` = 'jack@jackchristensen.com' AND `companies`.`name` = 'CCSalesPro';

INSERT INTO `programs` (`id`, `name`)
VALUES
    ('et',   'Engineering Technology'),
    ('ap',   'Administrative Professional'),
    ('baa',  'Business Administration - Accounting'),
    ('bamm', 'Business Administration - Management and Marketing'),
    ('cj',   'Criminal Justice'),
    ('dms',  'Diagnostic Medial Sonography'),
    ('ga',   'Graphic Arts'),
    ('it',   'Information Technology'),
    ('sdp',  'Software Development and Programming'),
    ('ma',   'Medical Assistant'),
    ('mcb',  'Medical Coding and Billing');

INSERT INTO `students` (`user_uuid`, `first_name`, `last_name`, `address`, `unit`, `city`, `state`, `zip`, `email`, `phone`, `campus_id`, `program_id`)
SELECT
    `users`.`uuid`              AS `user_uuid`,
    'Marcus'                    AS `first_name`,
    'Germano'                   AS `last_name`,
    '400 Grandview Rd'          AS `address`,
    'Apt 3'                     AS `unit`,
    'Altoona'                   AS `city`,
    'PA'                        AS `state`,
    '16601'                     AS `zip`,
    'mgermano79@southhills.edu' AS `email`,
    '8146310533'                AS `phone`,
    'alt'                       AS `campus_id`,
    'sdp'                       AS `program_id`
FROM `users`
WHERE `users`.`identity` = 'mgermano79';

INSERT INTO `students` (`user_uuid`, `first_name`, `last_name`, `address`, `city`, `state`, `zip`, `email`, `phone`, `campus_id`, `program_id`)
SELECT
    `users`.`uuid`               AS `user_uuid`,
    'Benjamin'                   AS `first_name`,
    'Kletzing'                   AS `last_name`,
    '148 Bradford Ln'            AS `address`,
    'Roaring Spring'             AS `city`,
    'PA'                         AS `state`,
    '16673'                      AS `zip`,
    'bkletzing32@southhills.edu' AS `email`,
    '8143093431'                 AS `phone`,
    'alt'                        AS `campus_id`,
    'sdp'                        AS `program_id`
FROM `users`
WHERE `users`.`identity` = 'bkletzing32';

-- +migrate Down
DELETE FROM `students`;

DELETE FROM `programs`;

DELETE FROM `supervisors`;

DELETE FROM `companies`;

DELETE FROM `instructors`;

DELETE FROM `campuses`;

DELETE FROM `administrators`;

DELETE FROM `users`;

DELETE FROM 'roles';