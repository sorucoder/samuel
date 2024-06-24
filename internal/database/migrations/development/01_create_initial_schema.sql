-- +migrate Up
CREATE TABLE `roles` (
    `id`
        VARCHAR(32)
        NOT NULL
        UNIQUE,
    `name`
        VARCHAR(64)
        NOT NULL,
    `priority`
        TINYINT UNSIGNED
        NOT NULL
        UNIQUE,
    PRIMARY KEY (`id`)
);

CREATE TABLE `users` (
    `uuid`
        CHAR(36)
        NOT NULL
        UNIQUE
        DEFAULT (UUID()),
    `identity`
        VARCHAR(254)
        NOT NULL,
    `password_hash`
        BINARY(60),
    `role_id`
        VARCHAR(32)
        NOT NULL,
    `created_on`
        DATETIME
        NOT NULL
        DEFAULT (NOW()),
    CONSTRAINT `check_only_supervisors_use_password_hash`
        CHECK (
            (`role_id` = 'supervisor' AND `password_hash` IS NOT NULL) OR
            (`role_id` != 'supervisor' AND `password_hash` IS NULL)
        ),
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (`role_id`)
        REFERENCES `roles`(`id`)
);

CREATE TABLE `sessions` (
    `token`
        CHAR(36)
        NOT NULL
        UNIQUE
        DEFAULT (UUID()),
    `user_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `started_on`
        DATETIME
        NOT NULL
        DEFAULT (NOW()),
    `expires_on`
        DATETIME
        NOT NULL
        DEFAULT (DATE_ADD(NOW(), INTERVAL 20 MINUTE)),
    PRIMARY KEY (`token`),
    FOREIGN KEY (`user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE
);

CREATE EVENT `event_delete_expired_sessions`
    ON SCHEDULE EVERY 40 MINUTE
DO
    DELETE FROM `sessions`
    WHERE `expired_on` < NOW();

CREATE TABLE `administrators` (
    `user_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `first_name`
        VARCHAR(64)
        NOT NULL,
    `last_name`
        VARCHAR(64)
        NOT NULL,
    `email`
        VARCHAR(254)
        NOT NULL,
    `phone`
        CHAR(10)
        NOT NULL
        CHECK (`phone` RLIKE '^[0-9]{10}$'),
    PRIMARY KEY (`user_uuid`),
    FOREIGN KEY (`user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE
);

CREATE TABLE `campuses` (
    `id`
        CHAR(3)
        NOT NULL
        UNIQUE
        CHECK (`id` RLIKE '^[a-z]{3}$'),
    `name`
        VARCHAR(64)
        NOT NULL,
    `address`
        TINYTEXT
        NOT NULL,
    `unit`
        TINYTEXT,
    `city`
        VARCHAR(64)
        NOT NULL,
    `state`
        ENUM(
            'AL', 'AK', 'AZ', 'AR', 'CA', 'CO', 'CT', 'DE', 'FL', 'GA',
            'HI', 'ID', 'IL', 'IN', 'IA', 'KS', 'KY', 'LA', 'ME', 'MD',
            'MA', 'MI', 'MN', 'MS', 'MO', 'MT', 'NE', 'NV', 'NH', 'NJ',
            'NM', 'NY', 'NC', 'ND', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC',
            'SD', 'TN', 'TX', 'UT', 'VT', 'VA', 'WA', 'WV', 'WI', 'WY'
        )
        NOT NULL,
    `zip`
        VARCHAR(10)
        NOT NULL
        CHECK (`zip` RLIKE '^[0-9]{5}(?:-[0-9]{4})?$'),
    `phone`
        CHAR(10)
        NOT NULL
        CHECK (`phone` RLIKE '^[0-9]{10}$'),
    PRIMARY KEY (`id`)
);

CREATE TABLE `instructors` (
    `user_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `first_name`
        VARCHAR(64)
        NOT NULL,
    `last_name`
        VARCHAR(64)
        NOT NULL,
    `email`
        VARCHAR(254)
        NOT NULL
        UNIQUE,
    `phone`
        CHAR(10)
        NOT NULL
        CHECK (`phone` RLIKE '^[0-9]{10}$'),
    `campus_id`
        VARCHAR(3)
        NOT NULL,
    PRIMARY KEY (`user_uuid`),
    FOREIGN KEY (`user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE,
    FOREIGN KEY (`campus_id`)
        REFERENCES `campuses`(`id`)
);

CREATE TABLE `programs` (
    `id`
        VARCHAR(4)
        NOT NULL
        UNIQUE
        CHECK (`id` RLIKE '^[a-z]{2,4}$'),
    `name`
        TINYTEXT
        NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `coordinators` (
    `campus_id`
        VARCHAR(3)
        NOT NULL,
    `program_id`
        VARCHAR(4)
        NOT NULL,
    `instructor_uuid`
        CHAR(36)
        NOT NULL,
    UNIQUE (`campus_id`, `program_id`),
    PRIMARY KEY (`campus_id`, `program_id`),
    FOREIGN KEY (`campus_id`)
        REFERENCES `campuses`(`id`),
    FOREIGN KEY (`program_id`)
        REFERENCES `programs`(`id`),
    FOREIGN KEY (`instructor_uuid`)
        REFERENCES `instructors`(`user_uuid`)
        ON DELETE CASCADE
);

CREATE TABLE `companies` (
    `uuid`
        CHAR(36)
        NOT NULL
        UNIQUE
        DEFAULT (UUID()),
    `name`
        VARCHAR(64)
        NOT NULL,
    `address`
        TINYTEXT
        NOT NULL,
    `unit`
        TINYTEXT,
    `city`
        VARCHAR(64)
        NOT NULL,
    `state`
        ENUM(
            'AL', 'AK', 'AZ', 'AR', 'CA', 'CO', 'CT', 'DE', 'FL', 'GA',
            'HI', 'ID', 'IL', 'IN', 'IA', 'KS', 'KY', 'LA', 'ME', 'MD',
            'MA', 'MI', 'MN', 'MS', 'MO', 'MT', 'NE', 'NV', 'NH', 'NJ',
            'NM', 'NY', 'NC', 'ND', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC',
            'SD', 'TN', 'TX', 'UT', 'VT', 'VA', 'WA', 'WV', 'WI', 'WY'
        )
        NOT NULL,
    `zip`
        VARCHAR(10)
        NOT NULL
        CHECK (`zip` RLIKE '^[0-9]{5}(?:-[0-9]{4})?$'),
    `phone`
        CHAR(10)
        NOT NULL
        CHECK (`phone` RLIKE '^[0-9]{10}$'),
    PRIMARY KEY (`uuid`)
);

CREATE TABLE `supervisors` (
    `user_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `first_name`
        VARCHAR(64)
        NOT NULL,
    `last_name`
        VARCHAR(64)
        NOT NULL,
    `title`
        VARCHAR(64)
        NOT NULL,
    `email`
        VARCHAR(254)
        NOT NULL
        UNIQUE,
    `phone`
        CHAR(10)
        NOT NULL
        CHECK (`phone` RLIKE '^[0-9]{10}$'),
    `company_uuid`
        CHAR(36)
        NOT NULL,
    PRIMARY KEY (`user_uuid`),
    FOREIGN KEY (`user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE,
    FOREIGN KEY (`company_uuid`)
        REFERENCES `companies`(`uuid`)
);

CREATE TABLE `password_changes` (
    `token`
        CHAR(36)
        NOT NULL
        UNIQUE
        DEFAULT (UUID()),
    `supervisor_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `expires_on`
        DATETIME
        NOT NULL
        DEFAULT (DATE_ADD(NOW(), INTERVAL 5 MINUTE)),
    PRIMARY KEY (`token`),
    FOREIGN KEY (`supervisor_uuid`)
        REFERENCES `supervisors`(`user_uuid`)
        ON DELETE CASCADE
);

CREATE EVENT `event_delete_expired_password_changes`
    ON SCHEDULE EVERY 10 MINUTE
DO
    DELETE FROM `password_changes`
    WHERE `expired_on` < NOW();

CREATE TABLE `students` (
    `user_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `first_name`
        VARCHAR(64)
        NOT NULL,
    `last_name`
        VARCHAR(64)
        NOT NULL,
    `address`
        TINYTEXT
        NOT NULL,
    `unit`
        TINYTEXT,
    `city`
        VARCHAR(64)
        NOT NULL,
    `state`
        ENUM(
            'AL', 'AK', 'AZ', 'AR', 'CA', 'CO', 'CT', 'DE', 'FL', 'GA',
            'HI', 'ID', 'IL', 'IN', 'IA', 'KS', 'KY', 'LA', 'ME', 'MD',
            'MA', 'MI', 'MN', 'MS', 'MO', 'MT', 'NE', 'NV', 'NH', 'NJ',
            'NM', 'NY', 'NC', 'ND', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC',
            'SD', 'TN', 'TX', 'UT', 'VT', 'VA', 'WA', 'WV', 'WI', 'WY'
        )
        NOT NULL,
    `zip`
        VARCHAR(10)
        NOT NULL
        CHECK (`zip` RLIKE '^[0-9]{5}(?:-[0-9]{4})?$'),
    `email`
        VARCHAR(256)
        NOT NULL
        UNIQUE,
    `phone`
        CHAR(10)
        NOT NULL
        CHECK (`phone` RLIKE '^[0-9]{10}$'),
    `campus_id`
        CHAR(3)
        NOT NULL,
    `program_id`
        VARCHAR(4)
        NOT NULL,
    PRIMARY KEY (`user_uuid`),
    FOREIGN KEY (`user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE,
    FOREIGN KEY (`campus_id`)
        REFERENCES `campuses`(`id`),
    FOREIGN KEY (`program_id`)
        REFERENCES `programs`(`id`)
);

CREATE TABLE `internships` (
    `uuid`
        CHAR(36)
        NOT NULL
        UNIQUE
        DEFAULT (UUID()),
    `student_uuid`
        CHAR(36)
        NOT NULL,
    `instructor_uuid`
        CHAR(36)
        NOT NULL,
    `supervisor_uuid`
        CHAR(36)
        NOT NULL,
    `start_on`
        DATE
        NOT NULL
        CHECK (DAYOFWEEK(`start_on`) = 1),
    `end_on`
        DATE
        NOT NULL
        CHECK (DAYOFWEEK(`end_on`) = 7),
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (`student_uuid`)
        REFERENCES `students`(`user_uuid`),
    FOREIGN KEY (`instructor_uuid`)
        REFERENCES `instructors`(`user_uuid`),
    FOREIGN KEY (`supervisor_uuid`)
        REFERENCES `supervisors`(`user_uuid`)
);

CREATE TABLE `timecards` (
    `internship_uuid`
        CHAR(36)
        NOT NULL,
    `week_of`
        DATE
        NOT NULL
        CHECK (DAYOFWEEK(`week_of`) = 1),
    `sunday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`sunday_hours` BETWEEN 0.00 AND 12.00),
    `monday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`monday_hours` BETWEEN 0.00 AND 12.00),
    `tuesday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`tuesday_hours` BETWEEN 0.00 AND 12.00),
    `wednesday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`wednesday_hours` BETWEEN 0.00 AND 12.00),
    `thursday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`thursday_hours` BETWEEN 0.00 AND 12.00),
    `friday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`friday_hours` BETWEEN 0.00 AND 12.00),
    `saturday_hours`
        DEC(4, 2)
        NOT NULL
        DEFAULT 0.00
        CHECK (`saturday_hours` BETWEEN 0.00 AND 12.00),
    `status`
        ENUM('submitted', 'approved', 'denied'),
    `status_changed_on`
        DATETIME,
    CONSTRAINT `valid_total_hours`
        CHECK (
            (
                `sunday_hours`    +
                `monday_hours`    +
                `tuesday_hours`   +
                `wednesday_hours` +
                `thursday_hours`  +
                `friday_hours`    +
                `saturday_hours`
            )
            BETWEEN 0.00 AND 72.00
        ),
    UNIQUE (`internship_uuid`, `week_of`),
    PRIMARY KEY (`internship_uuid`, `week_of`)
);

CREATE TABLE `supervisor_reports` (
    `internship_uuid`
        CHAR(36)
        NOT NULL,
    `week_of`
        DATE
        NOT NULL
        CHECK (DAYOFWEEK(`week_of`) = 1),
    `submitted_on`
        DATETIME
        NOT NULL,
    `knowledge_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `knowledge_response`
        TEXT,
    `quality_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `prioritization_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `quality_response`
        TEXT,
    `efficiency_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `multitasking_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `efficiency_response`
        TEXT,
    `communication_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `listening_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `communication_response`
        TEXT,
    `aptitude_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `inquisitiveness_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `aptitude_response`
        TEXT,
    `initiative_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `initiative_response`
        TEXT,
    `cooperation_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `attitude_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `cooperation_response`
        TEXT,
    `attendance_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `notification_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `attendance_response`
        TEXT,
    `professionalism_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `apperance_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `professionalism_response`
        TEXT,
    `overall_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `overall_response`
        TEXT,
    `accomplishment_response`
        TEXT,
    `requests_phone_call`
        BOOL
        NOT NULL
        DEFAULT FALSE,
    `visible_to_student`
        BOOL
        NOT NULL
        DEFAULT FALSE,
    UNIQUE (`internship_uuid`, `week_of`),
    PRIMARY KEY (`internship_uuid`, `week_of`),
    FOREIGN KEY (`internship_uuid`)
        REFERENCES `internships`(`uuid`)
);

CREATE TABLE `program_evaluation_questions` (
    `program_id`
        VARCHAR(4)
        NOT NULL,
    `number`
        TINYINT UNSIGNED
        NOT NULL
        CHECK (`number` > 1),
    `question`
        TEXT,
    UNIQUE (`program_id`, `number`),
    PRIMARY KEY (`program_id`, `number`),
    FOREIGN KEY (`program_id`)
        REFERENCES `programs`(`id`)
);

CREATE TABLE `supervisor_general_evaluations` (
    `internship_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `submitted_on`
        DATETIME
        NOT NULL,
    `knowledge_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `quality_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `efficiency_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `communication_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `aptitude_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `initiative_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `attitude_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `attendance_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `professionalism_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `overall_rating`
        ENUM('unacceptable', 'poor', 'satisfactory', 'good', 'superior'),
    `strengths_response`
        TEXT,
    `weaknesses_response`
        TEXT,
    `academic_suggestions_response`
        TEXT,
    `value_response`
        TEXT,
    `recommends_employment`
        BOOL
        NOT NULL,
    `recommendation_response`
        TEXT,
    `visible_to_student`
        BOOL
        NOT NULL
        DEFAULT FALSE,
    PRIMARY KEY (`internship_uuid`),
    FOREIGN KEY (`internship_uuid`)
        REFERENCES `internships`(`uuid`)
);

CREATE TABLE `supervisor_program_evaluation_responses` (
    `internship_uuid`
        CHAR(36)
        NOT NULL,
    `question_program_id`
        VARCHAR(4)
        NOT NULL,
    `question_number`
        TINYINT UNSIGNED
        NOT NULL,
    `response`
        TEXT,
    UNIQUE (`internship_uuid`, `question_program_id`, `question_number`),
    PRIMARY KEY (`internship_uuid`, `question_program_id`, `question_number`),
    FOREIGN KEY (`internship_uuid`)
        REFERENCES `internships`(`uuid`),
    FOREIGN KEY (`question_program_id`, `question_number`)
        REFERENCES `program_evaluation_questions`(`program_id`, `number`)
);

CREATE TABLE `student_reports` (
    `internship_uuid`
        CHAR(36)
        NOT NULL,
    `week_of`
        DATE
        NOT NULL
        CHECK (DAYOFWEEK(`week_of`) = 1),
    `submitted_on`
        DATETIME
        NOT NULL,
    `major_objectives_response`
        TEXT,
    `additional_accomplishments_response`
        TEXT,
    `unassigned_tasks_response`
        TEXT,
    `well_handled_activity_response`
        TEXT,
    `helpfulness_and_issues_response`
        TEXT,
    `problem_solving_response`
        TEXT,
    `learning_response`
        TEXT,
    UNIQUE (`internship_uuid`, `week_of`),
    PRIMARY KEY (`internship_uuid`, `week_of`),
    FOREIGN KEY (`internship_uuid`)
        REFERENCES `internships`(`uuid`)
);

CREATE TABLE `student_evaluations` (
    `internship_uuid`
        CHAR(36)
        NOT NULL
        UNIQUE,
    `submitted_on`
        DATETIME
        NOT NULL,
    `company_information_response`
        TEXT,
    `major_responsibilities_response`
        TEXT,
    `accomplishment_response`
        TEXT,
    `academic_training_benefits_response`
        TEXT,
    `academic_training_improvements_response`
        TEXT,
    `skill_development_response`
        TEXT,
    `attitude_change_response`
        TEXT,
    `comments_response`
        TEXT,
    PRIMARY KEY (`internship_uuid`),
    FOREIGN KEY (`internship_uuid`)
        REFERENCES `internships`(`uuid`)
);

CREATE TABLE `notifications` (
    `uuid`
        CHAR(36)
        NOT NULL
        UNIQUE
        DEFAULT (UUID()),
    `from_user_uuid`
        CHAR(36)
        NOT NULL,
    `to_user_uuid`
        CHAR(36)
        NOT NULL,
    `message`
        TEXT
        NOT NULL,
    `sent_on`
        DATETIME
        NOT NULL
        DEFAULT (NOW()),
    `seen`
        BOOL
        NOT NULL
        DEFAULT FALSE,
    `seen_on`
        DATETIME,
    `type`
        ENUM('system', 'personal')
        NOT NULL
        DEFAULT 'system',
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (`from_user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE,
    FOREIGN KEY (`to_user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE CASCADE
);

CREATE TABLE `audit` (
    `id`
        SERIAL,
    `description`
        TEXT
        NOT NULL,
    `user_uuid`
        CHAR(36),
    `timestamp`
        DATETIME
        NOT NULL
        DEFAULT (NOW()),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_uuid`)
        REFERENCES `users`(`uuid`)
        ON DELETE SET NULL
);

-- +migrate Down
DROP TABLE `audit`;

DROP TABLE `notifications`;

DROP TABLE `student_evaluations`;

DROP TABLE `student_reports`;

DROP TABLE `supervisor_program_evaluation_responses`;

DROP TABLE `supervisor_general_evaluations`;

DROP TABLE `program_evaluation_questions`;

DROP TABLE `supervisor_reports`;

DROP TABLE `timecards`;

DROP TABLE `internships`;

DROP TABLE `students`;

DROP TABLE `password_resets`;

DROP TABLE `supervisors`;

DROP TABLE `companies`;

DROP TABLE `coordinators`;

DROP TABLE `programs`;

DROP TABLE `instructors`;

DROP TABLE `campuses`;

DROP TABLE `administrators`;

DROP EVENT `event_delete_expired_password_changes`;

DROP EVENT `event_delete_expired_sessions`;

DROP TABLE `sessions`;

DROP TABLE `users`;

DROP TABLE `roles`;